package util

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/rotisserie/eris"
)

func BindJsonRequest[T any](ctx *gin.Context) (T, error) {
	var request T

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		return request, eris.Wrap(err, apperror.MsgBindJsonError)
	}

	return request, nil
}

func BindFormRequest[T any](ctx *gin.Context) (T, error) {
	var request T

	err := ctx.ShouldBind(&request)
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

func GetDatetimeParams(ctx *gin.Context, startKey string, endKey string) (dto.QueryParams, error) {
	var params dto.QueryParams

	startDate, err := parseDatetimeFromQuery(ctx, startKey, time.DateOnly)
	if err != nil {
		return params, err
	}

	endDate, err := parseDatetimeFromQuery(ctx, endKey, time.DateOnly)
	if err != nil {
		return params, err
	}

	if !endDate.IsZero() {
		endDate = endDate.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
	}

	params.StartDatetime = startDate
	params.EndDatetime = endDate

	return params, nil
}

func parseDatetimeFromQuery(ctx *gin.Context, key string, layout string) (time.Time, error) {
	val := ctx.Query(key)
	if val == "" {
		return time.Time{}, nil
	}

	result, err := time.Parse(layout, val)
	if err != nil {
		return time.Time{}, eris.Wrap(err, fmt.Sprintf("Invalid datetime format for %s", key))
	}

	return result, nil
}
