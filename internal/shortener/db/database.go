package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func InitializeDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func PingDatabase(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	return err
}
