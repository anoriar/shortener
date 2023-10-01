package repository

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/file/reader"
	"github.com/anoriar/shortener/internal/shortener/file/writer"
	"io"
)

type FileURLRepository struct {
	filename string
}

func NewFileURLRepository(filename string) URLRepositoryInterface {
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