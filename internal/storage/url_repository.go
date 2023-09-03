package storage

import (
	"fmt"
)

type URLRepository struct {
	//TODO: Нужно ли в данном случае делать переменную urls публичной? И почему? #MENTOR
	urls map[string]string
}

func newURLStorage(urls map[string]string) *URLRepository {
	return &URLRepository{urls: urls}
}

func (storage *URLRepository) AddURL(url string, key string) error {
	if _, exists := storage.FindURLByKey(key); exists {
		return fmt.Errorf("url with key %v exists", key)
	}
	storage.urls[key] = url

	return nil
}

func (storage *URLRepository) FindURLByKey(key string) (string, bool) {
	url, exists := storage.urls[key]
	return url, exists
}
