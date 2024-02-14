package stats

import (
	"context"

	"github.com/anoriar/shortener/internal/shortener/dto/response"
)

// StatsServiceInterface missing godoc.
type StatsServiceInterface interface {
	GetStats(ctx context.Context) (*response.StatsDto, error)
}
