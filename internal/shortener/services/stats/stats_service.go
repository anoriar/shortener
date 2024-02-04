package stats

import (
	"context"
	"fmt"

	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/repository/user"
)

// StatsService missing godoc.
type StatsService struct {
	urlRepository  url.URLRepositoryInterface
	userRepository user.UserRepositoryInterface
}

// NewStatsService missing godoc.
func NewStatsService(urlRepository url.URLRepositoryInterface, userRepository user.UserRepositoryInterface) *StatsService {
	return &StatsService{urlRepository: urlRepository, userRepository: userRepository}
}

// GetStats missing godoc.
func (s *StatsService) GetStats(ctx context.Context) (*response.StatsDto, error) {
	urlsCount, err := s.urlRepository.GetAllURLsCount(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetStats: %v", err)
	}
	usersCount, err := s.userRepository.GetAllUsersCount(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetStats: %v", err)
	}

	return &response.StatsDto{
		URLs:  urlsCount,
		Users: usersCount,
	}, nil
}
