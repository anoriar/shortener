package repository

import (
	"context"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

//go:generate mockgen -source=url_repository_interface.go -destination=mock/url_repository.go -package=mock URLRepositoryInterface
type URLRepositoryInterface interface {
	AddURL(url *entity.URL) error
	FindURLByShortURL(shortURL string) (*entity.URL, error)
	FindURLByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error)
	AddURLBatch(ctx context.Context, urls []entity.URL) error
	Ping(ctx context.Context) error
	Close() error
	DeleteURLBatch(ctx context.Context, shortURLs []string) error
}
