package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type projectRepositoryGorm struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepositoryGorm{db}
}

func (pr *projectRepositoryGorm) Insert(ctx context.Context, project *model.Project) (*model.Project, error) {
	err := pr.db.WithContext(ctx).Create(project).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return project, nil
}

func (pr *projectRepositoryGorm) GetAll(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project

	err := pr.db.WithContext(ctx).Find(&projects).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var project model.Project

	err := pr.db.WithContext(ctx).First(&project, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &project, nil
}

func (pr *projectRepositoryGorm) Update(ctx context.Context, project *model.Project) (*model.Project, error) {
	err := pr.db.WithContext(ctx).Save(project).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgUpdateError))
	}

	return project, nil
}

func (pr *projectRepositoryGorm) Delete(ctx context.Context, project *model.Project) error {
	err := pr.db.WithContext(ctx).Delete(project).Error
	if err != nil {
		return apperror.InternalServerError(eris.Wrap(err, apperror.MsgDeleteError))
	}

	return nil
}

func (pr *projectRepositoryGorm) Find(ctx context.Context, options *dto.FindProjectOptions) ([]model.Project, error) {
	var projects []model.Project

	query := pr.db.WithContext(ctx)
	if options != nil {
		if options.Name != "" {
			query = query.Where("name = ?", options.Name)
		}

		if options.Ids != nil {
			query = query.Where("id IN ?", options.Ids)
		}
	}

	err := query.Find(&projects).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) GetByName(ctx context.Context, name string) (*model.Project, error) {
	var project model.Project

	err := pr.db.WithContext(ctx).First(&project, "name = ?", name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &project, nil
}
