package reader

import (
	"encoding/json"
	"os"

	"github.com/anoriar/shortener/internal/shortener/entity"
)

// URLFileReader missing godoc.
type URLFileReader struct {
	file    *os.File
	decoder *json.Decoder
}

// NewURLFileReader missing godoc.
func NewURLFileReader(filename string) (*URLFileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &URLFileReader{file: file, decoder: json.NewDecoder(file)}, nil
}

// ReadURL missing godoc.
func (c *URLFileReader) ReadURL() (*entity.URL, error) {
	event := &entity.URL{}
	err := c.decoder.Decode(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Close missing godoc.
func (c *URLFileReader) Close() error {
	return c.file.Close()
}
