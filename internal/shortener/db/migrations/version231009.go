package migrations

import "database/sql"

func Version231009Up(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS urls (uuid VARCHAR(36) NOT NULL, short_url VARCHAR(255) NOT NULL, original_url VARCHAR(255) NOT NULL, PRIMARY KEY (uuid));")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS short_url_idx ON urls (short_url);")
	if err != nil {
		return err
	}

	return nil
}

func Version231009Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS urls")
	if err != nil {
		return err
	}

	return nil
}
