package file

import (
	"context"
	"errors"
	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/url/file/internal/reader"
	"github.com/anoriar/shortener/internal/shortener/repository/url/file/internal/writer"
	"io"
)

type FileURLRepository struct {
	filename string
}

func NewFileURLRepository(filename string) *FileURLRepository {
	return &FileURLRepository{
		filename: filename,
	}
}

func (repository *FileURLRepository) Ping(ctx context.Context) error {
	return nil
}

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

func (repository *FileURLRepository) DeleteURLBatch(ctx context.Context, shortURLs []string) error {
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

	//Удаляем лишние
	for _, shortURL := range shortURLs {
		delete(fileURLs, shortURL)
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

func (repository *FileURLRepository) Close() error {
	return nil
}
