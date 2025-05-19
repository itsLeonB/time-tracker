package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/mapper"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
)

type userTaskServiceImpl struct {
	userTaskRepository repository.UserTaskRepository
	taskService        TaskService
}

func NewUserTaskService(
	userTaskRepository repository.UserTaskRepository,
	taskService TaskService,
) UserTaskService {
	return &userTaskServiceImpl{
		userTaskRepository,
		taskService,
	}
}

func (uts *userTaskServiceImpl) Create(ctx context.Context, request dto.NewUserTaskRequest) (dto.UserTaskResponse, error) {
	var userTaskResponse dto.UserTaskResponse

	taskParams := dto.TaskQueryParams{
		Number: request.Number,
	}

	tasks, err := uts.taskService.Find(ctx, taskParams)
	if err != nil {
		return dto.UserTaskResponse{}, err
	}

	if len(tasks) == 0 {
		return dto.UserTaskResponse{}, apperror.NotFoundError(eris.Errorf("Task number %s is not found", request.Number))
	}

	newUserTask := model.UserTask{
		UserProjectId: request.UserProjectId,
		TaskId:        tasks[0].ID,
	}

	insertedUserTask, err := uts.userTaskRepository.Insert(ctx, newUserTask)
	if err != nil {
		return userTaskResponse, err
	}

	return mapper.UserTaskToResponse(insertedUserTask)
}

func (uts *userTaskServiceImpl) FindAll(ctx context.Context, params dto.UserTaskQueryParams) ([]dto.UserTaskResponse, error) {
	options := mapper.QueryParamsToOptions(params)

	userTasks, err := uts.userTaskRepository.FindAll(ctx, options)
	if err != nil {
		return nil, err
	}

	return util.MapSliceWithError(userTasks, mapper.UserTaskToResponse)
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

	return mapper.UserTaskToResponse(userTask)
}
