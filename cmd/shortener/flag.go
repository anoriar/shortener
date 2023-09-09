package main

import (
	"flag"
	"github.com/anoriar/shortener/internal/config"
	"github.com/caarlos0/env/v6"
)

func parseFlags(config *config.Config) {
	flag.StringVar(&config.Host, "a", "localhost:8080", "Host")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base url")

	flag.Parse()

	err := env.Parse(&config)
	if err != nil {
		return
	}
}
