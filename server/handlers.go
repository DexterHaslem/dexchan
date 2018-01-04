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
)

// TEMP: always parse template, makes editing easier during development
const alwaysParseTemplate = true
const baseTemplate = "templates/base.html"

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

	t.ExecuteTemplate(w, "base", struct {
		Boards []*model.Board
	}{s.boards})
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

func (s *Server) boardHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "board.html")
	t := s.getCachedTemplate(fp)

	b := s.tryGetBoard(w, r)
	tx, err := s.Db.GetThreads(b.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", struct {
		Board   *model.Board
		Threads []*model.Thread
	}{b, tx})
}

func (s *Server) threadHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "thread.html")
	t := s.getCachedTemplate(fp)

	b := s.tryGetBoard(w, r)

	v := mux.Vars(r)
	tns := v["thread"]
	// matched by regex, always good
	tn, _ := strconv.ParseInt(tns, 10, 64)
	thread, err := s.Db.GetThread(tn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, _ := s.Db.GetPosts(tn)

	t.ExecuteTemplate(w, "base", struct {
		Board  *model.Board
		Thread *model.Thread
		Posts  []*model.Post
	}{b, thread, posts})
}