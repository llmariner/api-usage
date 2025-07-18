package server

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/common/pkg/store"
	"github.com/llmariner/api-usage/server/internal/cache"
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

			UserID:   "iu0",
			APIKeyID: "api_key_id0",
			ModelID:  "model0",

			TimeToFirstTokenMS: 10,
			PromptTokens:       100,
			CompletionTokens:   100,
		},
	}
	err := store.CreateUsage(st.DB(), usages...)
	assert.NoError(t, err)

	cache := &fakeCache{
		apiKeysByID: map[string]*cache.K{
			"api_key_id0": {
				ID:   "api_key_id0",
				Name: "api_key0",
			},
		},
		usersByInternalID: map[string]*cache.U{
			"iu0": {
				ID:         "u0",
				InternalID: "iu0",
			},
		},
	}
	srv := New(st, cache, testr.New(t))
	ctx := fakeAuthInto(context.Background())

	exp := &v1.ListUsageDataResponse{
		Usages: []*v1.UsageDataByGroup{
			{
				UserId:     "u0",
				ApiKeyId:   "api_key_id0",
				ApiKeyName: "api_key0",
				ModelId:    "model0",

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

type fakeCache struct {
	apiKeysByID       map[string]*cache.K
	usersByInternalID map[string]*cache.U
}

func (c *fakeCache) GetAPIKeyByID(id string) (*cache.K, bool) {
	k, ok := c.apiKeysByID[id]
	return k, ok
}

func (c *fakeCache) GetUserByInternalID(internalID string) (*cache.U, bool) {
	u, ok := c.usersByInternalID[internalID]
	return u, ok
}
