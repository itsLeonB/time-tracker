package dto

import (
	"time"

	"github.com/google/uuid"
)

type QueryParams struct {
	StartDatetime time.Time `form:"startDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDatetime   time.Time `form:"endDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	UserId        uuid.UUID
}
