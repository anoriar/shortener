package benchmark

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"

	"github.com/anoriar/shortener/internal/e2e/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

const testURL = "https://github.com/"

func Benchmark_Shortener(b *testing.B) {
	memProfileFile, err := os.Create("memprofile.out")
	if err != nil {
		b.Fatalf("Error creating memory profile: %v", err)
	}
	defer memProfileFile.Close()

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
		panic(err)
	}
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addHandler := addurlhandler.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL)

	getHandler := geturlhandler.NewGetHandler(urlRepository, logger)

	b.ResetTimer()
	pprof.WriteHeapProfile(memProfileFile)
	b.Run("1", func(b *testing.B) {
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
