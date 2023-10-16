package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/ping"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	"github.com/anoriar/shortener/internal/shortener/repository"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/util"
	"go.uber.org/zap"
)

func InitializeRouter(cnf *config.Config, urlRepository repository.URLRepositoryInterface, logger *zap.Logger) (*Router, error) {

	return NewRouter(
		addurlhandler.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL),
		geturlhandler.NewGetHandler(urlRepository, logger),
		addURLHandlerV2.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL),
		addurlbatchhander.InitializeAddURLBatchHandler(urlRepository, util.NewKeyGen(), logger, cnf.BaseURL),
		ping.NewPingHandler(urlRepository, logger),
		loggerMiddlewarePkg.NewLoggerMiddleware(logger),
		compress.NewCompressMiddleware(),
	), nil
}
