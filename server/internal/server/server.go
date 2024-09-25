package server

import (
	"context"
	"fmt"
	"net"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// New creates a new server.
func New(logger logr.Logger) *Server {
	return &Server{
		logger: logger.WithName("grpc"),
	}
}

// Server is the server for the collection service.
type Server struct {
	v1.UnimplementedCollectonServiceServer

	logger logr.Logger
}

// Run starts the gRPC server.
func (s *Server) Run(ctx context.Context, port int) error {
	s.logger.Info("Starting the gRPC server...", "port", port)

	grpcServer := grpc.NewServer()
	v1.RegisterCollectonServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	healthCheck := health.NewServer()
	healthCheck.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthCheck)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen: %s", err)
	}
	if err := grpcServer.Serve(l); err != nil {
		return fmt.Errorf("serve: %s", err)
	}

	s.logger.Info("Stopped gRPC server")
	return nil
}
