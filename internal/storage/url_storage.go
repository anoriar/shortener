package storage

import (
	"fmt"
)

type URLStorage struct {
	urls map[string]string
}

func newURLStorage() *URLStorage {
	return &URLStorage{urls: make(map[string]string)}
}

func GetInstance() URLStorageInterface {
	once.Do(func() {
		urlStorageInstance = newURLStorage()
	})
	return urlStorageInstance
}

func (storage *URLStorage) AddURL(url string, key string) error {
	if _, exists := storage.FindURLByKey(key); exists {
		return fmt.Errorf("url with key %v exists", key)
	}
	storage.urls[key] = url

	return nil
}

func (storage *URLStorage) FindURLByKey(key string) (string, bool) {
	url, exists := storage.urls[key]
	return url, exists
}
