package main

import (
	"flag"
	"github.com/anoriar/shortener/internal/shortener/config"
)

func parseFlags(config *config.Config) {
	flag.StringVar(&config.Host, "a", "localhost:8080", "Host")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base url")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")

	flag.Parse()
}
