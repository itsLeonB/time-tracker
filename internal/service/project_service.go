package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
	"github.com/itsLeonB/time-tracker/internal/mapper"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type projectServiceImpl struct {
	projectRepository repository.ProjectRepository
	userService       UserService
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	userService UserService,
) ProjectService {
	return &projectServiceImpl{
		projectRepository,
		userService,
	}
}

func (ps *projectServiceImpl) Create(ctx context.Context, name string) (dto.ProjectResponse, error) {
	if strings.TrimSpace(name) == "" {
		return dto.ProjectResponse{}, ezutil.BadRequestError(eris.New("project name must not be empty"))
	}

	user, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	spec := entity.ProjectSpecification{}
	spec.Name = name
	spec.UserID = user.ID

	existingProject, err := ps.projectRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.ProjectResponse{}, err
	}
	if !existingProject.IsZero() {
		return dto.ProjectResponse{}, ezutil.ConflictError(eris.Errorf("project with name: %s already exists", name))
	}

	newProject := entity.Project{
		UserID: user.ID,
		Name:   name,
	}

	insertedProject, err := ps.projectRepository.Insert(ctx, newProject)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	return mapper.ProjectToResponse(insertedProject)
}

func (ps *projectServiceImpl) FirstByQuery(ctx context.Context, params dto.ProjectQueryParams) (dto.ProjectResponse, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	spec := entity.ProjectSpecification{}
	spec.ID = params.ID
	spec.Start = params.StartDatetime
	spec.End = params.EndDatetime

	project, err := ps.projectRepository.FindFirstPopulated(ctx, spec)
	if err != nil {
		return dto.ProjectResponse{}, err
	}
	if project.IsZero() {
		return dto.ProjectResponse{}, ezutil.NotFoundError(eris.Errorf("project not found"))
	}

	projectResponse, err := mapper.ProjectToResponse(project)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	return projectResponse, nil
}

func (ps *projectServiceImpl) FindByUserId(ctx context.Context, userId uuid.UUID) ([]dto.ProjectResponse, error) {
	spec := entity.ProjectSpecification{}
	spec.UserID = userId

	projects, err := ps.projectRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSliceWithError(projects, mapper.ProjectToResponse)
}
