package config

// TestConfig missing godoc.
type TestConfig struct {
	BaseURL string `env:"BASE_URL"`
}

// NewTestConfig missing godoc.
func NewTestConfig() *TestConfig {
	return &TestConfig{}
}
