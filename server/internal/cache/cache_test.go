package cache

import (
	"context"
	"testing"

	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestCache(t *testing.T) {
	ul := &fakeUserInfoLister{
		apikeys: &uv1.ListInternalAPIKeysResponse{
			ApiKeys: []*uv1.InternalAPIKey{
				{
					ApiKey: &uv1.APIKey{
						Id:   "s0",
						Name: "name0",
					},
				},
				{
					ApiKey: &uv1.APIKey{
						Id:   "s1",
						Name: "name1",
					},
					TenantId: "tid1",
				},
			},
		},
		users: &uv1.ListUsersResponse{
			Users: []*uv1.User{
				{
					Id:         "u0",
					InternalId: "iu0",
				},
				{
					Id:         "u1",
					InternalId: "iu1",
				},
			},
		},
	}

	c := NewStore(ul)
	err := c.updateCache(context.Background())
	assert.NoError(t, err)

	wantKeys := map[string]*K{
		"s0": {
			ID:   "s0",
			Name: "name0",
		},
		"s1": {
			ID:   "s1",
			Name: "name1",
		},
	}

	for id, w := range wantKeys {
		got, ok := c.GetAPIKeyByID(id)
		assert.True(t, ok)
		assert.Equal(t, w.ID, got.ID)
		assert.Equal(t, w.Name, got.Name)
	}

	wantUsers := map[string]*U{
		"iu0": {
			ID:         "u0",
			InternalID: "iu0",
		},
		"iu1": {
			ID:         "u1",
			InternalID: "iu1",
		},
	}
	for id, w := range wantUsers {
		got, ok := c.GetUserByInternalID(id)
		assert.True(t, ok)
		assert.Equal(t, w.ID, got.ID)
		assert.Equal(t, w.InternalID, got.InternalID)
	}
}

type fakeUserInfoLister struct {
	apikeys *uv1.ListInternalAPIKeysResponse
	users   *uv1.ListUsersResponse
}

func (l *fakeUserInfoLister) ListInternalAPIKeys(ctx context.Context, in *uv1.ListInternalAPIKeysRequest, opts ...grpc.CallOption) (*uv1.ListInternalAPIKeysResponse, error) {
	return l.apikeys, nil
}

func (l *fakeUserInfoLister) ListUsers(ctx context.Context, in *uv1.ListUsersRequest, opts ...grpc.CallOption) (*uv1.ListUsersResponse, error) {
	return l.users, nil
}
