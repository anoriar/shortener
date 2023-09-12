package client

import "github.com/anoriar/shortener/cmd/e2e/client/dto"

type ShortenerClientInterface interface {
	AddURL(url string) (*dto.AddResponseDto, error)
	GetURL(key string) (*dto.GetResponseDto, error)
}
