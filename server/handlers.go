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
	"io/ioutil"
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

	for _, t := range tx {
		parseThreadHelpers(t)
	}

	state := State{
		Boards:  s.boards,
		Board:   b,
		Threads: tx,
	}
	t.ExecuteTemplate(w, "base", state)
}

// helpers for templates
func parsePostHelpers(posts []*model.Post) {
	for _, p := range posts {
		p.HasAttachment = p.Location != "" && p.Size > 0
		if p.HasAttachment {
			ext := filepath.Ext(p.OriginalFilename)
			p.IsVideo = ext == ".webm"
		}
	}
}

func parseThreadHelpers(t *model.Thread) {
	// note: should always be true
	t.HasAttachment = t.Location != "" && t.Size > 0
	ext := filepath.Ext(t.OriginalFilename)
	t.IsVideo = ext == ".webm"
}

func (s *Server) threadHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "thread.html")
	t := s.getCachedTemplate(fp)

	b, thread := s.tryGetBoardAndThread(w, r)
	if b == nil || thread == nil {
		return
	}
	parseThreadHelpers(thread)

	posts, _ := s.Db.GetPosts(thread.ID)
	parsePostHelpers(posts)

	state := State{
		Boards: s.boards,
		Board:  b,
		Thread: thread,
		Posts:  posts,
	}

	t.ExecuteTemplate(w, "base", state)
}

func (s *Server) createThreadHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "createthread.html")
	t := s.getCachedTemplate(fp)

	b := s.tryGetBoard(w, r)
	if b == nil {
		return
	}

	state := State{
		Boards: s.boards,
		Board:  b,
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

	state := State{
		Boards: s.boards,
		Board:  b,
		Thread: thread,
	}
	t.ExecuteTemplate(w, "base", state)
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

	// todo: lookup board max
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*8)
	sentFile, fileHeader, err := r.FormFile("f")

	// attachment is always required for new post
	if err != nil {
		http.Error(w, "a file is required to create a thread", http.StatusBadRequest)
		return
	}

	// IMPORTANT: if you dont close body, redirect wont work
	r.Body.Close()
	newThread := &model.Thread{
		Subject:     sub,
		Description: desc,
		CreatedAt:   time.Now(),
		BoardID:     b.ID,
	}

	// todo: abstract
	saveDir := filepath.Join(s.Config.StaticDir, bn)
	if err := os.MkdirAll(saveDir, 0700); err != nil {
		return
	}

	// ensure dir for board content exists
	// security: we could get in a bogus header with traversal exploitation so chop to base
	// dont use original filename, timestamp fn := filepath.Base(fileHeader.Filename)
	origExt := filepath.Ext(fileHeader.Filename)
	// save this once so we dont get drift when making thumbnail
	timestamp := time.Now().Unix()
	fn := fmt.Sprintf("%d%s", timestamp, origExt)
	saveFn := filepath.Join(saveDir, fn)
	saveFile, err := os.Create(saveFn)
	if err != nil {
		return
	}

	newThread.OriginalFilename = filepath.Base(fileHeader.Filename)
	newThread.Size = fileHeader.Size
	// OUCH: cant use filename. its not safe cross platform. windows server will create
	// \static\bn\fn.ext which wont work as an url of course. build manually :-(
	// !! CDN HOOK
	newThread.Location = fmt.Sprintf("/%s/%s/%s", s.Config.StaticDir, bn, fn)

	// start saving it to disk. this is dumb atm
	// dont try anything cute like goroutines otherwise go will return a 200 for you
	// TODO: partial writes
	io.Copy(saveFile, sentFile)
	sentFile.Close()
	if origExt != ".webm" {
		// we need to reset file pos so resizer starts at correct spot
		saveFile.Seek(0, 0)
		tnBytes, err := s.createThumbnail(saveFile)
		saveFile.Close()
		if err == nil {
			tfn := fmt.Sprintf("%d_tn%s", timestamp, origExt)
			newThread.ThumbnailLocation = fmt.Sprintf("/%s/%s/%s", s.Config.StaticDir, bn, tfn)
			ioutil.WriteFile(filepath.Join(saveDir, tfn), tnBytes, 0600)
		}
	} else {
		newThread.ThumbnailLocation = newThread.Location
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
		// OUCH: cant use filename. its not safe cross platform. windows server will create
		// \static\bn\fn.ext which wont work as an url of course. build manually :-(
		// !! CDN HOOK
		newPost.Location = fmt.Sprintf("/%s/%s/%s", s.Config.StaticDir, bn, fn)
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

	createdID, err := s.Db.CreatePost(newPost, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: hardcoded proto atm :-(
	url := fmt.Sprintf("http://%s/%s/%d#%d", r.Host, bn, tid, createdID)
	http.Redirect(w, r, url, http.StatusSeeOther)
	//b, thread := s.tryGetBoardAndThread(w, r)
	//if b == nil || thread == nil {
	//	return
	//}
}
