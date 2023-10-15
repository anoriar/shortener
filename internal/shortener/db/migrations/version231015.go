package migrations

import "database/sql"

func Version231015Up(db *sql.DB) error {
	_, err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);")
	if err != nil {
		return err
	}

	return nil
}

func Version231015Down(db *sql.DB) error {
	_, err := db.Exec("DROP INDEX IF EXISTS original_url_idx;")
	if err != nil {
		return err
	}

	return nil
}
