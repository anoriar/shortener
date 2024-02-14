package getuserurlshandler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/anoriar/shortener/internal/shortener/usecases/getuserurlbatch"

	"github.com/google/uuid"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	"github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

func Example() {
	const testURL = "https://github.com/"

	logger, err := logger.Initialize("info")
	if err != nil {
		log.Fatalf("%s", err)
	}
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	userRepository := inmemory.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)

	getUserURLsHandler := NewGetUserURLsHandler(logger, getuserurlbatch.NewGetUserURLsService(urlRepository, userService, logger, "http://localhost:8080"))

	userURLShortKeys := make([]string, 0, 1)
	shortKey := "fb3fwi"
	urlEntity := entity.URL{
		UUID:        uuid.NewString(),
		ShortURL:    shortKey,
		OriginalURL: testURL,
		IsDeleted:   false,
	}
	err = urlRepository.AddURL(&urlEntity)
	if err != nil {
		log.Fatalf("%s", err)
	}
	userURLShortKeys = append(userURLShortKeys, urlEntity.ShortURL)

	userEntity := entity.User{
		UUID:        uuid.NewString(),
		SavedURLIDs: userURLShortKeys,
	}

	err = userRepository.AddUser(userEntity)
	if err != nil {
		log.Fatalf("%s", err)
	}
	ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userEntity.UUID)
	req, err := http.NewRequest(http.MethodGet, "/api/user/urls", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctxWithUser)

	w := httptest.NewRecorder()

	getUserURLsHandler.GetUserURLs(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// [{"short_url":"http://localhost:8080/fb3fwi","original_url":"https://github.com/"}]
}
