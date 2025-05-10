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
	projectRepository repository.ProjectRepository
	taskService       TaskService
	userService       UserService
	taskRepository    repository.TaskRepository
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	taskService TaskService,
	userService UserService,
	taskRepository repository.TaskRepository,
) ProjectService {
	return &projectServiceImpl{
		projectRepository,
		taskService,
		userService,
		taskRepository,
	}
}

func (ps *projectServiceImpl) Create(ctx context.Context, name string) (dto.ProjectResponse, error) {
	var projectResponse dto.ProjectResponse

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

	newProject := &model.Project{
		Name: name,
	}

	insertedProject, err := ps.projectRepository.Insert(ctx, newProject)
	if err != nil {
		return projectResponse, err
	}

	return mapper.ProjectToResponse(*insertedProject), nil
}

func (ps *projectServiceImpl) GetAll(ctx context.Context) ([]dto.ProjectResponse, error) {
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

func (ps *projectServiceImpl) GetByID(ctx context.Context, options *dto.QueryOptions) (dto.ProjectResponse, error) {
	var projectResponse dto.ProjectResponse

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

func (ps *projectServiceImpl) Update(ctx context.Context, id uuid.UUID, name string) (dto.ProjectResponse, error) {
	var projectResponse dto.ProjectResponse

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

func (ps *projectServiceImpl) FirstByQuery(ctx context.Context, options *dto.FindProjectOptions) (dto.ProjectResponse, error) {
	var projectResponse dto.ProjectResponse

	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return projectResponse, err
	}

	projects, err := ps.projectRepository.Find(ctx, options)
	if err != nil {
		return projectResponse, err
	}

	if len(projects) == 0 {
		return projectResponse, nil
	}

	return mapper.ProjectToResponse(*projects[0]), nil
}

func (ps *projectServiceImpl) getProject(ctx context.Context, id uuid.UUID) (*model.Project, error) {
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
