package db

import (
	"context"
	"database/sql"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"go.uber.org/zap"
)

type DatabaseURLRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewDBURLRepository(db *sql.DB, logger *zap.Logger) *DatabaseURLRepository {
	return &DatabaseURLRepository{db: db, logger: logger}
}

// #MENTOR Вопрос 1: Целесообразно ли передавать в функции db репозитория контекст запроса? В курсе видел, что лучше передавать
// С точки зрения DDD - плохая практика, а также замусоривает  параметры функции
// #MENTOR: Вопрос 2: Может ли сервер потерять связь с базой данных? Например, стоит консьюмер, который работает 20 дней.
// Нужно ли проверять коннект к бд и реконнектить в случае неуспешного пинга? Как отличить: потерялся коннект к бд или ошибка запроса?
func (repository *DatabaseURLRepository) AddURL(url *entity.URL) error {
	_, err := repository.db.Exec("INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3);", url.UUID, url.ShortURL, url.OriginalURL)
	if err != nil {
		repository.logger.Error(err.Error())
		return err
	}
	return nil
}

func (repository *DatabaseURLRepository) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	rows, err := repository.db.Query("SELECT uuid, short_url, original_url FROM urls WHERE short_url = $1 LIMIT 1", shortURL)
	if err != nil {
		repository.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var url entity.URL

	for rows.Next() {
		err := rows.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL)
		if err != nil {
			repository.logger.Error(err.Error())
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if url.UUID == "" {
		return nil, nil
	}

	return &url, err
}

func (repository *DatabaseURLRepository) AddURLBatch(ctx context.Context, urls []entity.URL) error {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		repository.logger.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO urls (uuid, short_url, original_url) VALUES ($1,$2,$3)")
	if err != nil {
		repository.logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	for _, url := range urls {
		_, err := stmt.ExecContext(ctx, url.UUID, url.ShortURL, url.OriginalURL)
		if err != nil {
			repository.logger.Error(err.Error())
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		repository.logger.Error(err.Error())
		return err
	}

	return nil
}
