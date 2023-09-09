package config

type Config struct {
	Host    string
	BaseURL string
}

func NewConfig() *Config {
	return &Config{}
}
