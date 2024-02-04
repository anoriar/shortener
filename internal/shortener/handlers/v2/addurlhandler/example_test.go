package addurlhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/anoriar/shortener/internal/shortener/dto/request"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/logger"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

// This example demonstrates how to use the GetStats.
func Example() {
	const testURL = "https://github.com/"

	logger, err := logger.Initialize("info")
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addHandler := NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), userService, logger, "http://localhost:8080")

	requestDto := request.AddURLRequestDto{URL: testURL}
	successRequestBody, err := json.Marshal(requestDto)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userID)

	req, err := http.NewRequest("POST", "/api/shorten", strings.NewReader(string(successRequestBody)))
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctxWithUser)

	w := httptest.NewRecorder()

	addHandler.AddURL(w, req)

	fmt.Println(w.Code)

	// Output:
	// 201
}
