package urlgen

import (
	"errors"
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/util"
)

const maxAttempts = 5

// #MENTOR: Как лучше оформлять кастомные ошибки? Делать отдельным типом или глобальной переменной?
var ShortKeyGenerationAttemptsExceededError = errors.New(fmt.Sprintf("max number of attempts for short url generation has been exhausted: %v", maxAttempts))

type ShortURLGenerator struct {
	urlRepository repository.URLRepositoryInterface
	keyGen        util.KeyGenInterface
}

func InitializeShortURLGenerator(urlRepository repository.URLRepositoryInterface) ShortURLGeneratorInterface {
	return NewShortURLGenerator(urlRepository, util.NewKeyGen())
}

func NewShortURLGenerator(urlRepository repository.URLRepositoryInterface, keyGen util.KeyGenInterface) *ShortURLGenerator {
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
	return "", ShortKeyGenerationAttemptsExceededError
}
