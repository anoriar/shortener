package client

import (
	dto2 "github.com/anoriar/shortener/internal/e2e/client/dto"
)

type ShortenerClientInterface interface {
	AddURL(url string) (*dto2.AddResponseDto, error)
	GetURL(key string) (*dto2.GetResponseDto, error)
}
