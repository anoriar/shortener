package writer

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"os"
)

type UrlFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewUrlFileWriter(filename string) (*UrlFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &UrlFileWriter{file: file, encoder: json.NewEncoder(file)}, nil
}

func (w *UrlFileWriter) WriteURL(url *entity.Url) error {
	err := w.encoder.Encode(url)
	if err != nil {
		return err
	}
	return nil
}

func (w *UrlFileWriter) Close() error {
	return w.file.Close()
}
