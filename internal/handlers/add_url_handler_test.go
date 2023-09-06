package handlers

import (
	"errors"
	"github.com/anoriar/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const expectedShortKey = "etw73C"
const successRequestBody = "https://github.com"
const baseURL = "http://localhost:8080/"
const successExpectedBody = baseURL + expectedShortKey

// TODO MENTOR Как лучше организовывать моки? Они видны во всем пакете и могут мешать друг другу
type mockURLRepositoryAddHandler struct{}

func (mcr *mockURLRepositoryAddHandler) AddURL(url string, key string) error {
	return nil
}

func (mcr *mockURLRepositoryAddHandler) FindURLByKey(key string) (string, bool) {
	return "https://github.com", true
}

type mockURLRepositoryErrorAddHandler struct{}

func (mcr *mockURLRepositoryErrorAddHandler) AddURL(url string, key string) error {
	return errors.New("test")
}

func (mcr *mockURLRepositoryErrorAddHandler) FindURLByKey(key string) (string, bool) {
	return "https://github.com", true
}

type mockKeyGenAddHandler struct{}

func (mock *mockKeyGenAddHandler) Generate() string {
	return expectedShortKey
}

func TestAddURL(t *testing.T) {

	type want struct {
		status      int
		body        string
		contentType string
	}
	tests := []struct {
		name           string
		requestBody    string
		repositoryMock storage.URLRepositoryInterface
		want           want
	}{
		{
			name:           "success",
			requestBody:    successRequestBody,
			repositoryMock: new(mockURLRepositoryAddHandler),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "text/plain",
			},
		},
		{
			name:           "not valid url",
			requestBody:    "/dd",
			repositoryMock: new(mockURLRepositoryAddHandler),
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:           "repository error",
			requestBody:    successRequestBody,
			repositoryMock: new(mockURLRepositoryErrorAddHandler),
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			keyGenMock := new(mockKeyGenAddHandler)

			NewAddHandler(tt.repositoryMock, keyGenMock).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}