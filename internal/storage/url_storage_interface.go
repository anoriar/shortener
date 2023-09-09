package storage

type URLStorageInterface interface {
	AddURL(url string, key string) error
	FindURLByKey(key string) (string, bool)
}
