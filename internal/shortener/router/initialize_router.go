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
	"github.com/anoriar/shortener/internal/shortener/middleware/auth"
	"github.com/anoriar/shortener/internal/shortener/middleware/compress"
	loggerMiddlewarePkg "github.com/anoriar/shortener/internal/shortener/middleware/logger"
	v1 "github.com/anoriar/shortener/internal/shortener/services/auth"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/util"
)

// InitializeRouter missing godoc.
func InitializeRouter(app *app.App) *Router {
	userRepository := app.UserRepository
	urlRepository := app.URLRepository
	userService := app.UserService
	return NewRouter(
		addurlhandler.NewAddHandler(urlRepository, userService, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), app.Logger, app.Config.BaseURL),
		geturlhandler.NewGetHandler(urlRepository, app.Logger),
		addURLHandlerV2.NewAddHandler(urlRepository, urlgen.NewShortURLGenerator(urlRepository, util.NewKeyGen()), userService, app.Logger, app.Config.BaseURL),
		addurlbatchhander.InitializeAddURLBatchHandler(urlRepository, userService, util.NewKeyGen(), app.Logger, app.Config.BaseURL),
		getuserurlshandler.InitializeGetUserURLsHandler(urlRepository, userService, app.Logger, app.Config.BaseURL),
		ping.NewPingHandler(urlRepository, app.Logger),
		deleteurlbatchhandler.NewDeleteURLBatchHandler(urlRepository, app.Logger),
		deleteuserurlshandler.NewDeleteUserURLsHandler(app.DeleteUserURLsProcessor, app.Logger),
		loggerMiddlewarePkg.NewLoggerMiddleware(app.Logger),
		compress.NewCompressMiddleware(),
		auth.NewAuthMiddleware(v1.NewSignService(app.Config.AuthSecretKey), userRepository),
	)
}
