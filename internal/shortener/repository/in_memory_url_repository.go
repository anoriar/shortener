package repository

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
)

type InMemoryURLRepository struct {
	urls map[string]*entity.Url
}

func NewInMemoryURLRepository() URLRepositoryInterface {
	return &InMemoryURLRepository{urls: make(map[string]*entity.Url)}
}

func (repository *InMemoryURLRepository) AddURL(url *entity.Url) (*entity.Url, error) {

	repository.urls[url.ShortURL] = url

	return url, nil
}

func (repository *InMemoryURLRepository) FindURLByShortURL(key string) (*entity.Url, error) {
	url, exists := repository.urls[key]
	if !exists {
		return nil, nil
	}
	return url, nil
}
