package user

import (
	"github.com/anoriar/shortener/internal/shortener/repository/user"
)

type UserService struct {
	userRepository user.UserRepositoryInterface
}

func NewUserService(userRepository user.UserRepositoryInterface) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) AddShortURLsToUser(userID string, shortURLs []string) error {
	user, exists, err := us.userRepository.FindUserByID(userID)
	if err != nil {
		return err
	}
	if exists {
		user.SavedURLIDs = append(user.SavedURLIDs, shortURLs...)
		err := us.userRepository.UpdateUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *UserService) GetUserShortURLs(userID string) ([]string, error) {
	user, exists, err := us.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return user.SavedURLIDs, nil
	}
	return []string{}, nil
}
