package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
	"gorm.io/gorm"
)

func PreloadRelations(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation, DefaultOrdering())
		}

		return db
	}
}

func FilterByColumns(filters map[string]any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for col, val := range filters {
			if val != nil {
				uuidVal, ok := val.(uuid.UUID)
				if !ok || (ok && uuidVal != uuid.Nil) {
					db = db.Where(fmt.Sprintf("%s = ?", col), val)
				}
			}
		}

		return db
	}
}

func WithFilters(filters []model.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, filter := range filters {
			clauseBuilder := strings.Builder{}
			clauseBuilder.WriteString(filter.Column)

			switch filter.Type {
			case model.FilterTypeEqual:
				clauseBuilder.WriteString(" = ?")
			case model.FilterTypeIn:
				clauseBuilder.WriteString(" IN ?")
			}

			db = db.Where(clauseBuilder.String(), filter.Value)
		}

		return db
	}
}

func GetTimeRangeFilter(timeCol string, start time.Time, end time.Time) (string, []any) {
	if start.IsZero() && end.IsZero() {
		return "", nil
	}

	if start.IsZero() {
		return fmt.Sprintf("%s <= ?", timeCol), []any{end}
	}

	if end.IsZero() {
		return fmt.Sprintf("%s >= ?", timeCol), []any{start}
	}

	return fmt.Sprintf("%s BETWEEN ? AND ?", timeCol), []any{start, end}
}

func WithinTimeRange(timeCol string, start time.Time, end time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start.IsZero() && end.IsZero() {
			return db
		}

		if start.IsZero() {
			return db.Where(fmt.Sprintf("%s <= ?", timeCol), end)
		}

		if end.IsZero() {
			return db.Where(fmt.Sprintf("%s >= ?", timeCol), start)
		}

		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", timeCol), start, end)
	}
}

func WithinTimeRanges(ranges map[string]model.DatetimeRange) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for col, rangeVal := range ranges {
			db = WithinTimeRange(col, rangeVal.Start, rangeVal.End)(db)
		}

		return db
	}
}

func Join(joins []model.Join) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, join := range joins {
			db = db.Joins(fmt.Sprintf("LEFT JOIN %s ON %s", join.Table, join.On), join.Params...)
		}

		return db
	}
}

func DefaultOrdering() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}
}
