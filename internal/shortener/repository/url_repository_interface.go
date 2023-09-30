package repository

import "github.com/anoriar/shortener/internal/shortener/entity"

type URLRepositoryInterface interface {
	AddURL(url string, key string) (*entity.Url, error)
	FindURLByKey(key string) (*entity.Url, error)
}
