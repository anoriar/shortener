package urlgen

import (
	"fmt"

	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/util"
)

const maxAttempts = 5

var ErrShortKeyGenerationAttemptsExceeded = fmt.Errorf("max number of attempts for short url generation has been exhausted: %v", maxAttempts)

type ShortURLGenerator struct {
	urlRepository url.URLRepositoryInterface
	keyGen        util.KeyGenInterface
}

func NewShortURLGenerator(urlRepository url.URLRepositoryInterface, keyGen util.KeyGenInterface) *ShortURLGenerator {
	return &ShortURLGenerator{urlRepository: urlRepository, keyGen: keyGen}
}

func (sug *ShortURLGenerator) GenerateShortURL() (string, error) {
	attempt := 0
	for attempt < maxAttempts {
		shortURL := sug.keyGen.Generate()
		url, err := sug.urlRepository.FindURLByShortURL(shortURL)
		if err != nil {
			return "", err
		}
		if url == nil {
			return shortURL, nil
		}
		attempt++
	}

	return "", ErrShortKeyGenerationAttemptsExceeded
}
