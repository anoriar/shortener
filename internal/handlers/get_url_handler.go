package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func GetURL(w http.ResponseWriter, req *http.Request) {
	shortKey := strings.Trim(req.URL.Path, "/")
	if shortKey == "" {
		http.Error(w, "Short key is empty", http.StatusBadRequest)
	}

	fmt.Println(shortKey)

	w.Header().Set("Location", "https://practicum.yandex.ru/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
