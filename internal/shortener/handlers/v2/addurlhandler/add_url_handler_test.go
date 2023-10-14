package addurlhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/mock"
	urlGenMock "github.com/anoriar/shortener/internal/shortener/services/url_gen/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const expectedShortKey = "etw73C"
const requestURL = "https://practicum.yandex.ru"
const baseURL = "http://localhost:8080"
const successExpectedBody = `{"result":"http://localhost:8080/etw73C"}`

type mockAddHandlerURLRepositoryError struct{}

func (mcr *mockAddHandlerURLRepositoryError) AddURL(url *entity.URL) (*entity.URL, error) {
	return nil, errors.New("test")
}

func (mcr *mockAddHandlerURLRepositoryError) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	return nil, nil
}

type mockAddHandlerShortURLGen struct{}

func (mock *mockAddHandlerShortURLGen) GenerateShortURL() (string, error) {
	return expectedShortKey, nil
}

func TestAddURL(t *testing.T) {

	requestDto := request.AddURLRequestDto{URL: requestURL}
	successRequestBody, err := json.Marshal(requestDto)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepoSuccessMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlRepoSuccessMock.EXPECT().AddURL(gomock.Any()).Return(nil)

	urlRepoNotCallsMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlRepoNotCallsMock.EXPECT().AddURL(gomock.Any()).Return(nil).MaxTimes(0)

	urlRepoErrorMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlRepoErrorMock.EXPECT().AddURL(gomock.Any()).Return(errors.New("exception")).MinTimes(1)

	logger, err := logger.Initialize("info")
	require.NoError(t, err)

	type want struct {
		status      int
		body        string
		contentType string
	}
	tests := []struct {
		name          string
		requestBody   []byte
		urlRepository repository.URLRepositoryInterface
		want          want
	}{
		{
			name:          "success",
			requestBody:   successRequestBody,
			urlRepository: urlRepoSuccessMock,
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "application/json",
			},
		},
		{
			name:          "not valid url_gen",
			requestBody:   []byte("sss"),
			urlRepository: urlRepoNotCallsMock,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "repository exception",
			requestBody:   successRequestBody,
			urlRepository: urlRepoErrorMock,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			urlGen := urlGenMock.NewMockShortURLGeneratorInterface(ctrl)
			urlGen.EXPECT().GenerateShortURL().Return(expectedShortKey, nil).AnyTimes()

			NewAddHandler(tt.urlRepository, urlGen, logger, baseURL).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
