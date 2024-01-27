package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anoriar/shortener/internal/shortener/config"
)

// LoadAndMergeConfig missing godoc.
func LoadAndMergeConfig(config *config.Config) error {
	if config.Config == "" {
		return nil
	}
	fileConfig, err := loadConfigFromFile(config.Config)
	if err != nil {
		return err
	}
	mergeConfig(*fileConfig, config)
	return nil
}

func mergeConfig(fileConfig FileConfig, config *config.Config) {
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
}

func loadConfigFromFile(filePath string) (*FileConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileConfig := &FileConfig{}
	err = json.Unmarshal(data, fileConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the file configuration data: %v", err)
	}

	return fileConfig, nil
}
