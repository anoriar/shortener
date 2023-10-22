package reader

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

const (
	fileName            = "/tmp/success_file.json"
	emptyFileName       = "/tmp/empty_file.json"
	syntaxErrorFileName = "/tmp/syntax_error_file.json"
)

func TestUrlFileReader_ReadUrl(t *testing.T) {
	testURL := &entity.URL{
		UUID:        "4e473abf-9ded-4b16-8d20-f0964c88a7b9",
		ShortURL:    "sS9fk2",
		OriginalURL: "https://practicum.yandex.ru/",
	}
	fileData, err := json.Marshal(testURL)
	if err != nil {
		t.Fatal("Error unmarshal json", err)
	}

	successFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal("Error create temporary file", err)
	}
	defer func() {
		successFile.Close()
		os.Remove(fileName)
	}()
	_, err = successFile.Write(fileData)
	if err != nil {
		t.Fatal("Error write file", err)
	}

	emptyFile, err := os.OpenFile(emptyFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal("Error create temporary file", err)
	}
	defer func() {
		emptyFile.Close()
		os.Remove(emptyFileName)
	}()

	syntaxErrorFile, err := os.OpenFile(syntaxErrorFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal("Error create temporary file", err)
	}
	defer func() {
		syntaxErrorFile.Close()
		os.Remove(syntaxErrorFileName)
	}()
	_, err = syntaxErrorFile.Write([]byte("syntax exception"))
	if err != nil {
		t.Fatal("Error write file", err)
	}

	tests := []struct {
		name     string
		filename string
		want     *entity.URL
		wantErr  bool
	}{
		{
			name:     "read success",
			filename: fileName,
			want:     testURL,
			wantErr:  false,
		},
		{
			name:     "read empty",
			filename: emptyFileName,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "read syntax exception",
			filename: syntaxErrorFileName,
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewURLFileReader(tt.filename)
			assert.NoError(t, err)

			got, err := c.ReadURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadURL() exception = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//#MENTOR: есть ли различие между reflect.DeepEqual и просто сравнением объектов через assert.equal(got, tt.want)?
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
