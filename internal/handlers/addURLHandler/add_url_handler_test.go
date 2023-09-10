package addURLHandler

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
const baseURL = "http://localhost:8080"
const successExpectedBody = baseURL + "/" + expectedShortKey

// TODO MENTOR Как лучше организовывать моки в структуре проекта? Они видны во всем пакете и могут мешать друг другу.
type mockAddHandlerURLStorage struct{}

func (mcr *mockAddHandlerURLStorage) AddURL(url string, key string) error {
	return nil
}

func (mcr *mockAddHandlerURLStorage) FindURLByKey(key string) (string, bool) {
	return "https://github.com", true
}

type mockAddHandlerURLStorageError struct{}

func (mcr *mockAddHandlerURLStorageError) AddURL(url string, key string) error {
	return errors.New("test")
}

func (mcr *mockAddHandlerURLStorageError) FindURLByKey(key string) (string, bool) {
	return "https://github.com", true
}

type mockAddHandlerKeyGen struct{}

func (mock *mockAddHandlerKeyGen) Generate() string {
	return expectedShortKey
}

func TestAddURL(t *testing.T) {

	type want struct {
		status      int
		body        string
		contentType string
	}
	tests := []struct {
		name        string
		requestBody string
		storageMock storage.URLStorageInterface
		want        want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			storageMock: new(mockAddHandlerURLStorage),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "text/plain",
			},
		},
		{
			name:        "not valid url",
			requestBody: "/dd",
			storageMock: new(mockAddHandlerURLStorage),
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "storage error",
			requestBody: successRequestBody,
			storageMock: new(mockAddHandlerURLStorageError),
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

			keyGenMock := new(mockAddHandlerKeyGen)

			NewAddHandler(tt.storageMock, keyGenMock, baseURL).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
