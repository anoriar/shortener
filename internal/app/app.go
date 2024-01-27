package app

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/config"
	deleteurlsprocessor "github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/repository/user"
	"github.com/anoriar/shortener/internal/shortener/repository/user/inmemory"
	"github.com/anoriar/shortener/internal/shortener/services/deleteuserurls"
	userServicePkg "github.com/anoriar/shortener/internal/shortener/services/user"
)

// App missing godoc.
type App struct {
	Config                  *config.Config
	Logger                  *zap.Logger
	UserRepository          user.UserRepositoryInterface
	UserService             userServicePkg.UserServiceInterface
	URLRepository           url.URLRepositoryInterface
	DeleteUserURLsService   *deleteuserurls.DeleteUserURLsService
	DeleteUserURLsProcessor *deleteurlsprocessor.DeleteUserURLsProcessor
}

// NewApp missing godoc.
func NewApp(cnf *config.Config, logger *zap.Logger) (*App, error) {
	userRepository := inmemory.NewInMemoryUserRepository()
	urlRepository, err := url.InitializeURLRepository(cnf, logger)
	if err != nil {
		return nil, fmt.Errorf("init repository error %v", err.Error())
	}
	userService := userServicePkg.NewUserService(userRepository)
	deleteUserURLsService := deleteuserurls.NewDeleteUserURLsService(urlRepository, userRepository)
	deleteUserURLsProcessor := deleteurlsprocessor.NewDeleteUserURLsProcessor(deleteUserURLsService, logger)
	return &App{
		Config:                  cnf,
		Logger:                  logger,
		UserRepository:          userRepository,
		UserService:             userService,
		URLRepository:           urlRepository,
		DeleteUserURLsService:   deleteUserURLsService,
		DeleteUserURLsProcessor: deleteUserURLsProcessor,
	}, nil
}
