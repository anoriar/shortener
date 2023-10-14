package repository

import (
	"context"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

type URLRepositoryInterface interface {
	AddURL(url *entity.URL) (*entity.URL, error)
	FindURLByShortURL(shortURL string) (*entity.URL, error)
	AddURLBatch(ctx context.Context, urls []entity.URL) error
}
