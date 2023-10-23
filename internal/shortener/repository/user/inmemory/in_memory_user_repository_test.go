package inmemory

import (
	"errors"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/repositoryerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryURLRepository_AddUser(t *testing.T) {
	const (
		userID1  = "46b8f9d2-b123-4f8e-aabb-f77dd764a00b"
		userID2  = "a077d055-6ec0-4cd5-adb2-6427f35045db"
		userURL1 = "76ca01d4-37db-4433-85d8-68cd59493430"
		userURL2 = "56dc7a32-7294-44d2-85bf-6cbbf0cbd563"
	)

	type args struct {
		user entity.User
	}
	tests := []struct {
		name         string
		existedUsers map[string]entity.User
		args         args
		wantErr      bool
		want         map[string]entity.User
		errorIs      error
	}{
		{
			name: "add item simple",
			existedUsers: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
			},
			args: args{
				entity.User{
					UUID: userID2,
					SavedURLIDs: []string{
						userURL1,
						userURL2,
					},
				},
			},
			want: map[string]entity.User{
				userID1: {
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
				userID2: {
					UUID: userID2,
					SavedURLIDs: []string{
						userURL1,
						userURL2,
					},
				},
			},
			wantErr: false,
			errorIs: nil,
		},
		{
			name: "conflict",
			existedUsers: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
			},
			args: args{
				entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
						userURL2,
					},
				},
			},
			want: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
			},
			wantErr: true,
			errorIs: repositoryerror.ErrConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &InMemoryUserRepository{
				tt.existedUsers,
			}

			err := repository.AddUser(tt.args.user)

			assert.Equal(t, tt.wantErr, err != nil)
			if tt.errorIs != nil {
				assert.True(t, errors.Is(err, tt.errorIs))
			}
			assert.Equal(t, tt.want, repository.users)
		})
	}
}

func TestInMemoryURLRepository_UpdateUser(t *testing.T) {
	const (
		userID1  = "46b8f9d2-b123-4f8e-aabb-f77dd764a00b"
		userURL1 = "76ca01d4-37db-4433-85d8-68cd59493430"
		userURL2 = "56dc7a32-7294-44d2-85bf-6cbbf0cbd563"
	)

	type args struct {
		user entity.User
	}
	tests := []struct {
		name         string
		existedUsers map[string]entity.User
		args         args
		wantErr      bool
		want         map[string]entity.User
		errorIs      error
	}{
		{
			name: "update item simple",
			existedUsers: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
			},
			args: args{
				entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
						userURL2,
					},
				},
			},
			want: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
						userURL2,
					},
				},
			},
			wantErr: false,
			errorIs: nil,
		},
		{
			name: "not found",
			existedUsers: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
			},
			args: args{
				entity.User{
					UUID: "31312",
					SavedURLIDs: []string{
						userURL1,
						userURL2,
					},
				},
			},
			want: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
					SavedURLIDs: []string{
						userURL1,
					},
				},
			},
			wantErr: true,
			errorIs: repositoryerror.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &InMemoryUserRepository{
				tt.existedUsers,
			}

			err := repository.UpdateUser(tt.args.user)

			assert.Equal(t, tt.wantErr, err != nil)
			if tt.errorIs != nil {
				assert.True(t, errors.Is(err, tt.errorIs))
			}
			assert.Equal(t, tt.want, repository.users)
		})
	}
}

func TestInMemoryURLRepository_FindUserById(t *testing.T) {

	const (
		userID1 = "46b8f9d2-b123-4f8e-aabb-f77dd764a00b"
		userID2 = "988cd458-b77e-4d68-b24f-29b1c462d6d3"
	)
	type want struct {
		user    entity.User
		existed bool
	}

	tests := []struct {
		name         string
		existedUsers map[string]entity.User
		userID       string
		want         want
		wantErr      bool
	}{
		{
			name: "item exists",
			existedUsers: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
				},
			},
			userID: userID1,
			want: want{
				user: entity.User{
					UUID: userID1,
				},
				existed: true,
			},
			wantErr: false,
		},
		{
			name: "item not exists",
			existedUsers: map[string]entity.User{
				userID1: entity.User{
					UUID: userID1,
				},
			},
			userID: userID2,
			want: want{
				user:    entity.User{},
				existed: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &InMemoryUserRepository{
				tt.existedUsers,
			}

			user, existed, err := repository.FindUserByID(tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.user, user)
			assert.Equal(t, tt.want.existed, existed)
		})
	}
}
