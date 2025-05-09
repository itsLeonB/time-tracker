package model

import (
	"time"

	"github.com/google/uuid"
)

type QueryOptions struct {
	Params   *QueryParams
	WithLogs bool
}

type QueryParams struct {
	Number    string
	ProjectID uuid.UUID
	Date      time.Time
}

type ExternalQueryOptions struct {
	Number string
}

type YoutrackQueryOptions struct {
	IdReadable string
}
