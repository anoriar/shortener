package addurlhandler

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
		urlStorage  storage.URLStorageInterface
		want        want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			urlStorage:  storage.GetInstance(),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "text/plain",
			},
		},
		{
			name:        "not valid url",
			requestBody: "/dd",
			urlStorage:  storage.GetInstance(),
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
