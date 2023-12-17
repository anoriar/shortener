package benchmark

import (
	"github.com/anoriar/shortener/internal/e2e/config"
	addurlhandler2 "github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const successRequestBody = "https://github.com"
const testURL = "https://github.com/"

func Benchmark_Shortener(b *testing.B) {
	cnf := config.NewTestConfig()
	cnf.BaseURL = "http://localhost:8080"

	logger, err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addHandler := addurlhandler2.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL)

	b.ResetTimer()
	b.Run("1", func(b *testing.B) {
		for i := 0; i < 1000; i++ {
			b.StopTimer()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(successRequestBody+strconv.Itoa(i)))
			w := httptest.NewRecorder()
			b.StartTimer()

			addHandler.AddURL(w, r)
		}
	})

}
