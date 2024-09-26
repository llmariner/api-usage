package server

import (
	"context"
	"fmt"
	"net"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/server/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// New creates a new server.
func New(store *store.Store, logger logr.Logger) *Server {
	return &Server{
		store:  store,
		logger: logger.WithName("grpc"),
	}
}

// Server is the server for the collection service.
type Server struct {
	v1.UnimplementedCollectionInternalServiceServer

	store  *store.Store
	logger logr.Logger
}

// Run starts the gRPC server.
func (s *Server) Run(ctx context.Context, port int) error {
	s.logger.Info("Starting the gRPC server...", "port", port)

	grpcServer := grpc.NewServer()
	v1.RegisterCollectionInternalServiceServer(grpcServer, s)
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

// CreateUsage creates usage.
func (s *Server) CreateUsage(ctx context.Context, req *v1.CreateUsageRequest) (*v1.CreateUsageResponse, error) {
	s.logger.V(4).WithName("api").Info("CreateUsage", "count", len(req.Records))
	// TODO: add authentication

	var records []*store.Usage
	for _, r := range req.Records {
		records = append(records, &store.Usage{
			User:         r.User,
			Tenant:       r.Tenant,
			Organization: r.Organization,
			Project:      r.Project,
			APIMethod:    r.ApiMethod,
			StatusCode:   r.StatusCode,
			Timestamp:    r.Timestamp,
			LatencyMS:    r.LatencyMs,
		})
	}
	if err := s.store.CreateUsage(records...); err != nil {
		return nil, status.Errorf(codes.Internal, "create usage: %s", err)
	}

	return &v1.CreateUsageResponse{}, nil
}
