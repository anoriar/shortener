package repository

import (
	"fmt"
)

type InMemoryURLRepository struct {
	urls map[string]string
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{urls: make(map[string]string)}
}

func (repository *InMemoryURLRepository) AddURL(url string, key string) error {
	if _, exists := repository.FindURLByKey(key); exists {
		return fmt.Errorf("url with key %v exists", key)
	}
	repository.urls[key] = url

	return nil
}

func (repository *InMemoryURLRepository) FindURLByKey(key string) (string, bool) {
	url, exists := repository.urls[key]
	return url, exists
}
