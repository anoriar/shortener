package db

import (
	"github.com/anoriar/shortener/internal/shortener/repository/db/internal"
	"go.uber.org/zap"
)

func InitializeDBURLRepository(dsn string, logger *zap.Logger) (*DatabaseURLRepository, error) {
	db, err := internal.InitializeDatabase(dsn)
	if err != nil {
		return nil, err
	}

	err = internal.PrepareDatabase(db)
	if err != nil {
		return nil, err
	}

	return NewDBURLRepository(db, logger), nil
}
