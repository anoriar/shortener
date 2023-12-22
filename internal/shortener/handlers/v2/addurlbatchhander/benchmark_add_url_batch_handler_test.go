package addurlbatchhander

import (
	"bytes"
	"context"
	"encoding/json"
	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/validator"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/google/uuid"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const testURL = "https://github.com/"
const urlCnt = 1000000
const batchSize = 1000

func Benchmark_AddURLBatch(b *testing.B) {
	batchRequests := make([][]byte, 0, int(math.Ceil(float64(urlCnt)/float64(batchSize))))
	generatedURLs := make([]request.AddURLBatchRequestDTO, 0, batchSize)
	for i := 0; i < urlCnt; i++ {

		if i%batchSize == 0 {
			batchRequest, err := json.Marshal(generatedURLs)
			if err != nil {
				b.Fatalf("%s", err)
			}
			batchRequests = append(batchRequests, batchRequest)
			generatedURLs = make([]request.AddURLBatchRequestDTO, 0, batchSize)
		}
		url := request.AddURLBatchRequestDTO{
			CorrelationID: uuid.NewString(),
			OriginalURL:   testURL + strconv.Itoa(i),
		}
		generatedURLs = append(generatedURLs, url)
	}

	logger, err := logger.Initialize("info")
	if err != nil {
		b.Fatalf("%s", err)
	}
	keyGen := util.NewKeyGen()
	userRepository := inmemoryuser.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	urlRepository := inmemoryurl.NewInMemoryURLRepository()
	addURLBatchHandler := NewAddURLBatchHandler(urlRepository, userService, factory.NewAddURLBatchFactory(keyGen), factory.NewAddURLBatchResponseFactory("http://localhost:8080"), logger, validator.NewAddURLBatchValidator())

	b.ResetTimer()
	b.Run("add url batch", func(b *testing.B) {
		for _, batchRequest := range batchRequests {

			b.StopTimer()
			ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, "1")

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(batchRequest))
			req.Header.Set("Content-Type", "application/json")
			req.WithContext(ctxWithUser)
			addURLWriter := httptest.NewRecorder()

			b.StartTimer()

			addURLBatchHandler.AddURLBatch(addURLWriter, req)
		}
	})

}
