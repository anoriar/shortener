package app

import (
	"fmt"

	"github.com/anoriar/shortener/internal/shortener/services/auth"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/anoriar/shortener/internal/shortener/usecases"
	"github.com/anoriar/shortener/internal/shortener/usecases/addurlbatch"
	"github.com/anoriar/shortener/internal/shortener/usecases/getuserurlbatch"
	"github.com/anoriar/shortener/internal/shortener/util"

	"github.com/anoriar/shortener/internal/shortener/services/stats"

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
	StatsService            stats.StatsServiceInterface
	KeyGen                  util.KeyGenInterface
	Authenticator           *auth.Authenticator
	AddURLServiceUC         *usecases.AddURLService
	AddURLBatchServiceUC    *addurlbatch.AddURLBatchService
	GetURLServiceUC         *usecases.GetURLService
	GetUserURLsServiceUC    *getuserurlbatch.GetUserURLsService
	DeleteUserURLsServiceUC *usecases.DeleteUserURLsService
}

// NewApp missing godoc.
func NewApp(cnf *config.Config, logger *zap.Logger) (*App, error) {
	userRepository := inmemory.NewInMemoryUserRepository()
	urlRepository, err := url.InitializeURLRepository(cnf, logger)
	if err != nil {
		return nil, fmt.Errorf("init repository error %v", err.Error())
	}
	keyGen := util.NewKeyGen()
	userService := userServicePkg.NewUserService(userRepository)
	deleteUserURLsService := deleteuserurls.NewDeleteUserURLsService(urlRepository, userRepository)
	deleteUserURLsProcessor := deleteurlsprocessor.NewDeleteUserURLsProcessor(deleteUserURLsService, logger)

	urlGen := urlgen.NewShortURLGenerator(urlRepository, keyGen)

	statsService := stats.NewStatsService(urlRepository, userRepository)
	return &App{
		Config:                  cnf,
		Logger:                  logger,
		UserRepository:          userRepository,
		UserService:             userService,
		URLRepository:           urlRepository,
		DeleteUserURLsService:   deleteUserURLsService,
		DeleteUserURLsProcessor: deleteUserURLsProcessor,
		StatsService:            statsService,
		KeyGen:                  keyGen,
		Authenticator:           auth.NewAuthenticator(auth.NewSignService(cnf.AuthSecretKey), userRepository, logger),
		AddURLServiceUC:         usecases.NewAddURLService(urlRepository, urlGen, userService, logger, cnf.BaseURL),
		GetURLServiceUC:         usecases.NewGetURLService(urlRepository, logger),
		AddURLBatchServiceUC:    addurlbatch.NewAddURLBatchService(urlRepository, userService, keyGen, cnf.BaseURL, logger),
		GetUserURLsServiceUC:    getuserurlbatch.NewGetUserURLsService(urlRepository, userService, logger, cnf.BaseURL),
		DeleteUserURLsServiceUC: usecases.NewDeleteUserURLsService(deleteUserURLsProcessor, logger),
	}, nil
}
