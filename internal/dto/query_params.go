package dto

import (
	"time"

	"github.com/google/uuid"
)

type QueryParams struct {
	StartDatetime time.Time `query:"start"`
	EndDatetime   time.Time `query:"end"`
	UserId        uuid.UUID
}
