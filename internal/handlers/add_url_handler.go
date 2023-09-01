package handlers

import (
	"github.com/anoriar/shortener/internal/util"
	"io"
	"net/http"
)

func AddUrl(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content type must be text/plain", http.StatusBadRequest)
		return
	}
	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := "http://localhost:8080/" + util.GenerateShortKey()

	_, err = w.Write([]byte(shortUrl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
