package cleaner

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/llmariner/api-usage/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestClearUsage(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	dur := []time.Duration{-time.Hour, -time.Minute, -time.Second}
	for i := 0; i < 3; i++ {
		record := &store.Usage{
			UserID:    fmt.Sprintf("u%d", i),
			Timestamp: time.Now().Add(dur[i]).UnixNano(),
		}
		err := store.CreateUsage(st.DB(), record)
		assert.NoError(t, err)
	}

	retentionPeriod := time.Second * 5
	interval := time.Second * 10
	cleaner := NewCleaner(st.DB(), retentionPeriod, interval, testr.New(t))

	err := cleaner.clearUsage()
	assert.NoError(t, err)
	got, err := store.FindUsages(st.DB())
	assert.NoError(t, err)
	assert.Len(t, got, 1)

	time.Sleep(interval)

	err = cleaner.clearUsage()
	assert.NoError(t, err)
	got, err = store.FindUsages(st.DB())
	assert.NoError(t, err)
	assert.Len(t, got, 0)
}
