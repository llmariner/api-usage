package store

import "gorm.io/gorm"

// Usage represents a API usage.
type Usage struct {
	gorm.Model

	UserID       string `gorm:"uniqueIndex:idx_user_ts"`
	Tenant       string
	Organization string
	Project      string

	APIMethod  string
	StatusCode int32
	Timestamp  int64 `gorm:"uniqueIndex:idx_user_ts"`
	LatencyMS  int32
}

// CreateUsage creates a new usage.
func (s *Store) CreateUsage(usage ...*Usage) error {
	return s.db.Create(usage).Error
}
