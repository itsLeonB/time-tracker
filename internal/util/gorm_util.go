package util

import (
	"fmt"

	"gorm.io/gorm"
)

func PreloadRelations(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}

		return db
	}
}

func FilterByColumns(filters map[string]any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for col, val := range filters {
			if val != nil {
				db = db.Where(fmt.Sprintf("%s = ?", col), val)
			}
		}

		return db
	}
}
