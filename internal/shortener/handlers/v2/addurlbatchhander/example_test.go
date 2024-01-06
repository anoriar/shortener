package addurlbatchhander

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/validator"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

func Example() {

	logger, err := logger.Initialize("info")
	if err != nil {
		log.Fatalf("%s", err)
	}
	keyGen := util.NewKeyGen()
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addURLBatchHandler := NewAddURLBatchHandler(urlRepository, userService, factory.NewAddURLBatchFactory(keyGen), factory.NewAddURLBatchResponseFactory("http://localhost:8080"), logger, validator.NewAddURLBatchValidator())

	ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userID)

	successRequestBody, err := json.Marshal([]request.AddURLBatchRequestDTO{
		{
			CorrelationID: "g0fsdf9fj",
			OriginalURL:   "https://practicum2.yandex.ru",
		},
		{
			CorrelationID: "ngfdsf3",
			OriginalURL:   "https://practicum3.yandex.ru",
		},
	})

	if err != nil {
		log.Fatalf("%s", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(successRequestBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctxWithUser)
	w := httptest.NewRecorder()

	addURLBatchHandler.AddURLBatch(w, req)

	fmt.Println(w.Code)

	// Output:
	// 201
}
