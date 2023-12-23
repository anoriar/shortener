package client

import (
	"github.com/anoriar/shortener/internal/e2e/client/dto/request"
	"github.com/anoriar/shortener/internal/e2e/client/dto/response"
)

// ShortenerClientInterface missing godoc.
type ShortenerClientInterface interface {
	AddURL(requestDto request.AddURLRequestDto) (*response.AddResponseDto, error)
	AddURLv2(requestDto request.AddURLRequestDto) (*response.AddResponseV2Dto, error)
	AddURLv2WithCompress(requestDto request.AddURLRequestDto, contentType string) (*response.AddResponseV2EncodingDto, error)
	GetURL(requestDto request.GetURLRequestDto) (*response.GetResponseDto, error)
	Ping() (response.PingResponseDto, error)
	DeleteURLBatch(shortURLs []string) (*response.DeleteURLBatchResponseDto, error)
	AddURLBatch(requestDto request.AddURLBatchRequestDTO) (*response.AddURLBatchResponseDto, error)
	GetUserURLs(requestDto request.AuthRequest) (*response.GetUserURLsResponseDto, error)
	DeleteUserURLs(requestDto request.DeleteUserURLsRequestDto) (*response.DeleteUserURLsResponseDto, error)
}
