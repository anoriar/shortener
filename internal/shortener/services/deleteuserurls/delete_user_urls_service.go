package deleteuserurls

import (
	"context"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/repository/user"
)

type DeleteUserURLsService struct {
	urlRepository  url.URLRepositoryInterface
	userRepository user.UserRepositoryInterface
}

func NewDeleteUserURLsService(urlRepository url.URLRepositoryInterface, userRepository user.UserRepositoryInterface) *DeleteUserURLsService {
	return &DeleteUserURLsService{urlRepository: urlRepository, userRepository: userRepository}
}

func (service *DeleteUserURLsService) DeleteUserURLs(ctx context.Context, userID string, shortURLs []string) error {
	user, exist, err := service.userRepository.FindUserByID(userID)
	if err != nil {
		return err
	}
	if exist && len(user.SavedURLIDs) > 0 {
		var shortURLsForDelete []string

		for _, requestShortURL := range shortURLs {
			if _, exists := user.SavedURLIDs.FindShortURL(requestShortURL); exists {
				shortURLsForDelete = append(shortURLsForDelete, requestShortURL)
			}
		}
		if len(shortURLsForDelete) > 0 {
			err = service.urlRepository.UpdateIsDeletedBatch(ctx, shortURLsForDelete, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
