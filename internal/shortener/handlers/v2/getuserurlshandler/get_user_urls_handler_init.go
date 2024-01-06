package getuserurlshandler

import (
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/handlers/v2/getuserurlshandler/internal/factory"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/services/user"
)

// InitializeGetUserURLsHandler missing godoc.
func InitializeGetUserURLsHandler(
	urlRepository url.URLRepositoryInterface,
	userService *user.UserService,
	logger *zap.Logger,
	baseURL string,
) *GetUserURLsHandler {
	return NewGetUserURLsHandler(urlRepository, userService, factory.NewGetUSerURLsResponseFactory(baseURL), logger)
}
