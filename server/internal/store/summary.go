package store

import (
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
