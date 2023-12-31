package factory

import (
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

// AddURLBatchResponseFactory missing godoc.
type AddURLBatchResponseFactory struct {
	baseURL string
}

// NewAddURLBatchResponseFactory missing godoc.
func NewAddURLBatchResponseFactory(baseURL string) *AddURLBatchResponseFactory {
	return &AddURLBatchResponseFactory{baseURL: baseURL}
}

// CreateResponse missing godoc.
func (factory *AddURLBatchResponseFactory) CreateResponse(urlsMap map[string]entity.URL, requestURLs []request.AddURLBatchRequestDTO) []response.AddURLBatchResponseDTO {
	responseURLs := make([]response.AddURLBatchResponseDTO, 0, len(urlsMap))
	for _, reqURL := range requestURLs {
		urlEntity, exists := urlsMap[reqURL.CorrelationID]
		if exists {
			responseURLs = append(responseURLs, response.AddURLBatchResponseDTO{
				CorrelationID: reqURL.CorrelationID,
				ShortURL:      factory.baseURL + "/" + urlEntity.ShortURL,
			})
		}
	}
	return responseURLs
}
