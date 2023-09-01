package main

import (
	"github.com/anoriar/shortener/internal/handlers"
	"net/http"
)

func main() {
	run()
}

func handleFunc(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		handlers.AddURL(w, req)
	case http.MethodGet:
		handlers.GetURL(w, req)
	default:
		http.Error(w, "Method must be POST or GET", http.StatusBadRequest)
	}
}

func run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleFunc)

	http.ListenAndServe("localhost:8080", mux)
}
