package urlgen

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/mock"
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

	urlRepositoryMock := mock.NewMockURLRepositoryInterface(ctrl)

	keyGenMock := utilMock.NewMockKeyGenInterface(ctrl)

	tests := []struct {
		name          string
		mockBehaviour func()
		want          string
		wantErr       error
	}{
		{
			name: "success",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(nil, nil).Times(1)
				keyGenMock.EXPECT().Generate().Return(expectedShortKey).Times(1)
			},
			want:    expectedShortKey,
			wantErr: nil,
		},
		{
			name: "exception repository",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(nil, errors.New("exception")).Times(1)
				keyGenMock.EXPECT().Generate().Return(expectedShortKey).Times(1)
			},
			want:    "",
			wantErr: errRuntime,
		},
		{
			name: "attempts exceeded",
			mockBehaviour: func() {
				urlRepositoryMock.EXPECT().FindURLByShortURL(gomock.Any()).Return(&entity.URL{
					UUID:        "111",
					ShortURL:    "222",
					OriginalURL: "333",
				}, nil).Times(maxAttempts)
				keyGenMock.EXPECT().Generate().Return(expectedShortKey).Times(5)
			},
			want:    "",
			wantErr: ErrShortKeyGenerationAttemptsExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			sug := &ShortURLGenerator{
				urlRepository: urlRepositoryMock,
				keyGen:        keyGenMock,
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
