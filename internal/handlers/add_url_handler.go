package handlers

import (
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"io"
	"net/http"
	neturl "net/url"
)

func AddURL(w http.ResponseWriter, req *http.Request) {
	url, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	parsedUrl, err := neturl.Parse(string(url))
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		http.Error(w, "Not valid URL", http.StatusBadRequest)
		return
	}

	shortKey := util.GenerateShortKey()
	err = storage.GetInstance().AddURL(string(url), shortKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte("http://localhost:8080/" + shortKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
