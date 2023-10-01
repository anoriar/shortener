package reader

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"os"
)

type UrlFileReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewUrlFileReader(filename string) (*UrlFileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &UrlFileReader{file: file, decoder: json.NewDecoder(file)}, nil
}

func (c *UrlFileReader) ReadUrl() (*entity.Url, error) {
	event := &entity.Url{}
	//#MENTOR - есть ли смысл читать файл чанками? например по 20 строк и разом декодировать их? Или по производительности будет одинаково?
	// Насколько я понимаю, чтобы прочитать файл чанками, все равно придется идти сканнером по строкам. Тут только вопрос к декодированию в джсон
	// Процедура декодирования json массива строк быстрее чем по одной строке?
	err := c.decoder.Decode(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (c *UrlFileReader) Close() error {
	return c.file.Close()
}
