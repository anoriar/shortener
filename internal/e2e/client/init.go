package client

import (
	"github.com/anoriar/shortener/internal/e2e/config"
	"net/http"
)

func InitializeShortenerClient(cnf *config.TestConfig) ShortenerClientInterface {
	return &ShortenerClient{
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Disable automatic redirects
				return http.ErrUseLastResponse
			},
		},
		baseURL: cnf.BaseURL,
	}
}
