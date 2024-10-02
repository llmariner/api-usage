package sender

import (
	"context"
	"time"

	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Unary is a interceptor that records API usage for gRPC calls with the user info.
func Unary(setter UsageSetter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ui, ok := auth.ExtractUserInfoFromContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		start := time.Now()
		resp, err := handler(ctx, req)
		setter.AddUsage(&v1.UsageRecord{
			UserId:       ui.InternalUserID,
			Tenant:       ui.TenantID,
			Organization: ui.OrganizationID,
			Project:      ui.ProjectID,
			ApiMethod:    info.FullMethod,
			StatusCode:   int32(status.Code(err)),
			Timestamp:    start.UnixNano(),
			LatencyMs:    int32(time.Since(start).Milliseconds()),
		})
		return resp, err
	}
}
