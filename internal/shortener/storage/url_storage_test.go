package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLStorage_AddURL(t *testing.T) {
	type args struct {
		url string
		key string
	}
	tests := []struct {
		name        string
		existedURLs map[string]string
		args        args
		wantErr     bool
	}{
		{
			name: "add item simple",
			existedURLs: map[string]string{
				"KZXdDY": "https://github.com",
			},
			args: args{
				url: "https://google.com",
				key: "aTgd1u",
			},
			wantErr: false,
		},
		{
			name: "item exists",
			existedURLs: map[string]string{
				"KZXdDY": "https://github.com",
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
			storage := &URLStorage{
				urls: tt.existedURLs,
			}

			err := storage.AddURL(tt.args.url, tt.args.key)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Contains(t, storage.urls, tt.args.key)
		})
	}
}

func TestURLStorage_FindURLByKey(t *testing.T) {

	type want struct {
		url   string
		exist bool
	}

	tests := []struct {
		name        string
		existedURLs map[string]string
		key         string
		want        want
	}{
		{
			name: "item exists",
			existedURLs: map[string]string{
				"KZXdDY": "https://github.com",
			},
			key: "KZXdDY",
			want: want{
				url:   "https://github.com",
				exist: true,
			},
		},
		{
			name: "item not exists",
			existedURLs: map[string]string{
				"KZXdDY": "https://github.com",
			},
			key: "1111",
			want: want{
				url:   "",
				exist: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &URLStorage{
				urls: tt.existedURLs,
			}

			url, exist := storage.FindURLByKey(tt.key)

			assert.Equal(t, tt.want.exist, exist)
			assert.Equal(t, tt.want.url, url)
		})
	}
}
