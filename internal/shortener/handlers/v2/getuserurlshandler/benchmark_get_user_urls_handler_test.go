package getuserurlshandler

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/uuid"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/getuserurlshandler/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	"github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

const testURL = "https://github.com/"
const urlCnt = 10000
const userCnt = 1000

func Benchmark_GetUserURLs(b *testing.B) {
	logger, err := logger.Initialize("info")
	if err != nil {
		b.Fatalf("%s", err)
	}
	keyGen := util.NewKeyGen()
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	userRepository := inmemory.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	getUserURLsResponseFactory := factory.NewGetUSerURLsResponseFactory("http://localhost:8080")

	getUserURLsHandler := NewGetUserURLsHandler(urlRepository, userService, getUserURLsResponseFactory, logger)

	urlCntPerUser := int(math.Ceil(float64(urlCnt) / float64(userCnt)))
	urlShortKeys := make([]string, 0, urlCntPerUser)
	userIDs := make([]string, 0, userCnt)

	for i := 0; i < urlCnt; i++ {

		if i%urlCntPerUser == 0 {
			userEntity := entity.User{
				UUID:        uuid.NewString(),
				SavedURLIDs: urlShortKeys,
			}
			err := userRepository.AddUser(userEntity)
			if err != nil {
				b.Fatalf("%s", err)
			}
			userIDs = append(userIDs, userEntity.UUID)

			urlShortKeys = make([]string, 0, urlCntPerUser)
		}
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
		urlShortKeys = append(urlShortKeys, urlEntity.ShortURL)
	}

	b.ResetTimer()

	b.Run("get user urls", func(b *testing.B) {
		for _, userID := range userIDs {
			b.StopTimer()

			ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userID)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(ctxWithUser)
			w := httptest.NewRecorder()

			b.StartTimer()

			getUserURLsHandler.GetUserURLs(w, req)
		}
	})
}
