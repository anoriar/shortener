package file

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	filename = "/tmp/file_storage.json"
)

func TestFileURLRepository_AddURL(t *testing.T) {

	defer os.Remove(filename)

	type fields struct {
		filename string
	}
	type args struct {
		url *entity.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.URL
		wantErr bool
	}{
		{
			name: "success add",
			fields: fields{
				filename: filename,
			},
			args: args{
				&entity.URL{
					UUID:        "4e473abf-9ded-4b16-8d20-f0964c88a7b9",
					ShortURL:    "sS9fk2",
					OriginalURL: "https://practicum.yandex.ru/",
				},
			},
			want: &entity.URL{
				UUID:        "4e473abf-9ded-4b16-8d20-f0964c88a7b9",
				ShortURL:    "sS9fk2",
				OriginalURL: "https://practicum.yandex.ru/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewFileURLRepository(filename)
			got, err := repository.AddURL(tt.args.url)
			if tt.wantErr != (err != nil) {
				t.Errorf("AddURL() exception = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileURLRepository_FindURLByShortURL(t *testing.T) {
	defer os.Remove(filename)

	successExistedURLs := []*entity.URL{
		&entity.URL{
			UUID:        "4e473abf-9ded-4b16-8d20-f0964c88a7b9",
			ShortURL:    "sS9fk2",
			OriginalURL: "https://practicum.yandex.ru/",
		},
		&entity.URL{
			UUID:        "b9d1113f-da5f-40d2-b9ef-15a3daf23668",
			ShortURL:    "ge9Yk2",
			OriginalURL: "https://google.com",
		},
		&entity.URL{
			UUID:        "936f1338-c817-4ce4-924a-58bc34f1dd4f",
			ShortURL:    "t8Fhd7",
			OriginalURL: "https://yandex.ru/",
		},
	}
	var fileData []byte
	if len(successExistedURLs) > 0 {

		for _, existedURL := range successExistedURLs {
			data, err := json.Marshal(existedURL)
			if err != nil {
				t.Fatal("Error when marshal json", err)
			}
			fileData = append(fileData, data...)
			fileData = append(fileData, []byte("\n")...)
		}
	}

	type fields struct {
		filename string
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name        string
		existedURLs []*entity.URL
		fileData    []byte
		fields      fields
		args        args
		want        *entity.URL
		wantErr     bool
	}{
		{
			name:     "existed from many urls",
			fileData: fileData,
			fields: fields{
				filename: filename,
			},
			args: args{
				shortURL: "ge9Yk2",
			},
			want: &entity.URL{
				UUID:        "b9d1113f-da5f-40d2-b9ef-15a3daf23668",
				ShortURL:    "ge9Yk2",
				OriginalURL: "https://google.com",
			},
			wantErr: false,
		},
		{
			name:     "not existed from many urls",
			fileData: fileData,
			fields: fields{
				filename: filename,
			},
			args: args{
				shortURL: "fdsfsf",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty file",
			fields: fields{
				filename: filename,
			},
			args: args{
				shortURL: "ge9Yk2",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:     "syntax exception file",
			fileData: []byte{21, 32, 44},
			fields: fields{
				filename: filename,
			},
			args: args{
				shortURL: "ge9Yk2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(filename)
			if len(tt.fileData) > 0 {
				file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					t.Fatal("Error when file create", err)
				}
				_, err = file.Write(tt.fileData)
				if err != nil {
					t.Fatal("Error when write data", err)
				}
			}

			repository := NewFileURLRepository(filename)
			got, err := repository.FindURLByShortURL(tt.args.shortURL)
			if tt.wantErr != (err != nil) {
				t.Errorf("FindURLByShortURL() exception = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equalf(t, tt.want, got, "FindURLByShortURL(%v)", tt.args.shortURL)
		})
	}
}
