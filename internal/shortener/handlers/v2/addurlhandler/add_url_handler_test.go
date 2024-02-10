package addurlhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anoriar/shortener/internal/shortener/usecases"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
	"github.com/anoriar/shortener/internal/shortener/repository/url/mock"
	urlGenMock "github.com/anoriar/shortener/internal/shortener/services/url_gen/mock"
	mock2 "github.com/anoriar/shortener/internal/shortener/services/user/mock"
)

const expectedShortKey = "etw73C"
const requestURL = "https://practicum.yandex.ru"
const baseURL = "http://localhost:8080"
const successExpectedBody = `{"result":"http://localhost:8080/etw73C"}`
const userID = "6daaf660-a160-4a5c-b99d-faca42c01ef6"

func TestAddURL(t *testing.T) {

	requestDto := request.AddURLRequestDto{URL: requestURL}
	successRequestBody, err := json.Marshal(requestDto)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepositoryMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlGeneratorMock := urlGenMock.NewMockShortURLGeneratorInterface(ctrl)
	userServiceMock := mock2.NewMockUserServiceInterface(ctrl)

	ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userID)

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
		ctx           context.Context
		want          want
	}{
		{
			name:        "success",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil).Times(1)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(nil).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(1)
			},
			ctx: ctxWithUser,
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
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			ctx: ctxWithUser,
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
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			ctx: ctxWithUser,
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
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			ctx: ctxWithUser,
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
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			ctx: ctxWithUser,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "success and empty user context",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil).Times(1)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(nil).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			ctx: context.WithValue(context.Background(), context2.UserIDContextKey, ""),
			want: want{
				status:      http.StatusCreated,
				body:        successExpectedBody,
				contentType: "application/json",
			},
		},
		{
			name:        "add short url to user error",
			requestBody: successRequestBody,
			mockBehaviour: func() {
				urlGeneratorMock.EXPECT().GenerateShortURL().Return(expectedShortKey, nil).Times(1)
				urlRepositoryMock.EXPECT().AddURL(gomock.Any()).Return(nil).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)
			},
			ctx: ctxWithUser,
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
			r = r.WithContext(tt.ctx)
			w := httptest.NewRecorder()

			NewAddHandler(logger, usecases.NewAddURLService(urlRepositoryMock, urlGeneratorMock, userServiceMock, logger, baseURL)).AddURL(w, r)

			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
