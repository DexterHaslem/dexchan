package server

import (
	"net/http"
	"dexchan/server/cfg"
	"strconv"
	dexDB "dexchan/server/db"
	"github.com/gorilla/mux"
	"dexchan/server/model"
	"log"
	"html/template"
)

type Server struct {
	Config *cfg.C
	Db     dexDB.D

	boards            []*model.Board
	boardsByShortcode map[string]*model.Board

	templatesCache map[string]*template.Template
}

func (s *Server) staticRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func (s *Server) Start() {
	s.templatesCache = map[string]*template.Template{}
	fs := http.FileServer(http.Dir(s.Config.StaticDir))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	router := mux.NewRouter()

	router.HandleFunc("/", s.homeHandler).Methods("GET")
	router.HandleFunc("/{board:[a-z]+}", s.boardHandler).Methods("GET")
	router.HandleFunc("/{board:[a-z]+}/createthread", s.createThreadHandler).Methods("GET")
	router.HandleFunc("/{board:[a-z]+}/{thread:[0-9]+}", s.threadHandler).Methods("GET")
	router.HandleFunc("/{board:[a-z]+}/{thread:[0-9]+}/reply", s.replyHandler).Methods("GET")

	router.HandleFunc("/{board:[a-z]+}/newthread", s.addThreadHandler).Methods("POST")
	router.HandleFunc("/{board:[a-z]+}/{thread:[0-9]+}/newreply", s.addReplyHandler).Methods("POST")

	http.Handle("/", router)

	// get boards right away
	b, err := s.Db.GetBoards()
	if err != nil {
		log.Fatal("failed to initialize")
	}
	s.boards = b
	s.boardsByShortcode = make(map[string]*model.Board)
	// build map by shortcode right away since we use this for all client rquests
	for _, b := range s.boards {
		s.boardsByShortcode[b.ShortCode] = b
	}
	http.ListenAndServe(":"+strconv.Itoa(s.Config.Port), nil)
}
