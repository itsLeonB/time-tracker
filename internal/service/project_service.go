package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
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

func (ps *projectServiceImpl) Create(ctx context.Context, name string) (*model.Project, error) {
	user, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	existingProject, err := ps.projectRepository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existingProject != nil {
		return nil, apperror.ConflictError(eris.Errorf("project with name: %s already exists", name))
	}

	newProject := &model.Project{
		UserID: user.ID,
		Name:   name,
	}

	return ps.projectRepository.Insert(ctx, newProject)
}

func (ps *projectServiceImpl) GetAll(ctx context.Context) ([]*model.Project, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	return ps.projectRepository.GetAll(ctx)
}

func (ps *projectServiceImpl) GetByID(ctx context.Context, options *model.QueryOptions) (*model.Project, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	project, err := ps.getProject(ctx, options.Params.ProjectID)
	if err != nil {
		return nil, err
	}

	tasks, err := ps.taskService.Find(ctx, options)
	if err != nil {
		return nil, err
	}

	timeSpent := new(model.TimeSpent)
	for _, task := range tasks {
		project.TotalPoints += task.Points
		timeSpent.Add(task.TimeSpent)
	}

	project.TimeSpent = timeSpent

	return project, nil
}

func (ps *projectServiceImpl) Update(ctx context.Context, id uuid.UUID, name string) (*model.Project, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	project, err := ps.getProject(ctx, id)
	if err != nil {
		return nil, err
	}

	project.Name = name

	return ps.projectRepository.Update(ctx, project)
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

func (ps *projectServiceImpl) FirstByQuery(ctx context.Context, options *model.FindProjectOptions) (*model.Project, error) {
	_, err := ps.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	projects, err := ps.projectRepository.Find(ctx, options)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, nil
	}

	return projects[0], nil
}

func (ps *projectServiceImpl) GetInProgressTasks(ctx context.Context, id uuid.UUID) ([]*model.Task, error) {
	_, err := ps.getProject(ctx, id)
	if err != nil {
		return nil, err
	}

	tasks, err := ps.taskRepository.GetInProgress(ctx, id)
	if err != nil {
		return nil, err
	}

	return tasks, nil
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
