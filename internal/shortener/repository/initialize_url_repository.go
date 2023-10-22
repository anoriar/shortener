package repository

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	dbURLRepository "github.com/anoriar/shortener/internal/shortener/repository/db"
	"github.com/anoriar/shortener/internal/shortener/repository/file"
	"go.uber.org/zap"
)

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
		return NewInMemoryURLRepository(), nil
	}
}
