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
