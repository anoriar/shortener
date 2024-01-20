package main

import (
	"flag"

	"github.com/anoriar/shortener/internal/shortener/config"
)

func parseFlags(config *config.Config) {
	flag.StringVar(&config.Host, "a", "localhost:8080", "Host")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base url_gen")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url_gen-db.json", "File storage path")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&config.AuthSecretKey, "ask", "secret-key", "Auth secret key")
	flag.StringVar(&config.ProfilerHost, "profiler", "", "Profiler host")
	flag.BoolVar(&config.EnableHTTPS, "s", false, "Enable https")

	flag.Parse()
}
