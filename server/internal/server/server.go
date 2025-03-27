package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/server/internal/cache"
	"github.com/llmariner/api-usage/server/internal/config"
	"github.com/llmariner/api-usage/server/internal/store"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	defaultOrgnizationID = "default"
	defaultProjectID     = "default"
	defaultClusterID     = "default"
	defaultTenantID      = "default-tenant-id"
	defaultNamespace     = "default"
)

type cacheGetter interface {
	GetAPIKeyByID(id string) (*cache.K, bool)
	GetUserByInternalID(internalID string) (*cache.U, bool)
}

// New creates a new server.
func New(store *store.Store, cache cacheGetter, logger logr.Logger) *Server {
	return &Server{
		store:  store,
		cache:  cache,
		logger: logger.WithName("server"),
	}
}

// Server is the server for the api-usage service.
type Server struct {
	v1.UnimplementedAPIUsageServiceServer

	srv    *grpc.Server
	store  *store.Store
	cache  cacheGetter
	logger logr.Logger
}

// Run starts the server.
func (s *Server) Run(ctx context.Context, port int, authConfig config.AuthConfig) error {
	s.logger.Info("Starting gRPC server...", "port", port)

	var opt grpc.ServerOption
	if authConfig.Enable {
		ai, err := auth.NewInterceptor(ctx, auth.Config{
			RBACServerAddr: authConfig.RBACInternalServerAddr,
			AccessResource: "api.api_usages",
		})
		if err != nil {
			return err
		}
		opt = grpc.ChainUnaryInterceptor(ai.Unary("/grpc.health.v1.Health/Check"))
	} else {
		fakeAuth := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return handler(fakeAuthInto(ctx), req)
		}
		opt = grpc.ChainUnaryInterceptor(fakeAuth)
	}

	grpcServer := grpc.NewServer(opt)
	v1.RegisterAPIUsageServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.srv = grpcServer

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
	return nil
}

// fakeAuthInto sets dummy user info and token into the context.
func fakeAuthInto(ctx context.Context) context.Context {
	// Set dummy user info and token
	ctx = auth.AppendUserInfoToContext(ctx, auth.UserInfo{
		OrganizationID: defaultOrgnizationID,
		ProjectID:      defaultProjectID,
		AssignedKubernetesEnvs: []auth.AssignedKubernetesEnv{
			{
				ClusterID: defaultClusterID,
				Namespace: defaultNamespace,
			},
		},
		TenantID: defaultTenantID,
	})
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", "Bearer token"))
	return ctx
}

// Stop stops the gRPC server.
func (s *Server) Stop() {
	s.logger.Info("Stopped server")
	s.srv.Stop()
}

// ListUsageData list the usages by groups.
func (s *Server) ListUsageData(ctx context.Context, req *v1.ListUsageDataRequest) (*v1.ListUsageDataResponse, error) {
	s.logger.V(4).Info("ListUsageData", "request", req)

	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user info from context")
	}

	now := time.Now()
	if req.StartTime == 0 {
		req.StartTime = now.Add(-31 * 24 * time.Hour).UnixNano()
	}
	if req.EndTime == 0 {
		req.EndTime = now.UnixNano()
	}
	if req.EndTime <= req.StartTime {
		return nil, status.Error(codes.InvalidArgument, "end time is before start time")
	}

	us, err := s.store.GetUsagesByGroups(userInfo.TenantID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	usProto := make([]*v1.UsageDataByGroup, len(us))
	for i, u := range us {
		usProto[i] = s.usageByGroupToProto(u)
	}

	return &v1.ListUsageDataResponse{
		Usages: usProto,
	}, nil
}

func (s *Server) usageByGroupToProto(ubg *store.UsageByGroup) *v1.UsageDataByGroup {
	var userID string
	u, ok := s.cache.GetUserByInternalID(ubg.UserID)
	if ok {
		userID = u.ID
	} else {
		userID = "unknown"
	}

	var apiKeyName string
	if ubg.APIKeyID != "" {
		k, ok := s.cache.GetAPIKeyByID(ubg.APIKeyID)
		if ok {
			apiKeyName = k.Name
		} else {
			apiKeyName = "unknown"
		}
	}

	return &v1.UsageDataByGroup{
		UserId:                userID,
		ApiKeyId:              ubg.APIKeyID,
		ApiKeyName:            apiKeyName,
		ModelId:               ubg.ModelID,
		TotalRequests:         ubg.TotalRequests,
		TotalPromptTokens:     ubg.TotalPromptTokens,
		TotalCompletionTokens: ubg.TotalCompletionTokens,
		AvgLatencyMs:          ubg.AverageLatency,
		AvgTimeToFirstTokenMs: ubg.AverageTimeToFirstToken,
	}
}
