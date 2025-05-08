package util

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/rotisserie/eris"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	id := ctx.Value(constant.ContextUserID)
	if id == nil {
		return uuid.Nil, eris.Errorf("id is not set in context")
	}

	userID, err := uuid.Parse(id.(string))
	if err != nil {
		return uuid.Nil, eris.Wrap(err, "error parsing UUID")
	}

	return userID, nil
}
