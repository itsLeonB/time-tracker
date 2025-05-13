package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/rotisserie/eris"
)

func BindRequest[T any](ctx *gin.Context) (T, error) {
	var request T

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		return request, eris.Wrap(err, apperror.MsgBindJsonError)
	}

	return request, nil
}

func GetUuidFromCtx(ctx *gin.Context, key string) (uuid.UUID, error) {
	ctxUserId, ok := ctx.Get(key)
	if !ok {
		return uuid.Nil, eris.Errorf("Missing user id in context")
	}

	userId, ok := ctxUserId.(string)
	if !ok {
		return uuid.Nil, eris.Errorf("Invalid user id type")
	}

	userUuid, err := ParseUuid(userId)
	if err != nil {
		return uuid.Nil, err
	}

	return userUuid, nil
}

func GetUuidParam(ctx *gin.Context, key string) (uuid.UUID, error) {
	val := ctx.Param(key)
	if val == "" || val == fmt.Sprintf(":%s", key) {
		return uuid.Nil, apperror.BadRequestError(eris.Errorf("missing path parameter"))
	}

	return ParseUuid(val)
}
