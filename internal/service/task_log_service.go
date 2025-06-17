package service

import (
	"context"

	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
	"github.com/itsLeonB/time-tracker/internal/mapper"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type taskLogServiceImpl struct {
	taskLogRepository repository.TaskLogRepository
}

func NewTaskLogService(
	taskLogRepository repository.TaskLogRepository,
) TaskLogService {
	return &taskLogServiceImpl{
		taskLogRepository,
	}
}

func (tls *taskLogServiceImpl) Create(ctx context.Context, request dto.NewTaskLogRequest) (dto.TaskLogResponse, error) {
	spec := entity.TaskLog{TaskID: request.TaskID}

	latestLog, err := tls.taskLogRepository.FindLatest(ctx, spec)
	if err != nil {
		return dto.TaskLogResponse{}, err
	}

	// should always start with START action
	if latestLog.IsZero() && request.Action == appconstant.LogAction.Stop {
		return dto.TaskLogResponse{}, ezutil.BadRequestError(eris.Errorf("User Task ID: %s is not yet started", request.TaskID))
	}

	// should alternate actions
	if !latestLog.IsZero() && latestLog.Action == request.Action {
		return dto.TaskLogResponse{}, ezutil.BadRequestError(eris.Errorf("user task ID: %s is already %s", request.TaskID, request.Action))
	}

	newLog := entity.TaskLog{
		TaskID: request.TaskID,
		Action: request.Action,
	}

	insertedLog, err := tls.taskLogRepository.Insert(ctx, newLog)
	if err != nil {
		return dto.TaskLogResponse{}, err
	}

	return mapper.TaskLogToResponse(insertedLog), nil
}
