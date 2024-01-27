package config

// Config missing godoc.
type Config struct {
	Host            string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL         string `env:"BASE_URL" json:"base_url"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	AuthSecretKey   string `env:"AUTH_SECRET_KEY"`
	ProfilerHost    string `env:"PROFILER_HOST"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	Config          string `env:"CONFIG"`
}

// NewConfig missing godoc.
func NewConfig() *Config {
	return &Config{
		LogLevel:      "info",
		AuthSecretKey: "secret-key",
	}
}
