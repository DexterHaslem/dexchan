package server

import (
	"net/http"
	"path/filepath"
	"html/template"
	"log"
	"dexchan/server/model"
)

// TEMP: always parse template, makes editing easier during development
const alwaysParseTemplate = true

func (s *Server) getCachedTemplate(name string) *template.Template {

	if alwaysParseTemplate {
		t, err := template.ParseFiles(name)
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		return t
	}

	t, ok := s.templatesCache[name]
	if !ok {
		var err error
		t, err = template.ParseFiles(name)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	} else {
		s.templatesCache[name] = t
	}

	return t
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "index.html")
	t := s.getCachedTemplate(fp)

	t.Execute(w, struct {
		Boards []*model.Board
	}{s.boards})
}
