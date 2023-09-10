package main

import (
	"flag"
	"github.com/anoriar/shortener/internal/config"
)

func parseFlags(config *config.Config) {
	flag.StringVar(&config.Host, "a", "localhost:8080", "Host")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base url")

	flag.Parse()
}
