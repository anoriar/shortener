package config

type Config struct {
	Host    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func NewConfig() *Config {
	return &Config{}
}
