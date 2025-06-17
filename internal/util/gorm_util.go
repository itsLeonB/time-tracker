package util

import (
	"gorm.io/gorm"
)

func DefaultOrdering() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}
}
