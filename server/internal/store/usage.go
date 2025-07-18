package store

import "github.com/llmariner/api-usage/common/pkg/store"

// UsageByGroup represents the aggregated usage data grouped by API key, user, and model.
type UsageByGroup struct {
	APIKeyID string
	UserID   string
	ModelID  string

	TotalRequests           int64
	TotalPromptTokens       int64
	TotalCompletionTokens   int64
	AverageLatency          float64
	AverageTimeToFirstToken float64
}

// GetUsagesByGroups returns the aggregated usage data by groups for the given tenant and time range.
func GetUsagesByGroups(s *store.Store, tenantID string, start, end int64) ([]*UsageByGroup, error) {
	var us []*UsageByGroup
	if err := s.DB().Model(&store.Usage{}).
		Select(
			`api_key_id`, `user_id`, `model_id`,
			`COUNT(*) AS total_requests`,
			`AVG(time_to_first_token_ms) AS average_time_to_first_token`,
			`SUM(prompt_tokens) AS total_prompt_tokens`,
			`SUM(completion_tokens) AS total_completion_tokens`,
			`AVG(latency_ms) AS average_latency`).
		Where("tenant = ?", tenantID).
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Where("model_id <> ?", "").
		Group("api_key_id,user_id,model_id").
		Order("api_key_id,user_id,model_id").
		Scan(&us).Error; err != nil {
		return nil, err
	}
	return us, nil
}
