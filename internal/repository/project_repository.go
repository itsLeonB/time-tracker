package repository

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
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
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return project, nil
}

func (pr *projectRepositoryGorm) GetAll() ([]*model.Project, error) {
	var projects []*model.Project

	err := pr.db.Find(&projects).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) GetByID(id uuid.UUID) (*model.Project, error) {
	var project model.Project

	err := pr.db.First(&project, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &project, nil
}

func (pr *projectRepositoryGorm) Update(project *model.Project) (*model.Project, error) {
	err := pr.db.Save(project).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgUpdateError))
	}

	return project, nil
}

func (pr *projectRepositoryGorm) Delete(project *model.Project) error {
	err := pr.db.Delete(project).Error
	if err != nil {
		return apperror.InternalServerError(eris.Wrap(err, apperror.MsgDeleteError))
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
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return projects, nil
}
