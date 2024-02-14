package addurlbatchhander

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch"

	"github.com/google/uuid"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/logger"
	inmemoryurl "github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
	inmemoryuser "github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
)

const testURL = "https://github.com/"
const urlCnt = 1000000
const batchSize = 1000

func Benchmark_AddURLBatch(b *testing.B) {
	batchRequests := make([][]byte, 0, int(math.Ceil(float64(urlCnt)/float64(batchSize))))
	generatedURLs := make([]request.AddURLBatchRequestDTO, 0, batchSize)
	for i := 0; i < urlCnt; i++ {

		if i%batchSize == 0 && i != 0 {
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

	addURLBatchService := addurlbatch.NewAddURLBatchService(
		urlRepository,
		userService,
		keyGen,
		baseURL,
		logger,
	)

	addURLBatchHandler := NewAddURLBatchHandler(logger, addURLBatchService)

	b.ResetTimer()
	b.Run("add url batch", func(b *testing.B) {
		for _, batchRequest := range batchRequests {

			b.StopTimer()
			ctxWithUser := context.WithValue(context.Background(), context2.UserIDContextKey, "1")

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(batchRequest))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctxWithUser)
			addURLWriter := httptest.NewRecorder()

			b.StartTimer()

			addURLBatchHandler.AddURLBatch(addURLWriter, req)
		}
	})

}
