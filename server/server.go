package server

import (

	"fmt"
	"net/http"
	"dexchan/server/cfg"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Start(c *cfg.C){
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}