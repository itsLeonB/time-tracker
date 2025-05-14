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

type taskServiceImpl struct {
	taskRepository         repository.TaskRepository
	userService            UserService
	externalTrackerService ExternalTrackerService
	projectService         ProjectService
	userTaskService        UserTaskService
}

func NewTaskService(
	taskRepository repository.TaskRepository,
	userService UserService,
	externalTrackerService ExternalTrackerService,
	projectService ProjectService,
	userTaskService UserTaskService,
) TaskService {
	return &taskServiceImpl{
		taskRepository,
		userService,
		externalTrackerService,
		projectService,
		userTaskService,
	}
}

func (ts *taskServiceImpl) Create(ctx context.Context, request *dto.NewTaskRequest) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.taskRepository.GetByNumber(ctx, request.Number)
	if err != nil {
		return taskResponse, err
	}
	if task != nil {
		return taskResponse, apperror.ConflictError(eris.Errorf("task number: %s already exists", request.Number))
	}

	newTask := &model.Task{
		ProjectID: request.ProjectID,
		Number:    request.Number,
		Name:      request.Name,
	}

	insertedTask, err := ts.taskRepository.Insert(ctx, newTask)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*insertedTask), nil
}

func (ts *taskServiceImpl) GetAll(ctx context.Context) ([]dto.TaskResponse, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(tasks, mapper.TaskToResponse), nil
}

func (ts *taskServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.taskRepository.GetByID(ctx, id)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*task), nil
}

func (ts *taskServiceImpl) Find(ctx context.Context, queryParams dto.TaskQueryParams) ([]dto.TaskResponse, error) {
	tasks, _, err := ts.findTask(ctx, queryParams)
	if err != nil {
		return nil, err
	}

	if len(tasks) > 0 {
		return util.MapSlice(tasks, mapper.TaskToResponse), nil
	}

	// If no tasks are found, check the external tracker
	return ts.findAndInsertFromExternal(ctx, queryParams.Number, uuid.Nil)
}

func (ts *taskServiceImpl) findAndInsertFromExternal(ctx context.Context, number string, projectId uuid.UUID) ([]dto.TaskResponse, error) {
	externalQueryOptions := dto.ExternalQueryOptions{
		Number: number,
	}

	externalTasks, err := ts.externalTrackerService.FindTask(ctx, externalQueryOptions)
	if err != nil {
		return nil, err
	}

	insertedTasks := make([]dto.TaskResponse, 0, len(externalTasks))

	for _, externalTask := range externalTasks {
		insertedTask, err := ts.insertFromExternal(ctx, externalTask, projectId)
		if err != nil {
			return nil, err
		}

		insertedTasks = append(insertedTasks, mapper.TaskToResponse(insertedTask))
	}

	return insertedTasks, nil
}

func (ts *taskServiceImpl) insertFromExternal(ctx context.Context, externalTask model.ExternalTask, projectId uuid.UUID) (model.Task, error) {
	var task model.Task

	if projectId == uuid.Nil {
		project, err := ts.projectService.GetOrCreate(ctx, externalTask.Project)
		if err != nil {
			return task, err
		}

		projectId = project.ID
	}

	newTask := model.Task{
		ProjectID: projectId,
		Number:    externalTask.Number,
		Name:      externalTask.Name,
	}

	insertedTask, err := ts.taskRepository.Insert(ctx, &newTask)
	if err != nil {
		return task, err
	}

	return *insertedTask, nil
}

func (ts *taskServiceImpl) GetByNumber(ctx context.Context, number string) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.getTaskByNumber(ctx, number)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*task), nil
}

func (ts *taskServiceImpl) Update(ctx context.Context, id uuid.UUID, name string) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.getTask(ctx, id)
	if err != nil {
		return taskResponse, err
	}

	task.Name = name

	updatedTask, err := ts.taskRepository.Update(ctx, task)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*updatedTask), nil
}

func (ts *taskServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return err
	}

	task, err := ts.getTask(ctx, id)
	if err != nil {
		return err
	}

	return ts.taskRepository.Delete(ctx, task)
}

func (ts *taskServiceImpl) getTask(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	task, err := ts.taskRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperror.NotFoundError(eris.Errorf("task ID: %s is not found", id))
	}

	return task, nil
}

func (ts *taskServiceImpl) getTaskByNumber(ctx context.Context, number string) (*model.Task, error) {
	task, err := ts.taskRepository.GetByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperror.NotFoundError(eris.Errorf("task number: %s is not found", number))
	}

	return task, nil
}

func (ts *taskServiceImpl) AddToUserProject(ctx context.Context, request dto.AddToProjectRequest) (dto.TaskResponse, error) {
	params := dto.TaskQueryParams{
		Number:    request.Number,
		ProjectID: request.ProjectID,
	}

	tasks, user, err := ts.findTask(ctx, params)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	if len(tasks) > 0 {
		for _, task := range tasks {
			userTaskParams := dto.UserTaskQueryParams{
				UserId: user.ID,
				TaskId: task.ID,
			}
			userTasks, err := ts.userTaskService.FindAll(ctx, userTaskParams)
			if err != nil {
				return dto.TaskResponse{}, err
			}

			if len(userTasks) > 0 {
				return userTasks[0].Task, nil
			}
		}
	}

	taskResponse, err := ts.findAndInsertFromExternal(ctx, request.Number, request.ProjectID)
	if err != nil {
		return dto.TaskResponse{}, err
	}
	if len(taskResponse) == 0 {
		return dto.TaskResponse{}, eris.Errorf("Task not found: %s", request.Number)
	}

	newUserTask := dto.NewUserTaskRequest{
		TaskId: taskResponse[0].ID,
		UserId: user.ID,
	}

	_, err = ts.userTaskService.Create(ctx, newUserTask)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	return taskResponse[0], nil
}

func (ts *taskServiceImpl) findTask(ctx context.Context, params dto.TaskQueryParams) ([]model.Task, model.User, error) {
	user, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, model.User{}, err
	}

	queryOptions := model.TaskQueryOptions{
		Number:    params.Number,
		ProjectID: params.ProjectID,
		Date:      params.Date,
	}

	tasks, err := ts.taskRepository.Find(ctx, queryOptions)
	if err != nil {
		return nil, model.User{}, err
	}

	return tasks, *user, nil
}
