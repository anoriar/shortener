package db

import (
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/repository/url/db/internal"
)

// InitializeDBURLRepository missing godoc.
func InitializeDBURLRepository(dsn string, logger *zap.Logger) (*DatabaseURLRepository, error) {
	db, err := internal.InitializeDatabase(dsn)
	if err != nil {
		return nil, err
	}

	return NewDBURLRepository(db, logger), nil
}
