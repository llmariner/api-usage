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

// NewInternal creates a new internal server.
func NewInternal(store *store.Store, logger logr.Logger) *InternalServer {
	return &InternalServer{
		store:  store,
		logger: logger.WithName("internal"),
	}
}

// InternalServer is the server for the collection service.
type InternalServer struct {
	v1.UnimplementedCollectionInternalServiceServer

	store  *store.Store
	logger logr.Logger
}

// Run starts the internal gRPC server.
func (s *InternalServer) Run(ctx context.Context, port int) error {
	s.logger.Info("Starting the internal server...", "port", port)

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

	s.logger.Info("Stopped internal server")
	return nil
}

// CreateUsage creates usage.
func (s *InternalServer) CreateUsage(ctx context.Context, req *v1.CreateUsageRequest) (*v1.Usage, error) {
	s.logger.V(4).WithName("api").Info("CreateUsage", "count", len(req.Records))
	// TODO: add authentication

	var records []*store.Usage
	for _, r := range req.Records {
		records = append(records, &store.Usage{
			UserID:       r.UserId,
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

	return &v1.Usage{}, nil
}
