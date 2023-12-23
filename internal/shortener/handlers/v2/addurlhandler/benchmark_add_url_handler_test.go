package addurlhandler

import (
	"bytes"
	"context"
	"encoding/json"
	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
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
const urlCnt = 10000

func Benchmark_AddOneURLV2(b *testing.B) {

	urlAddRequests := make([][]byte, 0, urlCnt)
	for i := 0; i < urlCnt; i++ {
		requestDto := request.AddURLRequestDto{URL: testURL + strconv.Itoa(i)}
		successRequestBody, err := json.Marshal(requestDto)
		if err != nil {
			b.Fatalf("%s", err)
		}

		urlAddRequests = append(urlAddRequests, successRequestBody)
	}

	logger, err := logger.Initialize("info")
	if err != nil {
		b.Fatalf("%s", err)
	}
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addHandler := NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), userService, logger, "http://localhost:8080")

	b.ResetTimer()
	b.Run("add url v2", func(b *testing.B) {
		for _, addURLRequest := range urlAddRequests {

			b.StopTimer()
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(addURLRequest))

			ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, "1")
			req = req.WithContext(ctxWithUser)
			addURLWriter := httptest.NewRecorder()

			b.StartTimer()

			addHandler.AddURL(addURLWriter, req)
		}
	})

}
