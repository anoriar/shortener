package storage

type UrlStorageInterface interface {
	AddUrl(url string, key string) error
	FindUrlByKey(key string) (string, bool)
}
