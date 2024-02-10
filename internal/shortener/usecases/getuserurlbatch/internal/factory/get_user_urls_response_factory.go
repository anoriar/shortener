package factory

import (
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

// GetUSerURLsResponseFactory missing godoc.
type GetUSerURLsResponseFactory struct {
	baseURL string
}

// NewGetUSerURLsResponseFactory missing godoc.
func NewGetUSerURLsResponseFactory(baseURL string) *GetUSerURLsResponseFactory {
	return &GetUSerURLsResponseFactory{baseURL: baseURL}
}

// CreateResponse missing godoc.
func (factory *GetUSerURLsResponseFactory) CreateResponse(urls []entity.URL) []response.GetUserURLsResponseDTO {
	responseURLs := make([]response.GetUserURLsResponseDTO, 0, len(urls))
	for _, url := range urls {
		responseURLs = append(responseURLs, response.GetUserURLsResponseDTO{
			ShortURL:    factory.baseURL + "/" + url.ShortURL,
			OriginalURL: url.OriginalURL,
		})
	}
	return responseURLs
}
