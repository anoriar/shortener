package internal

import (
	"database/sql"
	migrations2 "github.com/anoriar/shortener/internal/shortener/repository/db/internal/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitializeDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func PrepareDatabase(db *sql.DB) error {
	err := migrations2.Version231009Up(db)
	if err != nil {
		return err
	}

	err = migrations2.Version231015Up(db)
	if err != nil {
		return err
	}

	return nil
}
