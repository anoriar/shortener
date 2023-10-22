package addurlbatchhander

import (
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/validator"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/util"
	"go.uber.org/zap"
)

func InitializeAddURLBatchHandler(urlRepository repository.URLRepositoryInterface, keyGen util.KeyGenInterface, logger *zap.Logger, baseURL string) *AddURLBatchHandler {
	return NewAddURLBatchHandler(urlRepository, factory.NewAddURLBatchFactory(keyGen), factory.NewAddURLBatchResponseFactory(baseURL), logger, validator.NewAddURLBatchValidator())
}
