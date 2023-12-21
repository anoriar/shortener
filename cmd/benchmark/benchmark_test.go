package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/anoriar/shortener/internal/e2e/config"
	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	addHandlerV1 "github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	addHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const testURL = "https://github.com/"

func Benchmark_AddGetOneURLV1(b *testing.B) {
	const urlCnt = 10000
	urlAddRequests := make([]*http.Request, 0, urlCnt)
	for i := 0; i < urlCnt; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(testURL+strconv.Itoa(i))))
		urlAddRequests = append(urlAddRequests, req)
	}

	cnf := config.NewTestConfig()
	cnf.BaseURL = "http://localhost:8080"

	logger, err := logger.Initialize("info")
	if err != nil {
		b.Fatalf("%s", err)
	}
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addHandler := addHandlerV1.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL)

	getHandler := geturlhandler.NewGetHandler(urlRepository, logger)

	b.ResetTimer()
	b.Run("add url v1", func(b *testing.B) {
		for _, addURLRequest := range urlAddRequests {

			b.StopTimer()
			addURLWriter := httptest.NewRecorder()
			b.StartTimer()

			addHandler.AddURL(addURLWriter, addURLRequest)

			location := addURLWriter.Header().Get("Location")

			getURLRequest := httptest.NewRequest(http.MethodGet, "/"+location, nil)
			getURLWriter := httptest.NewRecorder()
			getHandler.GetURL(getURLWriter, getURLRequest)
		}
	})

}

func Benchmark_AddGetOneURLV2(b *testing.B) {
	const urlCnt = 10000
	urlAddRequests := make([][]byte, 0, urlCnt)
	for i := 0; i < urlCnt; i++ {
		requestDto := request.AddURLRequestDto{URL: testURL + strconv.Itoa(i)}
		successRequestBody, err := json.Marshal(requestDto)
		if err != nil {
			b.Fatalf("%s", err)
		}

		urlAddRequests = append(urlAddRequests, successRequestBody)
	}

	cnf := config.NewTestConfig()
	cnf.BaseURL = "http://localhost:8080"

	logger, err := logger.Initialize("info")
	if err != nil {
		b.Fatalf("%s", err)
	}
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addHandler := addHandlerV2.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), userService, logger, cnf.BaseURL)

	getHandler := geturlhandler.NewGetHandler(urlRepository, logger)

	b.ResetTimer()
	b.Run("add url v2", func(b *testing.B) {
		for _, addURLRequest := range urlAddRequests {

			b.StopTimer()
			ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, "1")

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(addURLRequest))
			req.Header.Set("Content-Type", "application/json")
			req.WithContext(ctxWithUser)
			addURLWriter := httptest.NewRecorder()

			b.StartTimer()

			addHandler.AddURL(addURLWriter, req)

			if addURLWriter.Code == http.StatusCreated {
				addURLResponseDto := &response.AddURLResponseDto{}
				err := json.Unmarshal(addURLWriter.Body.Bytes(), addURLResponseDto)
				if err != nil {
					b.Fatalf("%s", err)
				}
				location := addURLResponseDto.Result

				getURLRequest := httptest.NewRequest(http.MethodGet, location, nil)
				getURLWriter := httptest.NewRecorder()
				getHandler.GetURL(getURLWriter, getURLRequest)
			}

		}
	})

}
