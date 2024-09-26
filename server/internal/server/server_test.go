package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/server/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateUsage(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	var records []*v1.UsageRecord
	for i := 0; i < 3; i++ {
		records = append(records, &v1.UsageRecord{
			User:      fmt.Sprintf("u%d", i),
			Timestamp: time.Now().UnixNano(),
		})
	}
	srv := New(st, testr.New(t))
	ctx := context.Background()

	_, err := srv.CreateUsage(ctx, &v1.CreateUsageRequest{Records: records})
	assert.NoError(t, err)
	_, err = srv.CreateUsage(ctx, &v1.CreateUsageRequest{Records: records})
	assert.Error(t, err)
}
