package addurlhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/anoriar/shortener/internal/e2e/config"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

const testURL = "https://github.com/"
const urlCnt = 10000

func Benchmark_AddOneURLV1(b *testing.B) {
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
	addHandler := NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL)

	b.ResetTimer()
	b.Run("add url v1", func(b *testing.B) {
		for _, addURLRequest := range urlAddRequests {

			b.StopTimer()
			addURLWriter := httptest.NewRecorder()
			b.StartTimer()

			addHandler.AddURL(addURLWriter, addURLRequest)
		}
	})

}
