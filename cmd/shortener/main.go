package main

import (
	"github.com/anoriar/shortener/internal/handlers"
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"net/http"
)

func main() {
	run()
}

func handleFunc(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		handlers.NewAddHandler(storage.GetInstance(), util.NewKeyGen()).AddURL(w, req)
	case http.MethodGet:
		handlers.NewGetHandler(storage.GetInstance()).GetURL(w, req)
	default:
		http.Error(w, "Method must be POST or GET", http.StatusBadRequest)
	}
}

func run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleFunc)

	http.ListenAndServe("localhost:8080", mux)
}
