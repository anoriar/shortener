package config

// Config missing godoc.
type Config struct {
	Host            string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	AuthSecretKey   string `env:"AUTH_SECRET_KEY"`
	ProfilerHost    string `env:"PROFILER_HOST"`
}

// NewConfig missing godoc.
func NewConfig() *Config {
	return &Config{
		LogLevel:      "info",
		AuthSecretKey: "secret-key",
	}
}
