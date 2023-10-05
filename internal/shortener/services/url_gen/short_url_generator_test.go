package urlgen

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectedShortKey = "s7Fh4G"

var errRuntime = errors.New("exception")

type mockKeyGen struct{}

func (mock *mockKeyGen) Generate() string {
	return expectedShortKey
}

type mockURLRepositorySuccess struct{}

func (mcr *mockURLRepositorySuccess) AddURL(url *entity.URL) (*entity.URL, error) {
	return nil, errRuntime
}

func (mcr *mockURLRepositorySuccess) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	return nil, nil
}

type mockURLRepositoryError struct{}

func (mcr *mockURLRepositoryError) AddURL(url *entity.URL) (*entity.URL, error) {
	return nil, errors.New("exception")
}

func (mcr *mockURLRepositoryError) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	return nil, errors.New("exception")
}

type mockURLRepositoryEveryoneExisted struct{}

func (mcr *mockURLRepositoryEveryoneExisted) AddURL(url *entity.URL) (*entity.URL, error) {
	return &entity.URL{
		UUID:        "111",
		ShortURL:    "222",
		OriginalURL: "333",
	}, nil
}

func (mcr *mockURLRepositoryEveryoneExisted) FindURLByShortURL(shortURL string) (*entity.URL, error) {
	return &entity.URL{
		UUID:        "111",
		ShortURL:    "222",
		OriginalURL: "333",
	}, nil
}

func TestShortURLGenerator_generateShortURL(t *testing.T) {
	type fields struct {
		urlRepository repository.URLRepositoryInterface
		keyGen        util.KeyGenInterface
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				urlRepository: new(mockURLRepositorySuccess),
				keyGen:        new(mockKeyGen),
			},
			want:    expectedShortKey,
			wantErr: nil,
		},
		{
			name: "exception repository",
			fields: fields{
				urlRepository: new(mockURLRepositoryError),
				keyGen:        new(mockKeyGen),
			},
			want:    "",
			wantErr: errRuntime,
		},
		{
			name: "attempts exceeded",
			fields: fields{
				urlRepository: new(mockURLRepositoryEveryoneExisted),
				keyGen:        new(mockKeyGen),
			},
			want:    "",
			wantErr: ErrShortKeyGenerationAttemptsExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sug := &ShortURLGenerator{
				urlRepository: tt.fields.urlRepository,
				keyGen:        tt.fields.keyGen,
			}
			got, err := sug.GenerateShortURL()

			if tt.wantErr == nil && err != nil {
				t.Errorf("generateShortURL() exception:  %v", err)
			}

			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("generateShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
