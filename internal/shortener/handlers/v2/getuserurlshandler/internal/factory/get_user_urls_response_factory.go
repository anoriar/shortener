package factory

import (
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

type GetUSerURLsResponseFactory struct {
	baseURL string
}

func NewGetUSerURLsResponseFactory(baseURL string) *GetUSerURLsResponseFactory {
	return &GetUSerURLsResponseFactory{baseURL: baseURL}
}

func (factory *GetUSerURLsResponseFactory) CreateResponse(urls []entity.URL) []response.GetUserURLsResponseDTO {
	var responseURLs []response.GetUserURLsResponseDTO
	for _, url := range urls {
		responseURLs = append(responseURLs, response.GetUserURLsResponseDTO{
			ShortURL:    factory.baseURL + "/" + url.ShortURL,
			OriginalURL: url.OriginalURL,
		})
	}
	return responseURLs
}
