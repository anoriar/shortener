package writer

import (
	"encoding/json"
	"os"

	"github.com/anoriar/shortener/internal/shortener/entity"
)

type URLFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewURLFileWriter(filename string) (*URLFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &URLFileWriter{file: file, encoder: json.NewEncoder(file)}, nil
}

func NewURLFileEmptyWriter(filename string) (*URLFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &URLFileWriter{file: file, encoder: json.NewEncoder(file)}, nil
}

func (w *URLFileWriter) WriteURL(url *entity.URL) error {
	err := w.encoder.Encode(url)
	if err != nil {
		return err
	}
	return nil
}

func (w *URLFileWriter) Close() error {
	return w.file.Close()
}
