package config

type TestConfig struct {
	BaseURL string `env:"BASE_URL"`
}

func NewTestConfig() *TestConfig {
	return &TestConfig{}
}
