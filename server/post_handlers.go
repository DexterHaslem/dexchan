// dexChan copyright Dexter Haslem <dmh@fastmail.com> 2018
// This file is part of dexChan
//
// dexChan is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dexChan is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dexChan. If not, see <http://www.gnu.org/licenses/>.

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
	"mime/multipart"
	"net/url"
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

type fileHandle interface {
	io.ReadCloser
}

func fileInfoFromURL(urlStr string, max int64) (fileHandle, *multipart.FileHeader, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	getReq, err := http.Get(u.String())
	if err != nil {
		return nil, nil, err
	}

	h := &multipart.FileHeader{}
	lenStr := getReq.Header.Get("Content-Length")
	h.Size, _ = strconv.ParseInt(lenStr, 10, 64)

	if h.Size > max {
		// no easy way to do this elsewhere, but since we have length
		// we can optimistically handle
		return nil, nil, errors.New("file greater than max")
	}
	// luckily filepath.Base works on urls
	h.Filename = filepath.Base(urlStr)

	return getReq.Body, h, err
}

func getFileInfo(w http.ResponseWriter, r *http.Request, maxSize int64) (fileHandle, *multipart.FileHeader, error) {
	tryUrl := r.FormValue("url")
	if tryUrl != "" {
		return fileInfoFromURL(tryUrl, maxSize)
	}
	// set this before so we get error before trying to read it
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	sentFile, fileHeader, err := r.FormFile("f")
	return sentFile, fileHeader, err
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

	file, header, err := getFileInfo(w, r, b.MaxAttachmentSize)

	if err != nil {
		if isAttachmentRequired(a) {
			return false, errors.New("a file is required to create a thread")
		}
		return false, nil
	}

	defer file.Close()

	if !validAttachmentType(header.Filename, b.AttachmentTypes) {
		return true, errors.New("invalid attachment type")
	}

	var saveDir string
	if saveDir, err = s.ensureBoardSaveDir(bn); err != nil {
		return true, errors.New(fmt.Sprintf("failed to create board save dir: %s", err))
	}

	timestamp, ext, saveFn := createAttachmentName(header.Filename, saveDir)

	a.ParseFromHeader(header)
	a.SetLocation(s.createLocationLink(bn, ext, timestamp, false))

	saveFile, err := os.Create(saveFn)
	if err != nil {
		return true, err
	}
	defer saveFile.Close()

	// dont try anything cute like goroutines otherwise http will return a 200 for you
	// TODO: partial writes
	io.Copy(saveFile, file)

	if ext != ".webm" {
		// we need to reset file pos so resizer starts at correct spot
		saveFile.Seek(0, 0)
		tnBytes, err := createThumbnail(saveFile)
		if err == nil {
			// this stinks, tn location will be server root / url, not filepath so we cant use it
			// if the server is running windows, so try to use massage. this trick
			// will rebase to our current dir and flip the path to whatever host OS is
			a.SetThumbnail(s.createLocationLink(bn, ext, timestamp, true))
			tnFile := filepath.Clean("." + a.GetThumbnail())
			ioutil.WriteFile(tnFile, tnBytes, 0600)
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

	createdID, err := s.Db.CreateThread(newThread, getRemoteIP(r))
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

	createdID, err := s.Db.CreatePost(newPost, getRemoteIP(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: hardcoded proto atm :-(
	url := fmt.Sprintf("http://%s/%s/%d#%d", r.Host, bn, tid, createdID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func getRemoteIP(r *http.Request) string {
	// we cant always use remote addr
	// if we're behind a reverse proxy we want to use
	// X-Real-IP if present. this isnt fool proof, if a third party hit
	// server directly they could spoof it of course

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	return r.RemoteAddr
}
