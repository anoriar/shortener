package geturlhandler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	"github.com/anoriar/shortener/internal/shortener/util"
)

func Example() {
	const testURL = "https://github.com/"

	keyGen := util.NewKeyGen()
	urlRepository := inmemory.NewInMemoryURLRepository()
	logger, err := logger.Initialize("info")

	shortKey := keyGen.Generate()
	err = urlRepository.AddURL(&entity.URL{
		UUID:        uuid.NewString(),
		ShortURL:    shortKey,
		OriginalURL: testURL,
		IsDeleted:   false,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/"+shortKey, nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	NewGetHandler(urlRepository, logger).GetURL(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Location"))

	// Output:
	// 307
	// https://github.com/
}
