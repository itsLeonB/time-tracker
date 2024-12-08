package service

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type projectServiceImpl struct {
	projectRepository repository.ProjectRepository
	taskService       TaskService
}

func NewProjectService(
	projectRepository repository.ProjectRepository,
	taskService TaskService,
) ProjectService {
	return &projectServiceImpl{projectRepository, taskService}
}

func (ps *projectServiceImpl) Create(name string) (*model.Project, error) {
	return ps.projectRepository.Insert(&model.Project{Name: name})
}

func (ps *projectServiceImpl) GetAll() ([]*model.Project, error) {
	return ps.projectRepository.GetAll()
}

func (ps *projectServiceImpl) GetByID(options *model.QueryOptions) (*model.Project, error) {
	project, err := ps.projectRepository.GetByID(options.Params.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.NotFoundError(
			eris.Errorf("project with ID: %s not found", options.Params.ProjectID),
		)
	}

	tasks, err := ps.taskService.Find(options)
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

func (ps *projectServiceImpl) Update(id uuid.UUID, name string) (*model.Project, error) {
	return ps.projectRepository.Update(&model.Project{
		ID:   id,
		Name: name,
	})
}

func (ps *projectServiceImpl) Delete(id uuid.UUID) error {
	return ps.projectRepository.Delete(&model.Project{ID: id})
}

func (ps *projectServiceImpl) FirstByQuery(options model.FindProjectOptions) (*model.Project, error) {
	projects, err := ps.projectRepository.Find(&options)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, nil
	}

	return projects[0], nil
}
