package client

import (
	"github.com/anoriar/shortener/internal/e2e/client/dto/response"
)

type ShortenerClientInterface interface {
	AddURL(url string) (*response.AddResponseDto, error)
	AddURLv2(url string) (*response.AddResponseV2Dto, error)
	AddURLv2WithCompress(url string, contentType string) (*response.AddResponseV2EncodingDto, error)
	GetURL(key string) (*response.GetResponseDto, error)
	Ping() (*response.PingResponseDto, error)
}
