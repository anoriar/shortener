package handlers

import (
	"github.com/anoriar/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const successRedirectLocation = "https://github.com"

type mockURLRepositoryGetHandler struct{}

func (mcr *mockURLRepositoryGetHandler) AddURL(url string, key string) error {
	return nil
}

func (mcr *mockURLRepositoryGetHandler) FindURLByKey(key string) (string, bool) {
	return successRedirectLocation, true
}

type mockURLRepositoryNotExistsGetHandler struct{}

func (mcr *mockURLRepositoryNotExistsGetHandler) AddURL(url string, key string) error {
	return nil
}

func (mcr *mockURLRepositoryNotExistsGetHandler) FindURLByKey(key string) (string, bool) {
	return "", false
}

func TestGetHandler_GetURL(t *testing.T) {
	type want struct {
		status      int
		contentType string
		location    string
	}
	tests := []struct {
		name           string
		request        string
		repositoryMock storage.URLRepositoryInterface
		want           want
	}{
		{
			name:           "success",
			request:        "/sHde1e",
			repositoryMock: new(mockURLRepositoryGetHandler),
			want: want{
				status:      http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    successRedirectLocation,
			},
		},
		{
			name:           "empty short key",
			request:        "/",
			repositoryMock: new(mockURLRepositoryGetHandler),
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:           "empty short key",
			request:        "/",
			repositoryMock: new(mockURLRepositoryNotExistsGetHandler),
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

			NewGetHandler(tt.repositoryMock).GetURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, w.Header().Get("Location"))
			}
		})
	}
}