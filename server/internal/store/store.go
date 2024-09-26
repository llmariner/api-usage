package store

import "gorm.io/gorm"

// New creates a new store instance.
func New(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

// Store represents the data store.
type Store struct {
	db *gorm.DB
}

// AutoMigrate sets up the auto-migration task of the database.
func (s *Store) AutoMigrate() error {
	return autoMigrate(s.db)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Usage{})
}
