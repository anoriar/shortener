package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/anoriar/shortener/internal/shortener/config/file"
)

// LoadConfig missing godoc.
func LoadConfig() (*Config, error) {
	conf := NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("parse env error: %v", err)
	}

	err = loadAndMergeConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("parse config from file error: %v", err)
	}

	return conf, nil
}

// LoadAndMergeConfig missing godoc.
func loadAndMergeConfig(config *Config) error {
	if config.Config == "" {
		return nil
	}
	fileConfig, err := file.LoadConfigFromFile(config.Config)
	if err != nil {
		return err
	}
	mergeConfig(*fileConfig, config)
	return nil
}

func mergeConfig(fileConfig file.FileConfig, config *Config) {
	if config.Host == "" {
		config.Host = fileConfig.Host
	}

	if config.BaseURL == "" {
		config.BaseURL = fileConfig.BaseURL
	}
	if config.DatabaseDSN == "" {
		config.DatabaseDSN = fileConfig.DatabaseDSN
	}
	if config.FileStoragePath == "" {
		config.FileStoragePath = fileConfig.FileStoragePath
	}
	if !config.EnableHTTPS {
		config.EnableHTTPS = fileConfig.EnableHTTPS
	}

	if config.TrustedSubnet == "" {
		config.TrustedSubnet = fileConfig.TrustedSubnet
	}
}

func parseFlags(config *Config) {
	flag.StringVar(&config.Host, "a", "localhost:8080", "Host")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base url_gen")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url_gen-db.json", "File storage path")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&config.AuthSecretKey, "ask", "secret-key", "Auth secret key")
	flag.StringVar(&config.ProfilerHost, "profiler", "", "Profiler host")
	flag.BoolVar(&config.EnableHTTPS, "s", false, "Enable https")
	flag.StringVar(&config.Config, "c", "", "Config from file")
	flag.StringVar(&config.Config, "t", "", "Trusted utilip")

	flag.Parse()
}
