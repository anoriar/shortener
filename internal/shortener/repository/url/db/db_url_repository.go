package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/dto/repository"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"strings"
)

type DatabaseURLRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewDBURLRepository(db *sql.DB, logger *zap.Logger) *DatabaseURLRepository {
	return &DatabaseURLRepository{db: db, logger: logger}
}

func (repository *DatabaseURLRepository) Ping(ctx context.Context) error {
	err := repository.db.PingContext(ctx)
	return err
}

func (repository *DatabaseURLRepository) AddURL(url *entity.URL) error {
	_, err := repository.db.Exec("INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3);", url.UUID, url.ShortURL, url.OriginalURL)
	if err != nil {
		repository.logger.Error(err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return repositoryerror.ErrConflict
		}
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

func (repository *DatabaseURLRepository) FindURLByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	row := repository.db.QueryRowContext(ctx, "SELECT uuid, short_url, original_url FROM urls WHERE original_url = $1 LIMIT 1", originalURL)
	var url entity.URL
	err := row.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &url, nil
}

func (repository *DatabaseURLRepository) INQuery(ctx context.Context, urls []string) ([]entity.URL, error) {
	var resultUrls []entity.URL
	values := make([]interface{}, len(urls))
	for i, id := range urls {
		values[i] = id
	}
	query := "SELECT * FROM urls WHERE short_url IN ($1, $2, $3)"

	rows, err := repository.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url entity.URL
		err := rows.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL)
		if err != nil {
			repository.logger.Error(err.Error())
			return nil, err
		}
		resultUrls = append(resultUrls, url)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resultUrls, nil
}

func (repository *DatabaseURLRepository) GetURLsByQuery(ctx context.Context, urlQuery repository.Query) ([]entity.URL, error) {
	var resultUrls []entity.URL

	paramCounter := 1
	var filters []string
	var filterParams []string
	if len(urlQuery.OriginalURLs) > 0 {
		placeholders := make([]string, len(urlQuery.OriginalURLs))
		for i := range urlQuery.OriginalURLs {
			placeholders[i] = fmt.Sprintf("$%d", paramCounter)
			paramCounter++
		}

		filters = append(filters, fmt.Sprintf("original_url IN (%s)", strings.Join(placeholders, ", ")))
		filterParams = append(filterParams, urlQuery.OriginalURLs...)
	}
	if len(urlQuery.ShortURLs) > 0 {
		placeholders := make([]string, len(urlQuery.ShortURLs))
		for i := range urlQuery.ShortURLs {
			placeholders[i] = fmt.Sprintf("$%d", paramCounter)
			paramCounter++
		}
		filters = append(filters, fmt.Sprintf("short_url IN (%s)", strings.Join(placeholders, ", ")))
		filterParams = append(filterParams, urlQuery.ShortURLs...)
	}

	filterString := strings.Join(filters, " AND ")
	queryString := "SELECT uuid, short_url, original_url FROM urls"
	if filterString != "" && len(filterParams) != 0 {
		queryString += " WHERE " + filterString
	}

	//#MENTOR: вот эта штука супернеочевидная. нужно именно тип interface, string не подходит
	params := make([]interface{}, len(filterParams))
	for i, filterParam := range filterParams {
		params[i] = filterParam
	}

	rows, err := repository.db.QueryContext(ctx, queryString, params...)
	if err != nil {
		repository.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url entity.URL
		err := rows.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL)
		if err != nil {
			repository.logger.Error(err.Error())
			return nil, err
		}
		resultUrls = append(resultUrls, url)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resultUrls, nil
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

func (repository *DatabaseURLRepository) DeleteURLBatch(ctx context.Context, shortURLs []string) error {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		repository.logger.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("DELETE FROM urls WHERE short_url=$1")
	if err != nil {
		repository.logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	for _, shortURL := range shortURLs {
		_, err := stmt.ExecContext(ctx, shortURL)
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

func (repository *DatabaseURLRepository) Close() error {
	err := repository.db.Close()
	if err != nil {
		return err
	}
	return nil
}
