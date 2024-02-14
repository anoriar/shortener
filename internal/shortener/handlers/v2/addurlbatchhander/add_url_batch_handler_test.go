package addurlbatchhander

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url/mock"
	mock2 "github.com/anoriar/shortener/internal/shortener/services/user/mock"
	utilMock "github.com/anoriar/shortener/internal/shortener/util/mock"
)

const (
	shortKey1      = "4tH3FG"
	shortKey2      = "G7f6V19"
	shortKey3      = "m31Bfgd"
	baseURL        = "http/localhost"
	originalURL1   = "https://practicum.yandex.ru"
	correlationID1 = "g0fsdf9fj"
	originalURL2   = "https://practicum2.yandex.ru"
	correlationID2 = "ngfdsf3"
	originalURL3   = "https://practicum3.yandex.ru"
	correlationID3 = "by4564trg"
	userID         = "6daaf660-a160-4a5c-b99d-faca42c01ef6"
)

func TestAddURLBatchHandler_AddURLBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	successRequestBody, err := json.Marshal([]request.AddURLBatchRequestDTO{
		{
			CorrelationID: correlationID1,
			OriginalURL:   originalURL1,
		},
		{
			CorrelationID: correlationID2,
			OriginalURL:   originalURL2,
		},
		{
			CorrelationID: correlationID3,
			OriginalURL:   originalURL3,
		},
	})
	require.NoError(t, err)

	successResponseBody, err := json.Marshal([]response.AddURLBatchResponseDTO{
		{
			CorrelationID: correlationID1,
			ShortURL:      baseURL + "/" + shortKey1,
		},
		{
			CorrelationID: correlationID2,
			ShortURL:      baseURL + "/" + shortKey2,
		},
		{
			CorrelationID: correlationID3,
			ShortURL:      baseURL + "/" + shortKey3,
		},
	})
	require.NoError(t, err)

	notValidURLRequestBody, err := json.Marshal([]request.AddURLBatchRequestDTO{
		{
			CorrelationID: correlationID1,
			OriginalURL:   "fdsfsdf",
		},
	})
	require.NoError(t, err)

	reqBodyWithDuplicates, err := json.Marshal([]request.AddURLBatchRequestDTO{
		{
			CorrelationID: correlationID1,
			OriginalURL:   originalURL1,
		},
		{
			CorrelationID: correlationID1,
			OriginalURL:   originalURL2,
		},
	})
	require.NoError(t, err)

	keyGenMock := utilMock.NewMockKeyGenInterface(ctrl)
	logger, err := logger.Initialize("debug")
	require.NoError(t, err)

	urlRepositoryMock := mock.NewMockURLRepositoryInterface(ctrl)
	userServiceMock := mock2.NewMockUserServiceInterface(ctrl)

	ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userID)

	type args struct {
		requestBody []byte
	}
	type want struct {
		status      int
		body        string
		contentType string
	}
	tests := []struct {
		name          string
		mockBehaviour func()
		args          args
		ctx           context.Context
		want          want
	}{
		{
			name: "success",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey1).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey2).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey3).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				requestBody: successRequestBody,
			},
			ctx: ctxWithUser,
			want: want{
				status:      http.StatusCreated,
				body:        string(successResponseBody),
				contentType: "application/json",
			},
		},
		{
			name: "not valid body",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Times(0)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			args: args{
				requestBody: []byte("sss"),
			},
			ctx: ctxWithUser,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "not valid url",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Times(0)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			args: args{
				requestBody: notValidURLRequestBody,
			},
			ctx: ctxWithUser,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "has duplicates",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Times(0)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			args: args{
				requestBody: reqBodyWithDuplicates,
			},
			ctx: ctxWithUser,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "add repo error",
			mockBehaviour: func() {
				keyGenMock.EXPECT().Generate().Return(shortKey1).Times(3)
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Return(errors.New("exception")).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(nil).Times(0)
			},
			args: args{
				requestBody: successRequestBody,
			},
			ctx: ctxWithUser,
			want: want{
				status:      http.StatusBadRequest,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "success and empty user context",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey1).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey2).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey3).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			args: args{
				requestBody: successRequestBody,
			},
			ctx: context.WithValue(context.Background(), context2.UserIDContextKey, ""),
			want: want{
				status:      http.StatusCreated,
				body:        string(successResponseBody),
				contentType: "application/json",
			},
		},
		{
			name: "add short url to user error",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().AddURLBatch(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey1).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey2).Times(1)
				keyGenMock.EXPECT().Generate().Return(shortKey3).Times(1)
				userServiceMock.EXPECT().AddShortURLsToUser(userID, gomock.Any()).Return(errors.New("error")).Times(1)
			},
			args: args{
				requestBody: successRequestBody,
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
			r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(tt.args.requestBody))
			r = r.WithContext(tt.ctx)
			w := httptest.NewRecorder()

			tt.mockBehaviour()

			addURLBatchService := addurlbatch.NewAddURLBatchService(
				urlRepositoryMock,
				userServiceMock,
				keyGenMock,
				baseURL,
				logger,
			)
			handler := &AddURLBatchHandler{
				logger:             logger,
				addURLBatchService: addURLBatchService,
			}
			handler.AddURLBatch(w, r)
			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
		})
	}
}
