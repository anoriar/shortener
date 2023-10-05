package repository

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryURLRepository_AddURL(t *testing.T) {
	type args struct {
		url *entity.URL
	}
	tests := []struct {
		name        string
		existedURLs map[string]*entity.URL
		args        args
		wantErr     bool
	}{
		{
			name: "add item simple",
			existedURLs: map[string]*entity.URL{
				"KZXdDY": &entity.URL{
					UUID:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
			args: args{
				&entity.URL{
					UUID:        "4e473abf-9ded-4b16-8d20-f0964c88a7b9",
					ShortURL:    "sS9fk2",
					OriginalURL: "https://practicum.yandex.ru/",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &InMemoryURLRepository{
				urls: tt.existedURLs,
			}

			_, err := repository.AddURL(tt.args.url)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Contains(t, repository.urls, tt.args.url.ShortURL)
		})
	}
}

func TestInMemoryURLRepository_FindURLByKey(t *testing.T) {

	type want struct {
		url *entity.URL
	}

	tests := []struct {
		name        string
		existedURLs map[string]*entity.URL
		key         string
		want        want
	}{
		{
			name: "item exists",
			existedURLs: map[string]*entity.URL{
				"KZXdDY": &entity.URL{
					UUID:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
			key: "KZXdDY",
			want: want{
				url: &entity.URL{
					UUID:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
					ShortURL:    "KZXdDY",
					OriginalURL: "https://github.com",
				},
			},
		},
		{
			name: "item not exists",
			existedURLs: map[string]*entity.URL{
				"KZXdDY": &entity.URL{
					UUID:        "46b8f9d2-b123-4f8e-aabb-f77dd764a00b",
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

			newURL, err := repository.FindURLByShortURL(tt.key)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.url, newURL)
		})
	}
}
