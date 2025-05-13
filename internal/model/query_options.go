package model

import (
	"time"

	"github.com/google/uuid"
)

type FilterType string

const (
	FilterTypeEqual FilterType = "EQUAL"
	FilterTypeIn    FilterType = "IN"
)

type DatetimeRange struct {
	Start time.Time
	End   time.Time
}

func (dr *DatetimeRange) IsZero() bool {
	return dr.Start.IsZero() && dr.End.IsZero()
}

type QueryOptions struct {
	Id               uuid.UUID
	Filters          map[string]any
	PreloadRelations []string
	DatetimeRanges   map[string]DatetimeRange
	Joins            []Join
	Filterables      []Filter
}

type Filter struct {
	Column string
	Value  any
	Type   FilterType
}

type Join struct {
	Table  string
	On     string
	Params []any
}
