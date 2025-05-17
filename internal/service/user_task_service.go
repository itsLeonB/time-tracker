package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/mapper"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
)

type userTaskServiceImpl struct {
	userTaskRepository repository.UserTaskRepository
}

func NewUserTaskService(
	userTaskRepository repository.UserTaskRepository,
) UserTaskService {
	return &userTaskServiceImpl{
		userTaskRepository,
	}
}

func (uts *userTaskServiceImpl) Create(ctx context.Context, request dto.NewUserTaskRequest) (dto.UserTaskResponse, error) {
	var userTaskResponse dto.UserTaskResponse

	newUserTask := model.UserTask{
		UserProjectId: request.UserProjectId,
		TaskId:        request.TaskId,
	}

	insertedUserTask, err := uts.userTaskRepository.Insert(ctx, newUserTask)
	if err != nil {
		return userTaskResponse, err
	}

	return mapper.UserTaskToResponse(insertedUserTask), nil
}

func (uts *userTaskServiceImpl) FindAll(ctx context.Context, params dto.UserTaskQueryParams) ([]dto.UserTaskResponse, error) {
	options := mapper.QueryParamsToOptions(params)

	userTasks, err := uts.userTaskRepository.FindAll(ctx, options)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(userTasks, mapper.UserTaskToResponse), nil
}

func (uts *userTaskServiceImpl) GetById(ctx context.Context, id uuid.UUID) (dto.UserTaskResponse, error) {
	var userTaskResponse dto.UserTaskResponse

	userTask, err := uts.userTaskRepository.FindById(ctx, id)
	if err != nil {
		return userTaskResponse, err
	}
	if userTask.ID == uuid.Nil {
		return userTaskResponse, eris.Errorf("User task with id: %s not found", id)
	}

	return mapper.UserTaskToResponse(userTask), nil
}
