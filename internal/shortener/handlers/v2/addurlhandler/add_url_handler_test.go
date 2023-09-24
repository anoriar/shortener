package addurlhandler

import (
	"errors"
	storage2 "github.com/anoriar/shortener/internal/shortener/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const expectedShortKey = "etw73C"
const successRequestBody = `{"url":"http://localhost:8080/etw73C"}`
const baseURL = "http://localhost:8080"
const successExpectedBody = `{"result":"http://localhost:8080/etw73C"}`

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
		urlStorage  storage2.URLStorageInterface
		want        want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			urlStorage:  storage2.NewURLStorage(),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "application/json",
			},
		},
		{
			name:        "not valid url",
			requestBody: "/dd",
			urlStorage:  storage2.NewURLStorage(),
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "storage error",
			requestBody: successRequestBody,
			urlStorage:  new(mockAddHandlerURLStorageError),
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

			NewAddHandler(tt.urlStorage, keyGenMock, baseURL).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
