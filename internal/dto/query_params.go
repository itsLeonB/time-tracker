package dto

import (
	"time"

	"github.com/google/uuid"
)

type QueryParams struct {
	StartDatetime time.Time `form:"start" time_format:"2006-01-02"`
	EndDatetime   time.Time `form:"end" time_format:"2006-01-02"`
	UserId        uuid.UUID
}
