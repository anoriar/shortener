package geturlhandler

import (
	"github.com/anoriar/shortener/internal/e2e/config"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const testURL = "https://github.com/"
const urlCnt = 10000

func Benchmark_GetOneURLV1(b *testing.B) {
	cnf := config.NewTestConfig()
	cnf.BaseURL = "http://localhost:8080"

	logger, err := logger.Initialize("info")
	if err != nil {
		b.Fatalf("%s", err)
	}
	keyGen := util.NewKeyGen()
	urlRepository := inmemoryurl.NewInMemoryURLRepository()

	getHandler := NewGetHandler(urlRepository, logger)

	urlShortKeys := make([]string, 0, urlCnt)
	for i := 0; i < urlCnt; i++ {
		shortKey := keyGen.Generate()
		urlEntity := entity.URL{
			UUID:        uuid.NewString(),
			ShortURL:    shortKey,
			OriginalURL: testURL + strconv.Itoa(i),
			IsDeleted:   false,
		}
		err := urlRepository.AddURL(&urlEntity)
		if err != nil {
			b.Fatalf("%s", err)
		}
		urlShortKeys = append(urlShortKeys, shortKey)
	}

	b.ResetTimer()

	b.Run("get url v1", func(b *testing.B) {
		for _, key := range urlShortKeys {
			b.StopTimer()
			req := httptest.NewRequest(http.MethodGet, "/"+key, nil)
			w := httptest.NewRecorder()
			b.StartTimer()

			getHandler.GetURL(w, req)
		}
	})
}
