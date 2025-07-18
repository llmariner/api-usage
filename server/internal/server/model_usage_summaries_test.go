package server

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/google/go-cmp/cmp"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/common/pkg/store"
	"github.com/llmariner/api-usage/server/internal/cache"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestListModelUsageSummaries(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	startTime := time.Date(2025, 7, 17, 11, 0, 0, 0, time.Local)
	endTime := startTime.Add(3 * time.Hour)

	usages := []*store.Usage{
		{
			Tenant:    defaultTenantID,
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: startTime.UnixNano(),
		},
		// different model
		{
			Tenant:    defaultTenantID,
			ModelID:   "model1",
			UserID:    "user0",
			Timestamp: startTime.Add(10 * time.Minute).UnixNano(),
		},
		// hidden user
		{
			Tenant:    defaultTenantID,
			ModelID:   "model1",
			UserID:    "user1",
			Timestamp: startTime.Add(20 * time.Minute).UnixNano(),
		},
		{
			Tenant:    defaultTenantID,
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: startTime.Add(1*time.Hour + 1*time.Minute).UnixNano(),
		},
		{
			Tenant:    defaultTenantID,
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: startTime.Add(1*time.Hour + 30*time.Minute).UnixNano(),
		},
		{
			Tenant:    defaultTenantID,
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: startTime.Add(2*time.Hour + 10*time.Minute).UnixNano(),
		},
	}
	err := store.CreateUsage(st.DB(), usages...)
	assert.NoError(t, err)

	cache := &fakeCache{
		usersByInternalID: map[string]*cache.U{
			"iu0": {
				ID:         "u0",
				InternalID: "iu0",
				Hidden:     false,
			},
			"iu1": {
				ID:         "u1",
				InternalID: "iu1",
				Hidden:     true,
			},
		},
	}
	srv := New(st, cache, testr.New(t))
	ctx := fakeAuthInto(context.Background())

	got, err := srv.ListModelUsageSummaries(ctx, &v1.ListModelUsageSummariesRequest{
		Filter: &v1.RequestFilter{
			StartTimestamp: startTime.Unix(),
			EndTimestamp:   endTime.Unix(),
		},
	})
	assert.NoError(t, err)

	want := &v1.ListModelUsageSummariesResponse{
		Datapoints: []*v1.ListModelUsageSummariesResponse_Datapoint{
			{
				Timestamp: startTime.Unix(),
				Values: []*v1.ListModelUsageSummariesResponse_Value{
					{
						ModelId:       "model0",
						TotalRequests: 1,
					},
					{
						ModelId:       "model1",
						TotalRequests: 1,
					},
				},
			},
			{
				Timestamp: startTime.Add(1 * time.Hour).Unix(),
				Values: []*v1.ListModelUsageSummariesResponse_Value{
					{
						ModelId:       "model0",
						TotalRequests: 2,
					},
					{
						ModelId:       "model1",
						TotalRequests: 0,
					},
				},
			},
			{
				Timestamp: startTime.Add(2 * time.Hour).Unix(),
				Values: []*v1.ListModelUsageSummariesResponse_Value{
					{
						ModelId:       "model0",
						TotalRequests: 1,
					},
					{
						ModelId:       "model1",
						TotalRequests: 0,
					},
				},
			},
		},
	}
	assert.Truef(t, proto.Equal(got, want), cmp.Diff(got, want, protocmp.Transform()))
}

func TestGetStartEndTime(t *testing.T) {
	now := time.Date(2025, 7, 17, 11, 40, 0, 0, time.Local)
	nowT := now.Truncate(time.Hour)

	tcs := []struct {
		name      string
		filter    *v1.RequestFilter
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{
			name:      "no filter",
			filter:    nil,
			wantStart: nowT.Add(-23 * time.Hour),
			wantEnd:   nowT.Add(time.Hour),
		},
		{
			name: "both times zero",
			filter: &v1.RequestFilter{
				StartTimestamp: 0,
				EndTimestamp:   0,
			},
			wantStart: nowT.Add(-23 * time.Hour),
			wantEnd:   nowT.Add(time.Hour),
		},
		{
			name: "both set",
			filter: &v1.RequestFilter{
				StartTimestamp: now.Add(-12 * time.Hour).Unix(),
				EndTimestamp:   now.Add(-6 * time.Hour).Unix(),
			},
			wantStart: nowT.Add(-12 * time.Hour),
			wantEnd:   nowT.Add(-6 * time.Hour),
		},
		{
			name: "start time only",
			filter: &v1.RequestFilter{
				StartTimestamp: now.Add(-12 * time.Hour).Unix(),
				EndTimestamp:   0,
			},
			wantStart: nowT.Add(-12 * time.Hour),
			wantEnd:   nowT.Add(time.Hour),
		},
		{
			name: "end time only",
			filter: &v1.RequestFilter{
				StartTimestamp: 0,
				EndTimestamp:   now.Add(-6 * time.Hour).Unix(),
			},
			wantStart: nowT.Add(-30 * time.Hour),
			wantEnd:   nowT.Add(-6 * time.Hour),
		},
		{
			name: "invalid time range",
			filter: &v1.RequestFilter{
				StartTimestamp: now.Add(-6 * time.Hour).Unix(),
				EndTimestamp:   now.Add(-12 * time.Hour).Unix(),
			},
			wantErr: true,
		},
		{
			name: "negative start time",
			filter: &v1.RequestFilter{
				StartTimestamp: -1,
			},
			wantErr: true,
		},
		{
			name: "negative end time",
			filter: &v1.RequestFilter{
				EndTimestamp: -1,
			},
			wantErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			gotStart, gotEnd, err := getStartEndTime(tc.filter, now, 24*time.Hour)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantStart, gotStart)
			assert.Equal(t, tc.wantEnd, gotEnd)
		})
	}

}
