package cache

import (
	"context"
	"log"
	"sync"
	"time"

	uv1 "github.com/llmariner/user-manager/api/v1"
	"google.golang.org/grpc"
)

// K represents an API key.
type K struct {
	ID   string
	Name string
}

// U represents a user.
type U struct {
	ID         string
	InternalID string
}

type userInfoLister interface {
	ListInternalAPIKeys(ctx context.Context, in *uv1.ListInternalAPIKeysRequest, opts ...grpc.CallOption) (*uv1.ListInternalAPIKeysResponse, error)
	ListUsers(ctx context.Context, in *uv1.ListUsersRequest, opts ...grpc.CallOption) (*uv1.ListUsersResponse, error)
}

// NewStore creates a new cache store.
func NewStore(
	userInfoLister userInfoLister,
) *Store {
	return &Store{
		userInfoLister: userInfoLister,

		apiKeysByID:       map[string]*K{},
		usersByInternalID: map[string]*U{},
	}
}

// Store is a cache for API keys and organization users.
type Store struct {
	userInfoLister userInfoLister

	// apiKeysByID is a set of API keys, keyed by its ID.
	apiKeysByID map[string]*K

	// usersByInternalID is a set of users, keyed by its internal ID.
	usersByInternalID map[string]*U

	lastSuccessfulSyncTime time.Time

	mu sync.RWMutex
}

// GetAPIKeyByID returns an API key by its ID.
func (c *Store) GetAPIKeyByID(id string) (*K, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	k, ok := c.apiKeysByID[id]
	return k, ok
}

// GetUserByInternalID returns a user by its internal ID.
func (c *Store) GetUserByInternalID(internalID string) (*U, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	u, ok := c.usersByInternalID[internalID]
	return u, ok
}

// GetLastSuccessfulSyncTime returns the last successful sync time.
func (c *Store) GetLastSuccessfulSyncTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastSuccessfulSyncTime
}

// Sync synchronizes the cache.
func (c *Store) Sync(ctx context.Context, interval time.Duration) error {
	if err := c.updateCache(ctx); err != nil {
		// Gracefully ignore the error.
		log.Printf("Failed to update the cache: %s. Ignoring.", err)
	}

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := c.updateCache(ctx); err != nil {
				// Gracefully ignore the error.
				log.Printf("Failed to update the cache: %s. Ignoring.", err)
			}
		}
	}
}

func (c *Store) updateCache(ctx context.Context) error {
	resp, err := c.userInfoLister.ListInternalAPIKeys(ctx, &uv1.ListInternalAPIKeysRequest{})
	if err != nil {
		return err
	}

	k := map[string]*K{}
	for _, apiKey := range resp.ApiKeys {
		k[apiKey.ApiKey.Id] = &K{
			ID:   apiKey.ApiKey.Id,
			Name: apiKey.ApiKey.Name,
		}
	}

	uresp, err := c.userInfoLister.ListUsers(ctx, &uv1.ListUsersRequest{})
	if err != nil {
		return err
	}

	u := map[string]*U{}
	for _, user := range uresp.Users {
		u[user.InternalId] = &U{
			ID:         user.Id,
			InternalID: user.InternalId,
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.apiKeysByID = k
	c.usersByInternalID = u

	c.lastSuccessfulSyncTime = time.Now()

	return nil
}
