package store

import "gorm.io/gorm"

// Usage represents a API usage.
type Usage struct {
	gorm.Model

	UserID       string `gorm:"uniqueIndex:idx_user_ts"`
	Tenant       string `gorm:"index:idx_tenant_ts"`
	Organization string
	Project      string

	APIKeyID string

	APIMethod  string
	StatusCode int32
	Timestamp  int64 `gorm:"uniqueIndex:idx_user_ts;index:idx_tenant_ts"`
	LatencyMS  int32

	// The following fields are used for chat completions and completions.
	// TODO(kenji): Move these fields to a separate table if needed.
	ModelID            string
	TimeToFirstTokenMS int32
	PromptTokens       int32
	CompletionTokens   int32
}

// CreateUsage creates a new usage.
func (s *Store) CreateUsage(usage ...*Usage) error {
	return s.db.Create(usage).Error
}

// DeleteUsage deletes a usage.
func (s *Store) DeleteUsage(timestamp int64, limit int) (int64, error) {
	res := s.db.Unscoped().Where("timestamp < ?", timestamp).Limit(limit).Delete(&Usage{})
	return res.RowsAffected, res.Error
}

// FindUsages returns the usages. This is used for testing.
func (s *Store) FindUsages() ([]*Usage, error) {
	var usages []*Usage
	if err := s.db.Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

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
func (s *Store) GetUsagesByGroups(tenantID string, start, end int64) ([]*UsageByGroup, error) {
	var us []*UsageByGroup
	if err := s.db.Model(&Usage{}).
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
