package router

import (
	"github.com/anoriar/shortener/internal/app"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/geturlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v1/ping"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlbatchhander"
	addURLHandlerV2 "github.com/anoriar/shortener/internal/shortener/handlers/v2/addurlhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/deleteurlbatchhandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/deleteuserurlshandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/getuserurlshandler"
	"github.com/anoriar/shortener/internal/shortener/handlers/v2/statshandler"
	"github.com/anoriar/shortener/internal/shortener/middleware/auth"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	v1 "github.com/anoriar/shortener/internal/shortener/services/auth"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/usecases"
	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch"
	"github.com/anoriar/shortener/internal/shortener/usecases/getuserurlbatch"
	"github.com/anoriar/shortener/internal/shortener/util"
)

// InitializeRouter missing godoc.
func InitializeRouter(app *app.App) *Router {
	userRepository := app.UserRepository
	urlRepository := app.URLRepository
	userService := app.UserService
	keyGen := util.NewKeyGen()
	return NewRouter(
		addurlhandler.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, keyGen), app.Logger, app.Config.BaseURL),
		geturlhandler.NewGetHandler(app.Logger, usecases.NewGetURLService(urlRepository, app.Logger)),
		addURLHandlerV2.NewAddHandler(app.Logger, usecases.NewAddURLService(urlRepository, urlgen.NewShortURLGenerator(urlRepository, keyGen), userService, app.Logger, app.Config.BaseURL)),
		addurlbatchhander.NewAddURLBatchHandler(app.Logger, addurlbatch.NewAddURLBatchService(urlRepository, userService, keyGen, app.Config.BaseURL, app.Logger)),
		getuserurlshandler.NewGetUserURLsHandler(app.Logger, getuserurlbatch.NewGetUserURLsService(urlRepository, userService, app.Logger, app.Config.BaseURL)),
		ping.NewPingHandler(urlRepository, app.Logger),
		deleteurlbatchhandler.NewDeleteURLBatchHandler(urlRepository, app.Logger),
		deleteuserurlshandler.NewDeleteUserURLsHandler(usecases.NewDeleteUserURLsService(app.DeleteUserURLsProcessor, app.Logger), app.Logger),
		statshandler.NewStatsHandler(app.StatsService, app.Logger),
		loggerMiddlewarePkg.NewLoggerMiddleware(app.Logger),
		compress.NewCompressMiddleware(),
		auth.NewAuthMiddleware(v1.NewSignService(app.Config.AuthSecretKey), userRepository),
		auth.NewInternalAuthMiddleware(app.Config, app.Logger),
	)
}
