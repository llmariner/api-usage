package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsagesByGroups(t *testing.T) {
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

			UserID:   "u0",
			APIKeyID: "api_key0",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
		{
			Tenant:     "t0",
			APIMethod:  "GetFoo",
			StatusCode: 200,
			Timestamp:  start + 50,
			LatencyMS:  150,

			UserID:   "u0",
			APIKeyID: "api_key0",
			ModelID:  "model1",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
		{
			Tenant:     "t0",
			APIMethod:  "GetFoo",
			StatusCode: 404,
			Timestamp:  start + 100,
			LatencyMS:  200,

			UserID:   "u1",
			APIKeyID: "api_key0",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
		{
			Tenant:     "t0",
			APIMethod:  "UpdateFoo",
			StatusCode: 200,
			Timestamp:  start + 150,
			LatencyMS:  300,

			UserID:   "u1",
			APIKeyID: "api_key1",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
		{
			Tenant:     "t0",
			APIMethod:  "UpdateFoo",
			StatusCode: 200,
			Timestamp:  start + 190,
			LatencyMS:  100,

			UserID:   "u1",
			APIKeyID: "api_key1",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
		// out of range
		{
			Tenant:     "t0",
			APIMethod:  "UpdateFoo",
			StatusCode: 200,
			Timestamp:  start + 200,
			LatencyMS:  100,

			UserID:   "u0",
			APIKeyID: "api_key0",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       1000,
			CompletionTokens:   1000,
		},
		// different tenant
		{
			Tenant:     "t1",
			APIMethod:  "GetFoo",
			StatusCode: 404,
			Timestamp:  start + 130,
			LatencyMS:  200,

			UserID:   "u10",
			APIKeyID: "api_key10",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
		// empty model ID
		{
			Tenant:     "t0",
			APIMethod:  "GetFoo",
			StatusCode: 200,
			Timestamp:  start + 201,
			LatencyMS:  200,

			UserID:   "u10",
			APIKeyID: "api_key10",
			ModelID:  "",
		},
	}
	err := st.CreateUsage(usages...)
	assert.NoError(t, err)

	result, err := st.GetUsagesByGroups("t0", start, end)
	assert.NoError(t, err)
	assert.Len(t, result, 4)

	us := []*UsageByGroup{
		{
			APIKeyID: "api_key0",
			UserID:   "u0",
			ModelID:  "model0",

			TotalRequests:           1,
			TotalPromptTokens:       100,
			TotalCompletionTokens:   100,
			AverageLatency:          100,
			AverageTimeToFirstToken: 10,
		},
		{
			APIKeyID: "api_key0",
			UserID:   "u0",
			ModelID:  "model1",

			TotalRequests:           1,
			TotalPromptTokens:       100,
			TotalCompletionTokens:   100,
			AverageLatency:          150,
			AverageTimeToFirstToken: 10,
		},
		{
			APIKeyID: "api_key0",
			UserID:   "u1",
			ModelID:  "model0",

			TotalRequests:           1,
			TotalPromptTokens:       100,
			TotalCompletionTokens:   100,
			AverageLatency:          200,
			AverageTimeToFirstToken: 10,
		},
		{
			APIKeyID: "api_key1",
			UserID:   "u1",
			ModelID:  "model0",

			TotalRequests:           2,
			TotalPromptTokens:       200,
			TotalCompletionTokens:   200,
			AverageLatency:          200,
			AverageTimeToFirstToken: 10,
		},
	}

	for i, res := range result {
		assert.Equal(t, us[i], res)
	}
}
