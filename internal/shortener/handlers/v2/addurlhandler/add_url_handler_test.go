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
const successRequestBody = `{"url":"http://localhost:8080/etw73C"}`
const baseURL = "http://localhost:8080"
const successExpectedBody = `{"result":"http://localhost:8080/etw73C"}`

type mockAddHandlerURLRepositoryError struct{}

func (mcr *mockAddHandlerURLRepositoryError) AddURL(url *entity.Url) (*entity.Url, error) {
	return nil, errors.New("test")
}

func (mcr *mockAddHandlerURLRepositoryError) FindURLByShortURL(shortURL string) (*entity.Url, error) {
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
		name          string
		requestBody   string
		urlRepository repository.URLRepositoryInterface
		want          want
	}{
		{
			name:          "success",
			requestBody:   successRequestBody,
			urlRepository: repository.NewInMemoryURLRepository(),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "application/json",
			},
		},
		{
			name:          "not valid url",
			requestBody:   "/dd",
			urlRepository: repository.NewInMemoryURLRepository(),
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "repository error",
			requestBody:   successRequestBody,
			urlRepository: new(mockAddHandlerURLRepositoryError),
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

			NewAddHandler(tt.urlRepository, keyGenMock, baseURL).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
