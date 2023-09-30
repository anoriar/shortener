package geturlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const successRedirectLocation = "https://github.com"
const existedKey = "sHde1e"
const notExistedKey = "sdJ2f3"

type mockGetHandlerURLRepositoryNotExists struct{}

func (mcr *mockGetHandlerURLRepositoryNotExists) AddURL(url string, key string) (*entity.Url, error) {
	return nil, nil
}

func (mcr *mockGetHandlerURLRepositoryNotExists) FindURLByKey(key string) (*entity.Url, error) {
	return nil, nil
}

func TestGetHandler_GetURL(t *testing.T) {
	urlRepository := repository.NewInMemoryURLRepository()
	_, err := urlRepository.AddURL(successRedirectLocation, existedKey)
	assert.NoError(t, err)

	type want struct {
		status      int
		contentType string
		location    string
	}
	tests := []struct {
		name          string
		request       string
		urlRepository repository.URLRepositoryInterface
		want          want
	}{
		{
			name:          "success",
			request:       "/" + existedKey,
			urlRepository: urlRepository,
			want: want{
				status:      http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    successRedirectLocation,
			},
		},
		{
			name:          "empty short key",
			request:       "/",
			urlRepository: urlRepository,
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:          "not existed short key",
			request:       "/" + notExistedKey,
			urlRepository: urlRepository,
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:          "error when fetching",
			request:       "/" + notExistedKey,
			urlRepository: new(mockGetHandlerURLRepositoryNotExists),
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

			NewGetHandler(tt.urlRepository).GetURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, w.Header().Get("Location"))
			}
		})
	}
}
