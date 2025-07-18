package server

import (
	"context"
	"fmt"
	"sort"
	"time"

	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/server/internal/store"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultInterval = time.Hour
	defaultDuration = 7 * 24 * time.Hour
)

// ListModelUsageSummaries lists model usage summaries.
func (s *Server) ListModelUsageSummaries(ctx context.Context, req *v1.ListModelUsageSummariesRequest) (*v1.ListModelUsageSummariesResponse, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user info from context")
	}

	startTime, endTime, err := getStartEndTime(req.Filter, time.Now(), defaultDuration)
	if err != nil {
		return nil, err
	}

	summaries, err := store.ListModelUsageSummaries(
		s.store,
		userInfo.TenantID,
		startTime,
		endTime,
		defaultInterval,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list model usage summaries: %s", err)
	}

	// Find all models. Also filter out summaries belonging to hidden users.
	modelIDMap := make(map[string]struct{})
	var filteredSummaries []*store.ModelUsageSummary
	for _, sum := range summaries {
		if u, ok := s.cache.GetUserByInternalID(sum.UserID); ok && u.Hidden {
			continue
		}
		modelIDMap[sum.ModelID] = struct{}{}
		filteredSummaries = append(filteredSummaries, sum)
	}
	var modelIDs []string
	for modelID := range modelIDMap {
		modelIDs = append(modelIDs, modelID)
	}
	sort.Strings(modelIDs)

	intervalBuckets := make(map[int64]map[string]*store.ModelUsageSummary)
	for _, s := range filteredSummaries {
		m, ok := intervalBuckets[s.TruncatedTimestamp]
		if !ok {
			m = make(map[string]*store.ModelUsageSummary)
			intervalBuckets[s.TruncatedTimestamp] = m
		}
		m[s.ModelID] = s
	}

	var dps []*v1.ListModelUsageSummariesResponse_Datapoint
	for t := startTime; t.Before(endTime); t = t.Add(defaultInterval) {
		sums := intervalBuckets[t.UnixNano()]
		var vs []*v1.ListModelUsageSummariesResponse_Value
		for _, modelID := range modelIDs {
			var v int64
			if sum, ok := sums[modelID]; ok {
				v = sum.TotalRequests
			} else {
				v = 0
			}
			vs = append(vs, &v1.ListModelUsageSummariesResponse_Value{
				ModelId:       modelID,
				TotalRequests: v,
			})
		}

		dps = append(dps, &v1.ListModelUsageSummariesResponse_Datapoint{
			Timestamp: t.Unix(),
			Values:    vs,
		})
	}

	return &v1.ListModelUsageSummariesResponse{
		Datapoints: dps,
	}, nil
}

func getStartEndTime(filter *v1.RequestFilter, now time.Time, duration time.Duration) (time.Time, time.Time, error) {
	if filter == nil {
		filter = &v1.RequestFilter{}
	}

	var (
		startTime time.Time
		endTime   time.Time
	)

	switch t := filter.EndTimestamp; {
	case t > 0:
		endTime = time.Unix(t, 0)
	case t == 0:
		// Set the endtime so that it includes the most recent hour after truncation.
		//
		// But we also don't want to advance if there is no datapoint reported from the agent in
		// the most recent hour. So we add half of the default interval to the current time.
		endTime = now.Add(defaultInterval / 2)
	default:
		return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "endTimestamp must be a non-negative value")
	}
	endTime = endTime.Truncate(defaultInterval)

	switch t := filter.StartTimestamp; {
	case t > 0:
		startTime = time.Unix(t, 0)
	case t == 0:
		startTime = endTime.Add(-1 * duration)
	default:
		return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "startTimestamp must be a non-negative value")
	}
	startTime = startTime.Truncate(defaultInterval)

	if !startTime.Before(endTime) {
		return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "startTimestamp must be before endTimestamp")
	}

	return startTime, endTime, nil
}
