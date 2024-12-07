package service

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
)

type projectServiceImpl struct {
	projectRepository repository.ProjectRepository
}

func NewProjectService(projectRepository repository.ProjectRepository) ProjectService {
	return &projectServiceImpl{projectRepository}
}

func (ps *projectServiceImpl) Create(name string) (*model.Project, error) {
	return ps.projectRepository.Insert(&model.Project{Name: name})
}

func (ps *projectServiceImpl) GetAll() ([]*model.Project, error) {
	return ps.projectRepository.GetAll()
}

func (ps *projectServiceImpl) GetByID(id uuid.UUID) (*model.Project, error) {
	return ps.projectRepository.GetByID(id)
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
