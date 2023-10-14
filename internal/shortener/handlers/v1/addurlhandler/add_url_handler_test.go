package addurlhandler

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/mock"
	urlGenMock "github.com/anoriar/shortener/internal/shortener/services/url_gen/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const expectedShortKey = "etw73C"
const successRequestBody = "https://github.com"
const baseURL = "http://localhost:8080"
const successExpectedBody = baseURL + "/" + expectedShortKey

func TestAddURL(t *testing.T) {
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
		name        string
		requestBody string
		urlStorage  repository.URLRepositoryInterface
		want        want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			urlStorage:  urlRepoSuccessMock,
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "text/plain",
			},
		},
		{
			name:        "not valid url_gen",
			requestBody: "/dd",
			urlStorage:  urlRepoNotCallsMock,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "repository exception",
			requestBody: successRequestBody,
			urlStorage:  urlRepoErrorMock,
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

			urlGen := urlGenMock.NewMockShortURLGeneratorInterface(ctrl)
			urlGen.EXPECT().GenerateShortURL().Return(expectedShortKey, nil).AnyTimes()

			NewAddHandler(tt.urlStorage, urlGen, logger, baseURL).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
