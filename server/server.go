package server

import (
	"net/http"
	"dexchan/server/cfg"
	"strconv"
	dexDB "dexchan/server/db"
	"github.com/gorilla/mux"
	"dexchan/server/model"
	"log"
)
type Server struct {
	Config *cfg.C
	Db dexDB.D

	boards []*model.Board
	boardsByShortcode map[string]*model.Board
}

func (s *Server)  staticRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func (s *Server) Start() {
	//fs := http.FileServer(http.Dir(s.Config.StaticDir))
	//http.Handle("/static/", http.StripPrefix("/static", fs))
	//http.HandleFunc("/api/boards", s.getBoards)
	//http.HandleFunc("/", s.staticRoot)

	// didnt want to pull in any deps but argh, need them path args
	router := mux.NewRouter()
	router.HandleFunc("/api/boards", s.getBoards).Methods("GET")
	router.HandleFunc("/api/{board:[a-z]+}", s.getThreads).Methods("GET")
	router.HandleFunc("/api/{board:[a-z]+}/{thread:[0-9]+}", s.getPosts).Methods("GET")
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
