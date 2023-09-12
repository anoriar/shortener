package geturlhandler

import (
	"github.com/anoriar/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const successRedirectLocation = "https://github.com"
const existedKey = "sHde1e"

type mockGetHandlerURLStorageNotExists struct{}

func (mcr *mockGetHandlerURLStorageNotExists) AddURL(url string, key string) error {
	return nil
}

func (mcr *mockGetHandlerURLStorageNotExists) FindURLByKey(key string) (string, bool) {
	return "", false
}

func TestGetHandler_GetURL(t *testing.T) {
	urlStorage := storage.GetInstance()
	err := urlStorage.AddURL(successRedirectLocation, existedKey)
	assert.NoError(t, err)

	type want struct {
		status      int
		contentType string
		location    string
	}
	tests := []struct {
		name       string
		request    string
		urlStorage storage.URLStorageInterface
		want       want
	}{
		{
			name:       "success",
			request:    "/" + existedKey,
			urlStorage: urlStorage,
			want: want{
				status:      http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    successRedirectLocation,
			},
		},
		{
			name:       "empty short key",
			request:    "/",
			urlStorage: urlStorage,
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:       "empty short key",
			request:    "/",
			urlStorage: new(mockGetHandlerURLStorageNotExists),
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			NewGetHandler(tt.urlStorage).GetURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, w.Header().Get("Location"))
			}
		})
	}
}
