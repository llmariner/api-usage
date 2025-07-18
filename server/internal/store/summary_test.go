package store

import (
	"testing"
	"time"

	"github.com/llmariner/api-usage/common/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestAggregatedUsage(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	start := int64(1610000000)
	end := int64(1610000200)
	usages := []*store.Usage{
		{
			Tenant:     "t0",
			APIMethod:  "GetFoo",
			StatusCode: 200,
			Timestamp:  start,
			LatencyMS:  100,
		},
		{
			Tenant:     "t0",
			APIMethod:  "GetFoo",
			StatusCode: 200,
			Timestamp:  start + 50,
			LatencyMS:  150,
		},
		{
			Tenant:     "t0",
			APIMethod:  "GetFoo",
			StatusCode: 404,
			Timestamp:  start + 100,
			LatencyMS:  200,
		},
		{
			Tenant:     "t0",
			APIMethod:  "UpdateFoo",
			StatusCode: 200,
			Timestamp:  start + 150,
			LatencyMS:  300,
		},
		{
			Tenant:     "t0",
			APIMethod:  "UpdateFoo",
			StatusCode: 200,
			Timestamp:  start + 190,
			LatencyMS:  100,
		},
		// out of range
		{
			Tenant:     "t0",
			APIMethod:  "UpdateFoo",
			StatusCode: 200,
			Timestamp:  start + 200,
			LatencyMS:  100,
		},
		// different tenant
		{
			Tenant:     "t1",
			APIMethod:  "GetFoo",
			StatusCode: 404,
			Timestamp:  start + 130,
			LatencyMS:  200,
		},
	}
	err := store.CreateUsage(st.DB(), usages...)
	assert.NoError(t, err)

	result, err := AggregatedUsage(st, "t0", start, end)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, "GetFoo", result[0].APIMethod)
	assert.Equal(t, int64(3), result[0].TotalRequests)
	assert.Equal(t, int64(2), result[0].SuccessRequests)
	assert.Equal(t, int64(1), result[0].FailureRequests)
	assert.Equal(t, float64(150), result[0].AverageLatency) // avg(100 + 150 + 200)

	assert.Equal(t, "UpdateFoo", result[1].APIMethod)
	assert.Equal(t, int64(2), result[1].TotalRequests)
	assert.Equal(t, int64(2), result[1].SuccessRequests)
	assert.Equal(t, int64(0), result[1].FailureRequests)
	assert.Equal(t, float64(200), result[1].AverageLatency) // avg(300 + 100)
}

func TestListModelUsageSummaries(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	ts := time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC)
	usages := []*store.Usage{
		{
			Tenant:    "t0",
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: ts.Unix(),
		},
		// differnet model
		{
			Tenant:    "t0",
			ModelID:   "model1",
			UserID:    "user0",
			Timestamp: ts.Add(10 * time.Minute).Unix(),
		},
		// different user
		{
			Tenant:    "t0",
			ModelID:   "model1",
			UserID:    "user1",
			Timestamp: ts.Add(20 * time.Minute).Unix(),
		},
		// different tenant
		{
			Tenant:    "t1",
			ModelID:   "model2",
			UserID:    "user2",
			Timestamp: ts.Add(20 * time.Minute).Unix(),
		},
		// different timestamp
		{
			Tenant:    "t0",
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: ts.Add(1*time.Hour + 1*time.Minute).Unix(),
		},
		{
			Tenant:    "t0",
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: ts.Add(1*time.Hour + 30*time.Minute).Unix(),
		},
		{
			Tenant:    "t0",
			ModelID:   "model0",
			UserID:    "user0",
			Timestamp: ts.Add(2*time.Hour + 10*time.Minute).Unix(),
		},
	}
	err := store.CreateUsage(st.DB(), usages...)
	assert.NoError(t, err)

	startTime := ts
	endTime := ts.Add(24 * time.Hour)
	got, err := ListModelUsageSummaries(st, "t0", startTime, endTime, time.Hour)
	assert.NoError(t, err)

	want := []*ModelUsageSummary{
		{
			ModelID:            "model0",
			UserID:             "user0",
			TruncatedTimestamp: ts.Unix(),
			TotalRequests:      1,
		},
		{
			ModelID:            "model0",
			UserID:             "user0",
			TruncatedTimestamp: ts.Add(1 * time.Hour).Unix(),
			TotalRequests:      2,
		},
		{
			ModelID:            "model0",
			UserID:             "user0",
			TruncatedTimestamp: ts.Add(2 * time.Hour).Unix(),
			TotalRequests:      1,
		},
		{
			ModelID:            "model1",
			UserID:             "user0",
			TruncatedTimestamp: ts.Unix(),
			TotalRequests:      1,
		},
		{
			ModelID:            "model1",
			UserID:             "user1",
			TruncatedTimestamp: ts.Unix(),
			TotalRequests:      1,
		},
	}
	assert.ElementsMatch(t, want, got)

}
