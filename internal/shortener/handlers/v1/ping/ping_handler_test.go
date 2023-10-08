package ping

import (
	"database/sql"
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	dbName      = "shortener_test"
	validDSN    = "postgres://postgres:postgres@localhost:5434"
	notValidDSN = "postgres://postgres:postgres@193.22.11.33:5434"
)

func TestPingHandler_Ping(t *testing.T) {
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

	logger, err := logger.Initialize("info")
	require.NoError(t, err)

	type fields struct {
		db     *sql.DB
		logger *zap.Logger
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	type want struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "valid DSN",
			fields: fields{
				db:     validDB,
				logger: logger,
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "valid DSN",
			fields: fields{
				db:     notValidDB,
				logger: logger,
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			NewPingHandler(tt.fields.db, tt.fields.logger).Ping(w, r)

			assert.Equal(t, tt.want.status, w.Code)
		})
	}
}
