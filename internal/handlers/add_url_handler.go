package handlers

import (
	"github.com/anoriar/shortener/internal/util"
	"io"
	"net/http"
)

func AddURL(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content type must be text/plain", http.StatusBadRequest)
		return
	}
	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := "http://localhost:8080/" + util.GenerateShortKey()

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
