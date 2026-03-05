package repository

import "gorm.io/gorm"

// ReadRepository combines all read operations needed by the read-api.
type ReadRepository struct {
	db *gorm.DB
}

func NewReadRepository(db *gorm.DB) *ReadRepository {
	return &ReadRepository{db: db}
}
