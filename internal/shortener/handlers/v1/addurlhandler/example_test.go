package addurlhandler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemory2 "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

// This example demonstrates how to use the AddURL.
func Example() {

	urlRepository := inmemory.NewInMemoryURLRepository()
	userRepository := inmemory2.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlGenerator := urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen())
	logger, err := logger.Initialize("info")

	ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, userID)

	req, err := http.NewRequest("POST", "/", strings.NewReader("https://github.com"))
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctxWithUser)

	w := httptest.NewRecorder()

	NewAddHandler(urlRepository, userService, urlGenerator, logger, baseURL).AddURL(w, req)

	fmt.Println(w.Code)

	// Output:
	// 201
}
