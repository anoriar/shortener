package inmemory

import (
	"context"

	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/entity"
)

// InMemoryURLRepository missing godoc.
type InMemoryURLRepository struct {
	urls map[string]*entity.URL
}

// NewInMemoryURLRepository missing godoc.
func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{urls: make(map[string]*entity.URL)}
}

// Ping missing godoc.
func (repository *InMemoryURLRepository) Ping(ctx context.Context) error {
	return nil
}

// AddURL missing godoc.
func (repository *InMemoryURLRepository) AddURL(url *entity.URL) error {

	repository.urls[url.ShortURL] = url

	return nil
}

// FindURLByOriginalURL missing godoc.
func (repository *InMemoryURLRepository) FindURLByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	return repository.findOneByCondition(func(url entity.URL) bool {
		return url.OriginalURL == originalURL
	})
}

func (repository *InMemoryURLRepository) findOneByCondition(condition func(url entity.URL) bool) (*entity.URL, error) {
	for _, url := range repository.urls {
		if condition(*url) {
			return url, nil
		}
	}
	return nil, nil
}

// FindURLByShortURL missing godoc.
func (repository *InMemoryURLRepository) FindURLByShortURL(key string) (*entity.URL, error) {
	url, exists := repository.urls[key]
	if !exists {
		return nil, nil
	}
	return url, nil
}

// GetURLsByQuery missing godoc.
func (repository *InMemoryURLRepository) GetURLsByQuery(ctx context.Context, urlQuery repository.Query) ([]entity.URL, error) {
	var resultURLs []entity.URL

	for _, url := range repository.urls {
		if len(urlQuery.OriginalURLs) > 0 {
			for _, originalURL := range urlQuery.OriginalURLs {
				if url.OriginalURL == originalURL {
					resultURLs = append(resultURLs, *url)
					continue
				}
			}
		}

		if len(urlQuery.ShortURLs) > 0 {
			for _, shortURL := range urlQuery.ShortURLs {
				if url.ShortURL == shortURL {
					resultURLs = append(resultURLs, *url)
					continue
				}
			}
		}
	}
	return resultURLs, nil
}

// AddURLBatch missing godoc.
func (repository *InMemoryURLRepository) AddURLBatch(ctx context.Context, urls []entity.URL) error {
	for _, url := range urls {
		err := repository.AddURL(&url)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteURLBatch missing godoc.
func (repository *InMemoryURLRepository) DeleteURLBatch(ctx context.Context, shortURLs []string) error {
	for _, shortURL := range shortURLs {
		delete(repository.urls, shortURL)
	}
	return nil
}

// UpdateIsDeletedBatch missing godoc.
func (repository *InMemoryURLRepository) UpdateIsDeletedBatch(ctx context.Context, shortURLs []string, isDeleted bool) error {
	for _, shortURL := range shortURLs {
		if item, ok := repository.urls[shortURL]; ok {
			item.IsDeleted = isDeleted
		}
	}
	return nil
}

// Close missing godoc.
func (repository *InMemoryURLRepository) Close() error {
	return nil
}
