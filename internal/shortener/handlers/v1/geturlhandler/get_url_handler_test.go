package geturlhandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/repository/url/mock"
)

const successRedirectLocation = "https://github.com"
const existedKey = "sHde1e"
const notExistedKey = "sdJ2f3"

type mockGetHandlerURLRepositoryNotExists struct{}

func (mcr *mockGetHandlerURLRepositoryNotExists) AddURL(url *entity.URL) (*entity.URL, error) {
	return nil, nil
}

func (mcr *mockGetHandlerURLRepositoryNotExists) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	return nil, nil
}

func TestGetHandler_GetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepositoryMock := mock.NewMockURLRepositoryInterface(ctrl)

	logger, err := logger.Initialize("info")
	require.NoError(t, err)

	type want struct {
		status      int
		contentType string
		location    string
	}
	tests := []struct {
		name          string
		request       string
		mockBehaviour func()
		urlRepository url.URLRepositoryInterface
		want          want
	}{
		{
			name:    "success",
			request: "/" + existedKey,
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(&entity.URL{
					UUID:        "b9d1113f-da5f-40d2-b9ef-15a3daf23668",
					ShortURL:    existedKey,
					OriginalURL: successRedirectLocation,
				}, nil)
			},
			want: want{
				status:      http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    successRedirectLocation,
			},
		},
		{
			name:    "empty short key",
			request: "/",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Times(0)
			},
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:    "not existed short key",
			request: "/" + notExistedKey,
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(nil, nil)
			},
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:    "exception when fetching",
			request: "/" + notExistedKey,
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(nil, errors.New("exception")).Times(1)
			},
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockBehaviour()

			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			NewGetHandler(urlRepositoryMock, logger).GetURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, w.Header().Get("Location"))
			}
		})
	}
}
