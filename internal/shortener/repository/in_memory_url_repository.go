package repository

import (
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/google/uuid"
)

type InMemoryURLRepository struct {
	urls map[string]*entity.Url
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{urls: make(map[string]*entity.Url)}
}

func (repository *InMemoryURLRepository) AddURL(url string, key string) (*entity.Url, error) {
	entityURL, err := repository.FindURLByKey(key)
	if err != nil {
		return nil, err
	}
	if entityURL != nil {
		return nil, fmt.Errorf("url with key %v exists", key)
	}

	newURL := &entity.Url{
		Uuid:        uuid.NewString(),
		ShortURL:    key,
		OriginalURL: url,
	}
	repository.urls[key] = newURL

	return newURL, nil
}

func (repository *InMemoryURLRepository) FindURLByKey(key string) (*entity.Url, error) {
	url, exists := repository.urls[key]
	if !exists {
		return nil, nil
	}
	return url, nil
}
