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

type projectServiceImpl struct {
	projectRepository  repository.ProjectRepository
	userService        UserService
	userTaskService    UserTaskService
	userTaskLogService UserTaskLogService
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	userService UserService,
	userTaskService UserTaskService,
	userTaskLogService UserTaskLogService,
) ProjectService {
	return &projectServiceImpl{
		projectRepository,
		userService,
		userTaskService,
		userTaskLogService,
	}
}

func (ps *projectServiceImpl) Create(ctx context.Context, name string) (dto.UserProjectResponse, error) {
	var projectResponse dto.UserProjectResponse

	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return projectResponse, err
	}

	existingProject, err := ps.projectRepository.GetByName(ctx, name)
	if err != nil {
		return projectResponse, err
	}
	if existingProject != nil {
		return projectResponse, apperror.ConflictError(eris.Errorf("project with name: %s already exists", name))
	}

	newProject := &model.UserProject{
		Name: name,
	}

	insertedProject, err := ps.projectRepository.Insert(ctx, newProject)
	if err != nil {
		return projectResponse, err
	}

	return mapper.ProjectToResponse(*insertedProject), nil
}

func (ps *projectServiceImpl) GetAll(ctx context.Context) ([]dto.UserProjectResponse, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	projects, err := ps.projectRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(projects, mapper.ProjectToResponse), nil
}

func (ps *projectServiceImpl) GetByID(ctx context.Context, options *dto.QueryOptions) (dto.UserProjectResponse, error) {
	var projectResponse dto.UserProjectResponse

	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return projectResponse, err
	}

	project, err := ps.getProject(ctx, options.Params.ProjectID)
	if err != nil {
		return projectResponse, err
	}

	return mapper.ProjectToResponse(*project), nil
}

func (ps *projectServiceImpl) Update(ctx context.Context, id uuid.UUID, name string) (dto.UserProjectResponse, error) {
	var projectResponse dto.UserProjectResponse

	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return projectResponse, err
	}

	project, err := ps.getProject(ctx, id)
	if err != nil {
		return projectResponse, err
	}

	project.Name = name

	updatedProject, err := ps.projectRepository.Update(ctx, project)
	if err != nil {
		return projectResponse, err
	}

	return mapper.ProjectToResponse(*updatedProject), nil
}

func (ps *projectServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return err
	}

	project, err := ps.getProject(ctx, id)
	if err != nil {
		return err
	}

	return ps.projectRepository.Delete(ctx, project)
}

func (ps *projectServiceImpl) FirstByQuery(ctx context.Context, params dto.UserProjectQueryParams) (dto.UserProjectResponse, error) {
	var projectResponse dto.UserProjectResponse

	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return projectResponse, err
	}

	options := model.ProjectQueryOptions{}

	options.Filters = map[string]any{
		"projects.id":        params.ProjectID,
		"user_tasks.user_id": params.UserId,
	}

	options.Joins = []model.Join{
		{
			Table: "tasks",
			On:    "tasks.project_id = projects.id",
		},
		{
			Table: "user_tasks",
			On:    "user_tasks.task_id = tasks.id",
		},
	}

	options.PreloadRelations = []string{
		"Tasks",
		"Tasks.UserTasks",
	}

	project, err := ps.projectRepository.First(ctx, options)
	if err != nil {
		return projectResponse, err
	}
	if project.IsZero() {
		return projectResponse, apperror.NotFoundError(eris.Errorf("project not found"))
	}

	userTaskIds := make([]uuid.UUID, 0)
	for _, userTask := range project.UserTasks {
		userTaskIds = append(userTaskIds, userTask.ID)
	}

	userTaskLogParams := dto.UserTaskLogParams{
		UserTaskIds: userTaskIds,
	}

	userTaskLogParams.StartDatetime = params.StartDatetime
	userTaskLogParams.EndDatetime = params.EndDatetime

	userTaskLogs, err := ps.userTaskLogService.FindAll(ctx, userTaskLogParams)
	if err != nil {
		return projectResponse, err
	}

	projectResponse, err = mapper.PopulateProjectWithLogs(project, userTaskLogs)
	if err != nil {
		return projectResponse, err
	}

	return projectResponse, nil
}

func (ps *projectServiceImpl) GetOrCreate(ctx context.Context, name string) (dto.UserProjectResponse, error) {
	var projectResponse dto.UserProjectResponse

	project, err := ps.projectRepository.GetByName(ctx, name)
	if err != nil {
		return projectResponse, err
	}
	if project != nil {
		// Found project, return
		return mapper.ProjectToResponse(*project), nil
	}

	// Not found, insert new
	newProject := model.UserProject{
		Name: name,
	}

	project, err = ps.projectRepository.Insert(ctx, &newProject)
	if err != nil {
		return projectResponse, err
	}

	return mapper.ProjectToResponse(*project), nil
}

func (ps *projectServiceImpl) FindByUserId(ctx context.Context, userId uuid.UUID) ([]dto.UserProjectResponse, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	projectOptions := &dto.UserProjectQueryParams{}
	projectOptions.UserId = userId

	projects, err := ps.projectRepository.Find(ctx, *projectOptions)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(projects, mapper.ProjectToResponse), nil
}

func (ps *projectServiceImpl) getProject(ctx context.Context, id uuid.UUID) (*model.UserProject, error) {
	project, err := ps.projectRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.NotFoundError(
			eris.Errorf("project with ID: %s not found", id),
		)
	}

	return project, nil
}
