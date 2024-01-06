package url

import (
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/config"
	dbURLRepository "github.com/anoriar/shortener/internal/shortener/repository/url/db"
	"github.com/anoriar/shortener/internal/shortener/repository/url/file"
	"github.com/anoriar/shortener/internal/shortener/repository/url/inmemory"
)

// InitializeURLRepository missing godoc.
func InitializeURLRepository(cnf *config.Config, logger *zap.Logger) (URLRepositoryInterface, error) {
	switch {
	case cnf.DatabaseDSN != "":
		urlRepository, err := dbURLRepository.InitializeDBURLRepository(cnf.DatabaseDSN, logger)
		if err != nil {
			return nil, err
		}
		return urlRepository, nil
	case cnf.FileStoragePath != "":
		return file.NewFileURLRepository(cnf.FileStoragePath), nil
	default:
		return inmemory.NewInMemoryURLRepository(), nil
	}
}
