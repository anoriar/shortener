package addurlhandler

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
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

type mockAddHandlerURLRepositoryError struct{}

func (mcr *mockAddHandlerURLRepositoryError) AddURL(url string, key string) (*entity.Url, error) {
	return nil, errors.New("test")
}

func (mcr *mockAddHandlerURLRepositoryError) FindURLByKey(key string) (*entity.Url, error) {
	return nil, nil
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
		urlStorage  repository.URLRepositoryInterface
		want        want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			urlStorage:  repository.NewInMemoryURLRepository(),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "text/plain",
			},
		},
		{
			name:        "not valid url",
			requestBody: "/dd",
			urlStorage:  repository.NewInMemoryURLRepository(),
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "repository error",
			requestBody: successRequestBody,
			urlStorage:  new(mockAddHandlerURLRepositoryError),
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
