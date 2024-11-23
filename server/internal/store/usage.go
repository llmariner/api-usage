package store

import "gorm.io/gorm"

// Usage represents a API usage.
type Usage struct {
	gorm.Model

	UserID       string `gorm:"uniqueIndex:idx_user_ts"`
	Tenant       string
	Organization string
	Project      string

	APIKeyID string

	APIMethod  string
	StatusCode int32
	Timestamp  int64 `gorm:"uniqueIndex:idx_user_ts"`
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
