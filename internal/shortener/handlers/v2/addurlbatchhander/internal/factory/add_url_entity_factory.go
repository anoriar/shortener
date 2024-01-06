package factory

import (
	"github.com/google/uuid"

	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/util"
)

// AddURLEntityFactory missing godoc.
type AddURLEntityFactory struct {
	keyGen util.KeyGenInterface
}

// NewAddURLBatchFactory missing godoc.
func NewAddURLBatchFactory(keyGen util.KeyGenInterface) *AddURLEntityFactory {
	return &AddURLEntityFactory{keyGen: keyGen}
}

// CreateURLsFromBatchRequest missing godoc.
func (factory *AddURLEntityFactory) CreateURLsFromBatchRequest(requestURLs []request.AddURLBatchRequestDTO) map[string]entity.URL {
	urls := make(map[string]entity.URL, len(requestURLs))
	for _, reqURL := range requestURLs {
		urls[reqURL.CorrelationID] = entity.URL{
			UUID:        uuid.NewString(),
			ShortURL:    factory.keyGen.Generate(),
			OriginalURL: reqURL.OriginalURL,
		}
	}
	return urls
}
