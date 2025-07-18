package store

import (
	"fmt"
	"time"

	"github.com/llmariner/api-usage/common/pkg/store"
)

// Summary is a struct that represents the summary of the usage data.
type Summary struct {
	APIMethod       string
	TotalRequests   int64
	SuccessRequests int64
	FailureRequests int64
	AverageLatency  float64
}

// AggregatedUsage returns the aggregated usage data for the given tenant and time range.
func AggregatedUsage(s *store.Store, tenantID string, start, end int64) ([]Summary, error) {
	var summaries []Summary
	if err := s.DB().Model(&store.Usage{}).
		Select(
			`api_method`,
			`COUNT(*) AS total_requests`,
			`SUM(CASE WHEN status_code = '200' THEN 1 ELSE 0 END) AS success_requests`,
			`SUM(CASE WHEN status_code != '200' THEN 1 ELSE 0 END) AS failure_requests`,
			`AVG(latency_ms) AS average_latency`).
		Where("tenant = ?", tenantID).
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Group("api_method").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}
	return summaries, nil
}

// ModelUsageSummary is a struct that represents the usage summary for a specific model, user and truncated timestamp.
type ModelUsageSummary struct {
	UserID  string
	ModelID string

	TruncatedTimestamp int64

	TotalRequests int64
}

// ListModelUsageSummaries returns the usage summaries for models grouped by user and truncated by the specified interval.
func ListModelUsageSummaries(
	s *store.Store,
	tenantID string,
	startTime,
	endTime time.Time,
	interval time.Duration,
) ([]*ModelUsageSummary, error) {
	var us []*ModelUsageSummary
	if err := s.DB().Model(&store.Usage{}).
		Select(
			"model_id",
			"user_id",
			// Truncate by interval
			fmt.Sprintf("timestamp / %d * %d AS truncated_timestamp", interval.Nanoseconds(), interval.Nanoseconds()),
			"COUNT(*) AS total_requests",
		).
		Where("tenant = ?", tenantID).
		Where("model_id != ''"). // Exclude non-model related usage
		Where("timestamp >= ? AND timestamp < ?", startTime.UnixNano(), endTime.UnixNano()).
		Group("user_id, model_id, truncated_timestamp").
		Scan(&us).Error; err != nil {
		return nil, err
	}
	return us, nil
}
