package db

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	dbName      = "shortener_test"
	validDSN    = "postgres://postgres:postgres@localhost:5434"
	notValidDSN = "postgres://postgres:postgres@193.22.11.33:5434"
)

func TestPingDatabase(t *testing.T) {

	validDB, err := sql.Open("pgx", validDSN)
	require.NoError(t, err)
	defer func() {
		_, err := validDB.Exec(fmt.Sprintf("drop database if exists %s", dbName))
		require.NoError(t, err)
		validDB.Close()
	}()
	_, err = validDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err)

	notValidDB, err := sql.Open("pgx", notValidDSN)
	require.NoError(t, err)

	defer notValidDB.Close()

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid DSN",
			args: args{
				db: validDB,
			},
			wantErr: false,
		},
		{
			name: "not valid DSN",
			args: args{
				db: notValidDB,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PingDatabase(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("PingDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
