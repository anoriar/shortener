package client

import (
	"github.com/anoriar/shortener/internal/e2e/client/dto/response"
)

type ShortenerClientInterface interface {
	AddURL(url string) (*response.AddResponseDto, error)
	AddURLv2(url string) (*response.AddResponseV2Dto, error)
	GetURL(key string) (*response.GetResponseDto, error)
}
