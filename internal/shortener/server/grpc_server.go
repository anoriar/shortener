package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/anoriar/shortener/internal/app"
	grpcserver "github.com/anoriar/shortener/internal/shortener/handlers/grpc"
	grpcinterceptor "github.com/anoriar/shortener/internal/shortener/middleware/auth/grpc"

	pb "github.com/anoriar/shortener/proto/generated/shortener/proto"
)

// GRPCServer missing godoc.
type GRPCServer struct {
	app             *app.App
	authInterceptor *grpcinterceptor.AuthInterceptor
}

// NewGRPCServer missing godoc.
func NewGRPCServer(app *app.App) *GRPCServer {
	return &GRPCServer{app: app, authInterceptor: grpcinterceptor.NewAuthInterceptor(app.Authenticator)}
}

// RunGRPCServer missing godoc.
// RunServer missing godoc.
func (grpcServer *GRPCServer) RunGRPCServer() error {
	listen, err := net.Listen("tcp", grpcServer.app.Config.Host)
	if err != nil {
		log.Printf("Error listen: %v\n", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcServer.authInterceptor.Auth))

	pb.RegisterURLServiceServer(s, grpcserver.NewURLServiceServer(*grpcServer.app))
	pb.RegisterStatsServiceServer(s, grpcserver.NewStatsServiceServer(*grpcServer.app))

	fmt.Println("Сервер gRPC начал работу")

	if err := s.Serve(listen); err != grpc.ErrServerStopped {
		log.Printf("Error starting the server: %v\n", err)
	}
	return nil
}
