package db

import (
	"context"
	"database/sql"
	"github.com/anoriar/shortener/internal/shortener/db/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

// #MENTOR: Стоит ли делать здесь DatabaseInterface?
// А это как конкретная реализация - PGDatabase
// В этом случае придется переопределять все методы sql.DB
// Если здесь сделать интерфейс - можно будет протестить pingHandler через мок

// #MENTOR: sql.DB позволяет писать любые запросы. Например, CREATE DATABASE IF NOT EXISTS.
// Но для postgres не работает IF NOT EXISTS и выведет ошибку
func InitializeDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = PingDatabase(db)
	if err != nil {
		return nil, err
	}

	migrations.Up(db)

	return db, nil
}

func PingDatabase(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	return err
}
