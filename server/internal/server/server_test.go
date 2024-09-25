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

func TestCollectUsage(t *testing.T) {
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
	_, err := srv.CollectUsage(context.Background(), &v1.CollectUsageRequest{Records: records})
	assert.NoError(t, err)
}
