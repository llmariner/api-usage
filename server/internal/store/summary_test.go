package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregatedUsage(t *testing.T) {
	st, tearDown := NewTest(t)
	defer tearDown()

	start := int64(1610000000)
	end := int64(1610000200)
	usages := []*Usage{
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
	err := st.CreateUsage(usages...)
	assert.NoError(t, err)

	result, err := st.AggregatedUsage("t0", start, end)
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
