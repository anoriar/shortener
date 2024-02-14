package file

// FileConfig missing godoc.
type FileConfig struct {
	Host            string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// NewFileConfig missing godoc.
func NewFileConfig() *FileConfig {
	return &FileConfig{}
}
