package repository

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
)

type InMemoryURLRepository struct {
	urls map[string]*entity.URL
}

func NewInMemoryURLRepository() URLRepositoryInterface {
	return &InMemoryURLRepository{urls: make(map[string]*entity.URL)}
}

func (repository *InMemoryURLRepository) AddURL(url *entity.URL) (*entity.URL, error) {

	repository.urls[url.ShortURL] = url

	return url, nil
}

func (repository *InMemoryURLRepository) FindURLByShortURL(key string) (*entity.URL, error) {
	url, exists := repository.urls[key]
	if !exists {
		return nil, nil
	}
	return url, nil
}
