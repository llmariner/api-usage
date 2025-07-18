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
func CreateUsage(db *gorm.DB, usage ...*Usage) error {
	return db.Create(usage).Error
}

// FindUsages returns the usages. This is used for testing.
func FindUsages(db *gorm.DB) ([]*Usage, error) {
	var usages []*Usage
	if err := db.Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

// DeleteUsage deletes a usage.
func DeleteUsage(db *gorm.DB, timestamp int64, limit int) (int64, error) {
	res := db.Unscoped().Where("timestamp < ?", timestamp).Limit(limit).Delete(&Usage{})
	return res.RowsAffected, res.Error
}
