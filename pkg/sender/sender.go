package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/api-usage/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// UsageSetter is an interface that allows components to add usage data to the sender.
type UsageSetter interface {
	AddUsage(usage *v1.UsageRecord)
}

// New creates a new UsageSender.
func New(ctx context.Context, c Config, opt grpc.DialOption, log logr.Logger) (*UsageSender, error) {
	cc, err := grpc.NewClient(c.ServerAddr, opt)
	if err != nil {
		return nil, fmt.Errorf("create client: %s", err)
	}
	return &UsageSender{
		client:         v1.NewCollectionInternalServiceClient(cc),
		logger:         log.WithName("usage"),
		initialDelay:   c.InitialDelay,
		interval:       c.Interval,
		maxMessageSize: c.MaxMessageSize,
		usageCh:        make(chan *v1.UsageRecord, c.UsageChannelSize),
	}, nil
}

// UsageSender is a component that sends usage data to the API server.
type UsageSender struct {
	client v1.CollectionInternalServiceClient
	logger logr.Logger

	initialDelay   time.Duration
	interval       time.Duration
	maxMessageSize int

	usageCh chan *v1.UsageRecord
}

// Run starts the usage sender. It sends usage data to the usage collector.
// If the maximum message size is exceeded, the sender will send the data immediately.
func (s *UsageSender) Run(ctx context.Context) {
	s.logger.Info("Starting usage sender...", "interval", s.interval, "delay", s.initialDelay)
	time.Sleep(s.initialDelay)

	send := func(records []*v1.UsageRecord) {
		if err := s.sendUsageData(ctx, records); err != nil {
			s.logger.Error(err, "Failed to send usage data")
		}
	}

	var buffer []*v1.UsageRecord
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case record := <-s.usageCh:
			buffer = append(buffer, record)
			if size := proto.Size(&v1.CollectUsageRequest{Records: buffer}); size > s.maxMessageSize {
				s.logger.V(1).Info("Max message size exceeded", "size", size, "count", len(buffer))
				records := buffer[:len(buffer)-1]
				buffer = []*v1.UsageRecord{record}
				send(records)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				records := buffer
				buffer = make([]*v1.UsageRecord, 0)
				send(records)
			}
		case <-ctx.Done():
			s.logger.Info("Stopping usage sender...")
			if len(buffer) > 0 {
				send(buffer)
			}
			s.logger.Info("Stopped usage sender")
			return
		}
	}
}

// AddUsage adds a usage record to the sender.
func (s *UsageSender) AddUsage(usage *v1.UsageRecord) {
	select {
	case s.usageCh <- usage:
	default:
		s.logger.Error(nil, "Dropped usage record", "record", usage)
	}
}

func (s *UsageSender) sendUsageData(ctx context.Context, records []*v1.UsageRecord) error {
	req := &v1.CollectUsageRequest{Records: records}
	if _, err := s.client.CollectUsage(ctx, req); err != nil {
		return fmt.Errorf("collect usage: %s", err)
	}
	s.logger.V(4).Info("Sent API usage", "count", len(records))
	return nil
}
