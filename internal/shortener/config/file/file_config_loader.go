package file

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfigFromFile missing godoc.
func LoadConfigFromFile(filePath string) (*FileConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileConfig := NewFileConfig()
	err = json.Unmarshal(data, fileConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the file configuration data: %v", err)
	}

	return fileConfig, nil
}
