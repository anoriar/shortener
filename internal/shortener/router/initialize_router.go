package router

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/ping"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/deleteurlbatchhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/getuserurlshandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/auth"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	v1 "github.com/anoriar/shortener/internal/shortener/services/auth"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/services/user"
	"github.com/anoriar/shortener/internal/shortener/util"
	"go.uber.org/zap"
)

func InitializeRouter(cnf *config.Config, urlRepository url.URLRepositoryInterface, logger *zap.Logger) (*Router, error) {
	userRepository := inmemory.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	return NewRouter(
		addurlhandler.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), logger, cnf.BaseURL),
		geturlhandler.NewGetHandler(urlRepository, logger),
		addURLHandlerV2.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), userService, logger, cnf.BaseURL),
		addurlbatchhander.InitializeAddURLBatchHandler(urlRepository, userService, util.NewKeyGen(), logger, cnf.BaseURL),
		getuserurlshandler.InitializeGetUserURLsHandler(urlRepository, userService, logger, cnf.BaseURL),
		ping.NewPingHandler(urlRepository, logger),
		deleteurlbatchhandler.NewDeleteURLBatchHandler(urlRepository, logger),
		loggerMiddlewarePkg.NewLoggerMiddleware(logger),
		compress.NewCompressMiddleware(),
		auth.NewAuthMiddleware(v1.NewSignService(cnf.AuthSecretKey), userRepository),
	), nil
}
