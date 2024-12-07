package repository

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type projectRepositoryGorm struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepositoryGorm{db}
}

func (pr *projectRepositoryGorm) Insert(project *model.Project) (*model.Project, error) {
	err := pr.db.Create(project).Error
	if err != nil {
		return nil, eris.Wrap(err, "error inserting")
	}

	return project, nil
}

func (pr *projectRepositoryGorm) GetAll() ([]*model.Project, error) {
	var projects []*model.Project

	err := pr.db.Find(&projects).Error
	if err != nil {
		return nil, eris.Wrap(err, "error querying")
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) GetByID(id uuid.UUID) (*model.Project, error) {
	var project model.Project

	err := pr.db.First(&project, "id = ?", id).Error
	if err != nil {
		return nil, eris.Wrap(err, "error querying")
	}

	return &project, nil
}

func (pr *projectRepositoryGorm) Update(project *model.Project) (*model.Project, error) {
	err := pr.db.Save(project).Error
	if err != nil {
		return nil, eris.Wrap(err, "error updating")
	}

	return project, nil
}

func (pr *projectRepositoryGorm) Delete(project *model.Project) error {
	err := pr.db.Delete(project).Error
	if err != nil {
		return eris.Wrap(err, "error deleting")
	}

	return nil
}

func (pr *projectRepositoryGorm) Find(options *model.FindProjectOptions) ([]*model.Project, error) {
	var projects []*model.Project

	query := pr.db
	if options != nil {
		if options.Name != "" {
			query = query.Where("name ILIKE %?%", options.Name)
		}
	}

	err := query.Find(&projects).Error
	if err != nil {
		return nil, eris.Wrap(err, "error querying")
	}

	return projects, nil
}
