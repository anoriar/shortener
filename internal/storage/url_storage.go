package storage

import (
	"fmt"
)

type UrlStorage struct {
	urls map[string]string
}

func NewUrlStorage() *UrlStorage {
	return &UrlStorage{urls: make(map[string]string)}
}

func (storage *UrlStorage) AddUrl(url string, key string) error {
	if _, exists := storage.FindUrlByKey(key); exists {
		return fmt.Errorf("url with key %v exists", key)
	}
	storage.urls[key] = url

	return nil
}

func (storage *UrlStorage) FindUrlByKey(key string) (string, bool) {
	url, exists := storage.urls[key]
	return url, exists
}
