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
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/util"
)

// InitializeRouter missing godoc.
func InitializeRouter(app *app.App) *Router {
	urlRepository := app.URLRepository
	userService := app.UserService
	keyGen := util.NewKeyGen()
	return NewRouter(
		addurlhandler.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, keyGen), app.Logger, app.Config.BaseURL),
		geturlhandler.NewGetHandler(app.Logger, app.GetURLServiceUC),
		addURLHandlerV2.NewAddHandler(app.Logger, app.AddURLServiceUC),
		addurlbatchhander.NewAddURLBatchHandler(app.Logger, app.AddURLBatchServiceUC),
		getuserurlshandler.NewGetUserURLsHandler(app.Logger, app.GetUserURLsServiceUC),
		ping.NewPingHandler(urlRepository, app.Logger),
		deleteurlbatchhandler.NewDeleteURLBatchHandler(urlRepository, app.Logger),
		deleteuserurlshandler.NewDeleteUserURLsHandler(app.DeleteUserURLsServiceUC, app.Logger),
		statshandler.NewStatsHandler(app.StatsService, app.Logger),
		loggerMiddlewarePkg.NewLoggerMiddleware(app.Logger),
		compress.NewCompressMiddleware(),
		auth.NewAuthMiddleware(app.Authenticator),
		auth.NewInternalAuthMiddleware(app.Config, app.Logger),
	)
}
