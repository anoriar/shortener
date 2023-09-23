package config

type Config struct {
	Host     string `env:"SERVER_ADDRESS"`
	BaseURL  string `env:"BASE_URL"`
	LogLevel string `env:"LOG_LEVEL"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "info",
	}
}
