package repository

type URLRepositoryInterface interface {
	AddURL(url string, key string) error
	FindURLByKey(key string) (string, bool)
}
