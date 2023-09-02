package handlers

import (
	"github.com/anoriar/shortener/internal/storage"
	"net/http"
	"strings"
)

func GetURL(w http.ResponseWriter, req *http.Request) {
	shortKey := strings.Trim(req.URL.Path, "/")
	if shortKey == "" {
		http.Error(w, "Short key is empty", http.StatusBadRequest)
	}

	url, exists := storage.GlobalUrlStorage.FindUrlByKey(shortKey)
	if !exists {
		http.Error(w, "Url does not exists", http.StatusBadRequest)
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
