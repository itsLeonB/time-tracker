package service

import (
	"context"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/mapper"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
)

type userTaskLogServiceImpl struct {
	userTaskLogRepository repository.UserTaskLogRepository
}

func NewUserTaskLogService(
	userTaskLogRepository repository.UserTaskLogRepository,
) UserTaskLogService {
	return &userTaskLogServiceImpl{
		userTaskLogRepository,
	}
}

func (utls *userTaskLogServiceImpl) Create(ctx context.Context, request dto.NewUserTaskLogRequest) (dto.UserTaskLogResponse, error) {
	var response dto.UserTaskLogResponse

	queryOptions := model.UserTaskLogQueryOptions{}

	queryOptions.Filters = map[string]any{
		"user_task_id": request.UserTaskId,
	}

	latestLog, err := utls.userTaskLogRepository.FindLatest(ctx, queryOptions)
	if err != nil {
		return response, err
	}

	// should always start with START action
	if latestLog.IsZero() && request.Action == constant.LogAction.Stop {
		return response, apperror.BadRequestError(eris.Errorf("User Task ID: %s is not yet started", request.UserTaskId))
	}

	// should alternate actions
	if !latestLog.IsZero() && latestLog.Action == request.Action {
		return response, apperror.BadRequestError(eris.Errorf("user task ID: %s is already %s", request.UserTaskId, request.Action))
	}

	newLog := model.UserTaskLog{
		UserTaskId: request.UserTaskId,
		Action:     request.Action,
	}

	insertedLog, err := utls.userTaskLogRepository.Insert(ctx, newLog)
	if err != nil {
		return response, err
	}

	return mapper.UserTaskLogToResponse(insertedLog), nil
}

func (utls *userTaskLogServiceImpl) FindAll(ctx context.Context, params dto.UserTaskLogParams) ([]dto.UserTaskLogResponse, error) {
	queryOptions := model.UserTaskLogQueryOptions{}

	queryOptions.Filterables = []model.Filter{
		{
			Column: "user_task_logs.user_task_id",
			Type:   model.FilterTypeIn,
			Value:  params.UserTaskIds,
		},
	}

	queryOptions.DatetimeRanges = map[string]model.DatetimeRange{
		"user_task_logs.created_at": {
			Start: params.StartDatetime,
			End:   params.EndDatetime,
		},
	}

	logs, err := utls.userTaskLogRepository.FindAll(ctx, queryOptions)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(logs, mapper.UserTaskLogToResponse), nil
}
