package file

import (
	"context"
	"errors"
	"io"

	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/url/file/internal/reader"
	"github.com/anoriar/shortener/internal/shortener/repository/url/file/internal/writer"
)

// FileURLRepository missing godoc.
type FileURLRepository struct {
	filename string
}

// NewFileURLRepository missing godoc.
func NewFileURLRepository(filename string) *FileURLRepository {
	return &FileURLRepository{
		filename: filename,
	}
}

// Ping missing godoc.
func (repository *FileURLRepository) Ping(ctx context.Context) error {
	return nil
}

// AddURL missing godoc.
func (repository *FileURLRepository) AddURL(url *entity.URL) error {

	fileWriter, err := writer.NewURLFileWriter(repository.filename)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	err = fileWriter.WriteURL(url)
	if err != nil {
		return err
	}
	return nil
}

// FindURLByShortURL missing godoc.
func (repository *FileURLRepository) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	fileReader, err := reader.NewURLFileReader(repository.filename)
	if err != nil {
		return nil, err
	}

	defer fileReader.Close()

	for {
		url, err := fileReader.ReadURL()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}

		if url.ShortURL == shortURL {
			return url, nil
		}
	}
}

// GetURLsByQuery missing godoc.
func (repository *FileURLRepository) GetURLsByQuery(ctx context.Context, urlQuery repository.Query) ([]entity.URL, error) {
	var resultURLs []entity.URL

	fileReader, err := reader.NewURLFileReader(repository.filename)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	for {
		url, err := fileReader.ReadURL()
		if err != nil {
			if errors.Is(io.EOF, err) {
				break
			}
			return nil, err
		}

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

// FindURLByOriginalURL missing godoc.
func (repository *FileURLRepository) FindURLByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	return repository.findOneByCondition(func(url entity.URL) bool {
		return url.OriginalURL == originalURL
	})
}

func (repository *FileURLRepository) findOneByCondition(condition func(url entity.URL) bool) (*entity.URL, error) {
	fileReader, err := reader.NewURLFileReader(repository.filename)
	if err != nil {
		return nil, err
	}

	defer fileReader.Close()

	for {
		url, err := fileReader.ReadURL()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}

		if condition(*url) {
			return url, nil
		}
	}
}

// AddURLBatch missing godoc.
func (repository *FileURLRepository) AddURLBatch(ctx context.Context, urls []entity.URL) error {
	fileWriter, err := writer.NewURLFileWriter(repository.filename)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	for _, url := range urls {
		err = fileWriter.WriteURL(&url)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteURLBatch missing godoc.
func (repository *FileURLRepository) DeleteURLBatch(ctx context.Context, shortURLs []string) error {
	return repository.rewriteFile(ctx, func(fileURLs map[string]*entity.URL) error {
		for _, shortURL := range shortURLs {
			delete(fileURLs, shortURL)
		}
		return nil
	})
}

// UpdateIsDeletedBatch missing godoc.
func (repository *FileURLRepository) UpdateIsDeletedBatch(ctx context.Context, shortURLs []string, isDeleted bool) error {
	return repository.rewriteFile(ctx, func(fileURLs map[string]*entity.URL) error {
		for _, shortURL := range shortURLs {
			if item, ok := fileURLs[shortURL]; ok {
				item.IsDeleted = isDeleted
			}
		}
		return nil
	})
}

func (repository *FileURLRepository) rewriteFile(ctx context.Context, callback func(fileURLs map[string]*entity.URL) error) error {
	fileReader, err := reader.NewURLFileReader(repository.filename)
	if err != nil {
		return nil
	}
	fileURLs := make(map[string]*entity.URL)

	defer fileReader.Close()

	//Считываем все данные с файла
	for {
		url, err := fileReader.ReadURL()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		fileURLs[url.ShortURL] = url
	}

	err = callback(fileURLs)
	if err != nil {
		return err
	}

	//Перезаписываем файл заново
	fileWriter, err := writer.NewURLFileEmptyWriter(repository.filename)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	for _, url := range fileURLs {
		err = fileWriter.WriteURL(url)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close missing godoc.
func (repository *FileURLRepository) Close() error {
	return nil
}
