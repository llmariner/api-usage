package cleaner

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/llmariner/api-usage/common/pkg/store"
	"gorm.io/gorm"
)

// Cleaner is the struct that deletes records outside the retention period.
type Cleaner struct {
	db              *gorm.DB
	retentionPeriod time.Duration
	ticker          *time.Ticker
	logger          logr.Logger
}

// NewCleaner returns a new cleaner.
func NewCleaner(db *gorm.DB, retentionPeriod, interval time.Duration, logger logr.Logger) *Cleaner {
	return &Cleaner{
		db:              db,
		retentionPeriod: retentionPeriod,
		ticker:          time.NewTicker(interval),
		logger:          logger.WithName("cleaner"),
	}
}

// Run runs the poller for the Cleaner.
func (r *Cleaner) Run(ctx context.Context) error {
	return r.runDeletion(ctx, r.clearUsage)
}

func (r *Cleaner) runDeletion(ctx context.Context, f func() error) error {
	if err := f(); err != nil {
		return err
	}

	for {
		select {
		case <-r.ticker.C:
			if err := f(); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *Cleaner) clearUsage() error {
	t := time.Now().Add(-r.retentionPeriod).UnixNano()
	const limit = 100
	for {
		deleted, err := store.DeleteUsage(r.db, t, limit)
		if err != nil {
			return err
		}
		r.logger.Info("Deleted usage", "records", deleted)
		if deleted == 0 {
			break
		}
	}

	return nil
}
