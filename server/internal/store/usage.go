package store

import "gorm.io/gorm"

// Usage represents a API usage.
type Usage struct {
	gorm.Model

	User         string
	Tenant       string
	Organization string
	Project      string

	APIMethod  string
	StatusCode int32
	Timestamp  int64
	LatencyMS  int32
}

// CreateUsage creates a new usage.
func (s *Store) CreateUsage(usage ...*Usage) error {
	return s.db.Create(usage).Error
}
