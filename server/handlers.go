package server

import (
	"net/http"
	"path/filepath"
	"html/template"
	"log"
	"dexchan/server/model"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"io"
	"os"
	"time"
	"html"
)

// TEMP: always parse template, makes editing easier during development
const alwaysParseTemplate = true
const baseTemplate = "templates/base.html"

type State struct {
	Boards  []*model.Board
	Threads []*model.Thread
	Board   *model.Board
	Thread  *model.Thread
	Posts   []*model.Post
}

func (s *Server) getCachedTemplate(name string) *template.Template {
	t, ok := s.templatesCache[name]
	if !ok || alwaysParseTemplate {
		var err error
		t, err = template.ParseFiles(baseTemplate, name)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	} else {
		s.templatesCache[name] = t
	}

	return t
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "index.html")
	t := s.getCachedTemplate(fp)

	state := State{
		Boards: s.boards,
	}
	t.ExecuteTemplate(w, "base", state)
}

func (s *Server) tryGetBoard(w http.ResponseWriter, r *http.Request) *model.Board {
	v := mux.Vars(r)

	bn := v["board"]
	found, ok := s.boardsByShortcode[bn]

	if !ok {
		http.Error(w,
			fmt.Sprintf("bad request: board %s doesnt exist", bn), http.StatusBadRequest)
		return nil
	}

	return found
}

func (s *Server) tryGetBoardAndThread(w http.ResponseWriter, r *http.Request) (*model.Board, *model.Thread) {
	b := s.tryGetBoard(w, r)
	if b == nil {
		return nil, nil
	}

	v := mux.Vars(r)
	tns := v["thread"]
	// matched by regex, always good
	tn, _ := strconv.ParseInt(tns, 10, 64)
	thread, err := s.Db.GetThread(tn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil
	}

	return b, thread
}

func (s *Server) boardHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "board.html")
	t := s.getCachedTemplate(fp)

	b := s.tryGetBoard(w, r)
	tx, err := s.Db.GetThreads(b.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := State{
		Boards:  s.boards,
		Board:   b,
		Threads: tx,
	}
	t.ExecuteTemplate(w, "base", state)
}

func (s *Server) threadHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "thread.html")
	t := s.getCachedTemplate(fp)

	b, thread := s.tryGetBoardAndThread(w, r)
	if b == nil || thread == nil {
		return
	}

	posts, _ := s.Db.GetPosts(thread.ID)

	state := State{
		Boards: s.boards,
		Board:  b,
		Thread: thread,
		Posts:  posts,
	}

	t.ExecuteTemplate(w, "base", state)
}

func (s *Server) replyHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "reply.html")
	t := s.getCachedTemplate(fp)

	b, thread := s.tryGetBoardAndThread(w, r)
	if b == nil || thread == nil {
		return
	}

	posts, _ := s.Db.GetPosts(thread.ID)

	state := State{
		Boards: s.boards,
		Board:  b,
		Thread: thread,
		Posts:  posts,
	}
	t.ExecuteTemplate(w, "base", state)
}

func (s *Server) addReplyHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	bn := v["board"]
	tidstr := v["thread"]
	tid, _ := strconv.ParseInt(tidstr, 10, 64)
	newPost := &model.Post{}

	// todo: lookup board max
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*8)
	sentFile, fileHeader, err := r.FormFile("f")
	if err == nil {
		// todo: abstract
		saveDir := filepath.Join(s.Config.StaticDir, bn)
		if err := os.MkdirAll(saveDir, 0700); err != nil {
			return
		}

		// ensure dir for board content exists
		// security: we could get in a bogus header with traversal exploitation so chop to base
		// dont use original filename, timestamp fn := filepath.Base(fileHeader.Filename)
		origExt := filepath.Ext(fileHeader.Filename)
		fn := fmt.Sprintf("%d%s", time.Now().Unix(), origExt)
		saveFn := filepath.Join(saveDir, fn)
		saveFile, err := os.Create(saveFn)
		if err != nil {
			return
		}

		newPost.OriginalFilename = filepath.Base(fileHeader.Filename)
		newPost.Size = fileHeader.Size
		// !! CDN HOOK
		newPost.Location = filepath.Join("/", saveFn)
		// TODO: thumbnail creation
		newPost.ThumbnailLocation = newPost.Location

		// start saving it to disk. this is dumb atm
		go func() {
			defer r.Body.Close()
			defer sentFile.Close()
			defer saveFile.Close()
			// TODO: partial writes
			io.Copy(saveFile, sentFile)
		}()
	}

	postContent := r.FormValue("post")
	newPost.Content = html.EscapeString(postContent)
	newPost.ThreadID = tid
	newPost.PostedAt = time.Now()

	// TODO: inet only handles ipv4. argh, inet6 localhost coming in
	// not sure if i should even bother with ip based user ids. might have to text
	_, err = s.Db.CreatePost(newPost, "192.168.1.1") //r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: hardcoded proto atm :-(
	url := fmt.Sprintf("http://%s/%s/%d", r.Host, bn, tid)
	http.Redirect(w, r, url, http.StatusSeeOther)
	//b, thread := s.tryGetBoardAndThread(w, r)
	//if b == nil || thread == nil {
	//	return
	//}
}
