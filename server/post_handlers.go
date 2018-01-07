package server

import (
	"github.com/gorilla/mux"
	"dexchan/server/model"
	"time"
	"path/filepath"
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"net/http"
	"strings"
	"errors"
	"log"
)

func validAttachmentType(fn string, validAttachmentTypes string) bool {
	origExt := filepath.Ext(fn)
	valid := strings.Split(validAttachmentTypes, ",")

	for _, v := range valid {
		if v == origExt {
			return true
		}
	}

	return false
}

func isAttachmentRequired(a model.AttachmentEntity) bool {
	_, ok := a.(*model.Thread)
	return ok
}

func (s *Server) ensureBoardSaveDir(bn string) (string, error) {
	saveDir := filepath.Join(s.Config.StaticDir, bn)
	if err := os.MkdirAll(saveDir, 0700); err != nil && err != os.ErrExist {
		return saveDir, err
	}
	return saveDir, nil
}

func createAttachmentName(origFileName, saveDir string) (int64, string, string) {
	// security: we could get in a bogus header with traversal exploitation so chop to base
	// dont use original filename, use timestamp
	// we return the timestamp we determined for file so we dont get drift when creating
	// thumbnail for it later
	timestamp := time.Now().Unix()
	origExt := filepath.Ext(origFileName)
	fn := fmt.Sprintf("%d%s", timestamp, origExt)
	saveFn := filepath.Join(saveDir, fn)
	return timestamp, origExt, saveFn
}

func (s *Server) createLocationLink(bn, ext string, timestamp int64, isThumbnail bool) string {
	// cant use filepath.Join here, on windows it will use wrong slashes
	base := fmt.Sprintf("/%s/%s/", s.Config.StaticDir, bn)
	fmtstring := ""
	if isThumbnail {
		fmtstring = "%s%d_tn%s"
	} else {
		fmtstring = "%s%d%s"
	}
	return fmt.Sprintf(fmtstring, base, timestamp, ext)
}

// handleAttachment returns true if there was a an attachment in response
func (s *Server) handleAttachment(w http.ResponseWriter, r *http.Request, a model.AttachmentEntity) (bool, error) {
	v := mux.Vars(r)
	bn := v["board"]
	b := s.tryGetBoard(w, r)
	defer r.Body.Close()

	if b == nil {
		// this means no board was found.
		return false, errors.New("no board found")
	}

	r.Body = http.MaxBytesReader(w, r.Body, b.MaxAttachmentSize)
	sentFile, fileHeader, err := r.FormFile("f")

	if err != nil {
		if isAttachmentRequired(a) {
			return false, errors.New("a file is required to create a thread")
		}
		return false, nil
	}

	defer sentFile.Close()

	if !validAttachmentType(fileHeader.Filename, b.AttachmentTypes) {
		return true, errors.New("invalid attachment type")
	}

	var saveDir string
	if saveDir, err = s.ensureBoardSaveDir(bn); err != nil {
		return true, errors.New(fmt.Sprintf("failed to create board save dir: %s", err))
	}

	timestamp, ext, saveFn := createAttachmentName(fileHeader.Filename, saveDir)

	a.ParseFromHeader(fileHeader)
	a.SetLocation(s.createLocationLink(bn, ext, timestamp, false))

	saveFile, err := os.Create(saveFn)
	if err != nil {
		return true, err
	}
	defer saveFile.Close()

	// dont try anything cute like goroutines otherwise http will return a 200 for you
	// TODO: partial writes
	io.Copy(saveFile, sentFile)

	if ext != ".webm" {
		// we need to reset file pos so resizer starts at correct spot
		saveFile.Seek(0, 0)
		tnBytes, err := createThumbnail(saveFile)
		if err == nil {
			a.SetThumbnail(s.createLocationLink(bn, ext, timestamp, true))
			ioutil.WriteFile(filepath.Join(saveDir, a.GetThumbnail()), tnBytes, 0600)
		} else {
			log.Printf("Warning: failed to generate a thumbnail: %s\n", err)
		}
	} else {
		a.SetThumbnail(a.GetLocation())
	}

	return true, nil
}

func (s *Server) addThreadHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	bn := v["board"]

	b := s.tryGetBoard(w, r)

	sub := r.FormValue("subject")
	desc := r.FormValue("description")
	if sub == "" || desc == "" {
		http.Error(w, "subject and description required to create a thread", http.StatusBadRequest)
		return
	}

	newThread := &model.Thread{
		Subject:     sub,
		Description: desc,
		CreatedAt:   time.Now(),
		BoardID:     b.ID,
	}

	if _, err := s.handleAttachment(w, r, newThread); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdID, err := s.Db.CreateThread(newThread, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: hardcoded proto atm :-(
	url := fmt.Sprintf("http://%s/%s/%d", r.Host, bn, createdID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (s *Server) addReplyHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	bn := v["board"]
	tidstr := v["thread"]
	tid, _ := strconv.ParseInt(tidstr, 10, 64)
	newPost := &model.Post{}

	hadAttachment, err := s.handleAttachment(w, r, newPost)

	if hadAttachment && err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postContent := r.FormValue("post")
	if postContent == "" && !hadAttachment {
		http.Error(w, "A reply requires either post content or an attachment!", http.StatusBadRequest)
		return
	}

	newPost.Content = postContent
	newPost.ThreadID = tid
	newPost.PostedAt = time.Now()

	createdID, err := s.Db.CreatePost(newPost, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: hardcoded proto atm :-(
	url := fmt.Sprintf("http://%s/%s/%d#%d", r.Host, bn, tid, createdID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
