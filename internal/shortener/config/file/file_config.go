package file

// FileConfig missing godoc.
type FileConfig struct {
	Host            string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL         string `env:"BASE_URL" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
}

// NewFileConfig missing godoc.
func NewFileConfig() *FileConfig {
	return &FileConfig{}
}
