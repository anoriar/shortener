package addurlbatchhander

import (
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander/internal/validator"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
	"go.uber.org/zap"
)

func InitializeAddURLBatchHandler(
	urlRepository url.URLRepositoryInterface,
	userService *user.UserService,
	keyGen util.KeyGenInterface,
	logger *zap.Logger,
	baseURL string,
) *AddURLBatchHandler {
	return NewAddURLBatchHandler(urlRepository, userService, factory.NewAddURLBatchFactory(keyGen), factory.NewAddURLBatchResponseFactory(baseURL), logger, validator.NewAddURLBatchValidator())
}
