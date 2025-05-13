package util

import (
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func ParseUuid(val string) (uuid.UUID, error) {
	parsedUuid, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, eris.Wrap(err, "Failed to parse UUID")
	}

	return parsedUuid, nil
}
