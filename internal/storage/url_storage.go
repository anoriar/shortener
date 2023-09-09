package storage

import (
	"fmt"
	"sync"
)

var urlStorage URLStorageInterface
var urlStorageSyncOnce sync.Once

type URLStorage struct {
	urls map[string]string
}

func newURLStorage(urls map[string]string) *URLStorage {
	return &URLStorage{urls: urls}
}

func GetInstance() URLStorageInterface {
	urlStorageSyncOnce.Do(func() {
		urlStorage = newURLStorage(make(map[string]string))
	})
	return urlStorage
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
