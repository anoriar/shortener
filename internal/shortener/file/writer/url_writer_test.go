package writer

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const fileName = "/tmp/success_file.json"

func TestUrlFileWriter_WriteURL(t *testing.T) {
	type args struct {
		url *entity.Url
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal("Error create temporary file", err)
	}
	defer func() {
		file.Close()
		os.Remove(fileName)
	}()

	tests := []struct {
		name     string
		filename string
		args     args
		wantErr  bool
	}{
		{
			name:     "success write",
			filename: fileName,
			args: args{
				&entity.Url{
					Uuid:        "4e473abf-9ded-4b16-8d20-f0964c88a7b9",
					ShortURL:    "sS9fk2",
					OriginalURL: "https://practicum.yandex.ru/",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewUrlFileWriter(tt.filename)
			assert.NoError(t, err)

			if err := w.WriteURL(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("WriteURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			readFile, err := os.ReadFile(tt.filename)
			if err != nil {
				return
			}

			urlFromFile := &entity.Url{}
			err = json.Unmarshal(readFile, urlFromFile)
			assert.NoError(t, err)

			assert.Equal(t, tt.args.url, urlFromFile)
		})
	}
}
