// Package url Пакет работы с хранилище URL
package url

import (
	"context"

	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

// URLRepositoryInterface Интерфейс работы с хранилищем URL
//
//go:generate mockgen -source=url_repository_interface.go -destination=mock/url_repository.go -package=mock URLRepositoryInterface
type URLRepositoryInterface interface {
	// AddURL Добавляет URL в хранилище
	AddURL(url *entity.URL) error
	// FindURLByShortURL Получает URL по его короткой версии (т.е. HnsSMA)
	FindURLByShortURL(shortURL string) (*entity.URL, error)
	// FindURLByOriginalURL Получает запись URL по его оригинальной версии (т.е. https://www.google1.ru/)
	FindURLByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error)
	// AddURLBatch Добавление нескольких URL
	AddURLBatch(ctx context.Context, urls []entity.URL) error
	// Ping Проверка жизнеспособности хранилища
	Ping(ctx context.Context) error
	// Close Закрывает работу с хранилищем
	Close() error
	// DeleteURLBatch Удаление нескольких URL
	DeleteURLBatch(ctx context.Context, shortURLs []string) error
	// GetURLsByQuery Получение URL по запросу
	GetURLsByQuery(ctx context.Context, urlQuery repository.Query) ([]entity.URL, error)
	// UpdateIsDeletedBatch Проставление флага is_deleted у нескольких URL
	UpdateIsDeletedBatch(ctx context.Context, shortURLs []string, isDeleted bool) error
}
