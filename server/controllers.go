package server

//
//import (
//	"net/http"
//	"io"
//	"encoding/json"
//	"github.com/gorilla/mux"
//	"fmt"
//	"dexchan/server/model"
//	"strconv"
//)
//
//func writeJson(w http.ResponseWriter, r *http.Request, v interface{}) {
//	w.Header().Set("Content-Type", "application/json")
//	result, _ := json.Marshal(v)
//	io.WriteString(w, string(result))
//}
//
//func (s *Server) getBoards(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	writeJson(w, r, s.boards)
//}
//
//func (s *Server) tryGetBoard(w http.ResponseWriter, r *http.Request) *model.Board {
//	v := mux.Vars(r)
//
//	bn := v["board"]
//	found, ok := s.boardsByShortcode[bn]
//
//	if !ok {
//		http.Error(w,
//			fmt.Sprintf("bad request: board %s doesnt exist", bn), http.StatusBadRequest)
//		return nil
//	}
//
//	return found
//}
//func (s *Server) getThreads(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	board := s.tryGetBoard(w, r)
//	if board == nil {
//		return
//	}
//	tx, err := s.Db.GetThreads(board.ID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	writeJson(w, r, tx)
//}
//
//func (s *Server) getPosts(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	board := s.tryGetBoard(w, r)
//	if board == nil {
//		return
//	}
//
//	v := mux.Vars(r)
//	threadIDstr := v["thread"]
//	// regex matched, will always be good
//	threadID, _ := strconv.Atoi(threadIDstr)
//	posts, err := s.Db.GetPosts(int64(threadID))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	writeJson(w, r, posts)
//}
//
//func (s *Server) createThread(w http.ResponseWriter, r *http.Request) {
//	board := s.tryGetBoard(w, r)
//	if board == nil {
//		return
//	}
//
//	//s.Db.Creat
//}
//
//func (s *Server) createPost(w http.ResponseWriter, r *http.Request) {
//	board := s.tryGetBoard(w, r)
//	if board == nil {
//		return
//	}
//
//	//v := mux.Vars(r)
//	//threadIDs := v["thread"]
//}
