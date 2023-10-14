package file

import (
	"context"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/file/internal/reader"
	"github.com/anoriar/shortener/internal/shortener/repository/file/internal/writer"
	"io"
)

type FileURLRepository struct {
	filename string
}

func NewFileURLRepository(filename string) repository.URLRepositoryInterface {
	return &FileURLRepository{
		filename: filename,
	}
}

func (repository *FileURLRepository) AddURL(url *entity.URL) (*entity.URL, error) {

	fileWriter, err := writer.NewURLFileWriter(repository.filename)
	if err != nil {
		return nil, err
	}
	defer fileWriter.Close()

	err = fileWriter.WriteURL(url)
	if err != nil {
		return nil, err
	}
	return url, nil
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
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}

		if url.ShortURL == shortURL {
			return url, nil
		}
	}
}

func (repository *FileURLRepository) AddURLBatch(ctx context.Context, urls []entity.URL) error {
	for _, url := range urls {
		_, err := repository.AddURL(&url)
		if err != nil {
			return err
		}
	}
	return nil
}
