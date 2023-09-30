package repository

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryURLRepository_AddURL(t *testing.T) {
	type args struct {
		url string
		key string
	}
	tests := []struct {
		name        string
		existedURLs map[string]*entity.Url
		args        args
		wantErr     bool
	}{
		{
			name: "add item simple",
			existedURLs: map[string]*entity.Url{
				"KZXdDY": &entity.Url{
					Uuid:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
			args: args{
				url: "https://google.com",
				key: "aTgd1u",
			},
			wantErr: false,
		},
		{
			name: "item exists",
			existedURLs: map[string]*entity.Url{
				"KZXdDY": &entity.Url{
					Uuid:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
			args: args{
				url: "https://google.com",
				key: "KZXdDY",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &InMemoryURLRepository{
				urls: tt.existedURLs,
			}

			_, err := repository.AddURL(tt.args.url, tt.args.key)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Contains(t, repository.urls, tt.args.key)
		})
	}
}

func TestInMemoryURLRepository_FindURLByKey(t *testing.T) {

	type want struct {
		url *entity.Url
	}

	tests := []struct {
		name        string
		existedURLs map[string]*entity.Url
		key         string
		want        want
	}{
		{
			name: "item exists",
			existedURLs: map[string]*entity.Url{
				"KZXdDY": &entity.Url{
					Uuid:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
			key: "KZXdDY",
			want: want{
				url: &entity.Url{
					Uuid:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
		},
		{
			name: "item not exists",
			existedURLs: map[string]*entity.Url{
				"KZXdDY": &entity.Url{
					Uuid:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
			key: "1111",
			want: want{
				url: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &InMemoryURLRepository{
				urls: tt.existedURLs,
			}

			newURL, err := repository.FindURLByKey(tt.key)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.url, newURL)
		})
	}
}
