package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anoriar/shortener/internal/app"
	"github.com/anoriar/shortener/internal/shortener/services/stats"
	pb "github.com/anoriar/shortener/proto/generated/shortener/proto"
)

// StatsServiceServer missing godoc.
type StatsServiceServer struct {
	pb.UnimplementedStatsServiceServer

	statsService stats.StatsServiceInterface
	logger       *zap.Logger
}

// NewStatsServiceServer missing godoc.
// StatsServiceServer missing godoc.
func NewStatsServiceServer(
	app app.App,
) *StatsServiceServer {
	return &StatsServiceServer{
		statsService: app.StatsService,
		logger:       app.Logger,
	}
}

// GetStats возвращает статистику сервера
// на выход:
//
//	{
//	  "urls": 3,
//	  "users": 2
//	},
func (service StatsServiceServer) GetStats(ctx context.Context, empty *pb.Empty) (*pb.StatsResponse, error) {
	var response *pb.StatsResponse

	result, err := service.statsService.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, `internal error`)
	}

	response.Urls = int32(result.URLs)
	response.Users = int32(result.Users)

	return response, nil
}
