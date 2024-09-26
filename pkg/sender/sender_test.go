package sender

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestRun(t *testing.T) {
	client := &fakeCollectonServiceClient{}
	s := &UsageSender{
		client:         client,
		logger:         testr.NewWithOptions(t, testr.Options{Verbosity: 5}),
		interval:       time.Second,
		maxMessageSize: 5,
		usageCh:        make(chan *v1.UsageRecord, defaultUsageChannelSize),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go func() {
		s.AddUsage(&v1.UsageRecord{}) // size=2
		s.AddUsage(&v1.UsageRecord{})
		s.AddUsage(&v1.UsageRecord{})
	}()

	s.Run(ctx)
	assert.Equal(t, 2, client.counter)
	assert.Equal(t, 3, client.totalRecords)
}

type fakeCollectonServiceClient struct {
	counter      int
	totalRecords int
}

func (c *fakeCollectonServiceClient) CreateUsage(ctx context.Context, req *v1.CreateUsageRequest, opts ...grpc.CallOption) (*v1.Usage, error) {
	c.counter++
	c.totalRecords += len(req.Records)
	return &v1.Usage{}, nil
}
