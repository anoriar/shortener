package repository

import "github.com/anoriar/shortener/internal/shortener/entity"

type URLRepositoryInterface interface {
	AddURL(url *entity.Url) (*entity.Url, error)
	FindURLByShortURL(shortURL string) (*entity.Url, error)
}
