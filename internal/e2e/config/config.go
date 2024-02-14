package config

// TestConfig missing godoc.
type TestConfig struct {
	BaseURL    string `env:"BASE_URL"`
	ServerAddr string `env:"SERVER_ADDR"`
}

// NewTestConfig missing godoc.
func NewTestConfig() *TestConfig {
	return &TestConfig{}
}
