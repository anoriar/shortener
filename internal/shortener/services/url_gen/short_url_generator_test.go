package urlgen

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/repository/mock"
	"github.com/anoriar/shortener/internal/shortener/util"
	utilMock "github.com/anoriar/shortener/internal/shortener/util/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectedShortKey = "s7Fh4G"

var errRuntime = errors.New("exception")

func TestShortURLGenerator_generateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	urlRepoSuccessMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlRepoSuccessMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(nil, nil).MaxTimes(1).MinTimes(1)

	urlRepoErrorMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlRepoErrorMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(nil, errors.New("exception")).MaxTimes(1).MinTimes(1)

	urlRepoEveryoneExistedMock := mock.NewMockURLRepositoryInterface(ctrl)
	urlRepoEveryoneExistedMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(&entity.URL{
		UUID:        "111",
		ShortURL:    "222",
		OriginalURL: "333",
	}, nil).MaxTimes(maxAttempts).MinTimes(maxAttempts)

	keyGenMock := utilMock.NewMockKeyGenInterface(ctrl)
	keyGenMock.EXPECT().Generate().Return(expectedShortKey).AnyTimes()

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
				urlRepository: urlRepoSuccessMock,
				keyGen:        keyGenMock,
			},
			want:    expectedShortKey,
			wantErr: nil,
		},
		{
			name: "exception repository",
			fields: fields{
				urlRepository: urlRepoErrorMock,
				keyGen:        keyGenMock,
			},
			want:    "",
			wantErr: errRuntime,
		},
		{
			name: "attempts exceeded",
			fields: fields{
				urlRepository: urlRepoEveryoneExistedMock,
				keyGen:        keyGenMock,
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
