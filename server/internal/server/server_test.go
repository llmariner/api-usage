package server

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/server/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestListUsageData(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	start := int64(1610000000)
	end := int64(1610000200)
	usages := []*store.Usage{
		{
			Tenant:     defaultTenantID,
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
	}
	err := st.CreateUsage(usages...)
	assert.NoError(t, err)

	srv := New(st, testr.New(t))
	ctx := fakeAuthInto(context.Background())

	exp := &v1.ListUsageDataResponse{
		Usages: []*v1.UsageDataByGroup{
			{
				ApiKeyId: "api_key0",
				UserId:   "u0",
				ModelId:  "model0",

				TotalRequests:         1,
				TotalPromptTokens:     100,
				TotalCompletionTokens: 100,
				AvgLatencyMs:          100,
				AvgTimeToFirstTokenMs: 10,
			},
		},
	}
	got, err := srv.ListUsageData(ctx, &v1.ListUsageDataRequest{
		StartTime: start,
		EndTime:   end,
	})
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}
