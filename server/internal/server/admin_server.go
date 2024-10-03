package server

import (
	"context"
	"fmt"
	"net"
	"time"

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

// NewAdmin creates a new admin server.
func NewAdmin(store *store.Store, logger logr.Logger) *AdminServer {
	return &AdminServer{
		store:  store,
		logger: logger.WithName("admin"),
	}
}

// AdminServer is the server for the api-usage service.
type AdminServer struct {
	v1.UnimplementedAPIUsageServiceServer

	store  *store.Store
	logger logr.Logger
}

// Run starts the admin server.
func (s *AdminServer) Run(ctx context.Context, port int) error {
	s.logger.Info("Starting the admin server...", "port", port)

	grpcServer := grpc.NewServer()
	v1.RegisterAPIUsageServiceServer(grpcServer, s)
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

	s.logger.Info("Stopped admin server")
	return nil
}

// GetAggregatedSummary aggregates the API usage in the specified time range.
func (s *AdminServer) GetAggregatedSummary(ctx context.Context, req *v1.GetAggregatedSummaryRequest) (*v1.AggregatedSummary, error) {
	s.logger.V(4).Info("GetAggregatedSummary", "request", req)

	st := time.Now()

	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartTime == 0 {
		req.StartTime = st.Add(-24 * time.Hour).UnixNano()
	}
	if req.EndTime == 0 {
		req.EndTime = st.UnixNano()
	}

	summary, err := s.store.AggregatedUsage(req.TenantId, req.StartTime, req.EndTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cumulativeSum float64
	total := &v1.Summary{}
	methods := make([]*v1.Summary, len(summary))
	for _, s := range summary {
		total.TotalRequests += s.TotalRequests
		total.SuccessRequests += s.SuccessRequests
		total.FailureRequests += s.FailureRequests
		cumulativeSum += s.AverageLatency * float64(s.TotalRequests)

		methods = append(methods, &v1.Summary{
			Method:          s.APIMethod,
			TotalRequests:   s.TotalRequests,
			SuccessRequests: s.SuccessRequests,
			FailureRequests: s.FailureRequests,
			AverageLatency:  s.AverageLatency,
		})
	}
	if total.TotalRequests > 0 {
		total.AverageLatency = cumulativeSum / float64(total.TotalRequests)
	}

	return &v1.AggregatedSummary{
		Summary:         total,
		MethodSummaries: methods,
	}, nil
}
