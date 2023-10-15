package addurlhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/mock"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
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

func TestAddURL(t *testing.T) {

	requestDto := request.AddURLRequestDto{URL: requestURL}
	successRequestBody, err := json.Marshal(requestDto)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepositoryMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlGeneratorMock := urlGenMock.NewMockShortURLGeneratorInterface(ctrl)

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
		mockBehaviour func()
		want          want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil).Times(1)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(nil).Times(1)
			},
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "application/json",
			},
		},
		{
			name:        "not valid url body",
			requestBody: []byte("sss"),
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Times(0)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Times(0)
			},
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "repository exception",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(errors.New("exception")).Times(1)
			},
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "conflict",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(repositoryerror.ErrConflict).Times(1)
				urlRepositoryMock.EXPECT().FindURLByOriginalURL(gomock.Any(), requestURL).Return(&entity.URL{
					UUID:        "8fh34uf349f",
					ShortURL:    expectedShortKey,
					OriginalURL: requestURL,
				}, nil).Times(1)
			},
			want: want{
				status:      http.StatusConflict,
				body:        successExpectedBody,
				contentType: "application/json",
			},
		},
		{
			name:        "conflict find by original url error",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(repositoryerror.ErrConflict).Times(1)
				urlRepositoryMock.EXPECT().FindURLByOriginalURL(gomock.Any(), requestURL).Return(nil, errors.New("error")).Times(1)
			},
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockBehaviour()

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			NewAddHandler(urlRepositoryMock, urlGeneratorMock, logger, baseURL).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
