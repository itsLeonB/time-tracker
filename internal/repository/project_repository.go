package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/util"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type projectRepositoryGorm struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepositoryGorm{db}
}

func (pr *projectRepositoryGorm) Insert(ctx context.Context, project *model.UserProject) (*model.UserProject, error) {
	err := pr.db.WithContext(ctx).Create(project).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return project, nil
}

func (pr *projectRepositoryGorm) GetAll(ctx context.Context) ([]model.UserProject, error) {
	var projects []model.UserProject

	err := pr.db.WithContext(ctx).
		Scopes(util.DefaultOrdering()).
		Find(&projects).
		Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*model.UserProject, error) {
	var project model.UserProject

	err := pr.db.WithContext(ctx).
		Scopes(util.DefaultOrdering()).
		First(&project, "id = ?", id).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &project, nil
}

func (pr *projectRepositoryGorm) Update(ctx context.Context, project *model.UserProject) (*model.UserProject, error) {
	err := pr.db.WithContext(ctx).Save(project).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgUpdateError))
	}

	return project, nil
}

func (pr *projectRepositoryGorm) Delete(ctx context.Context, project *model.UserProject) error {
	err := pr.db.WithContext(ctx).Delete(project).Error
	if err != nil {
		return apperror.InternalServerError(eris.Wrap(err, apperror.MsgDeleteError))
	}

	return nil
}

func (pr *projectRepositoryGorm) Find(ctx context.Context, options dto.UserProjectQueryParams) ([]model.UserProject, error) {
	var projects []model.UserProject

	query := pr.db.WithContext(ctx).Scopes(util.DefaultOrdering())

	if options.Name != "" {
		query = query.Where("name = ?", options.Name)
	}

	if options.Ids != nil {
		query = query.Where("id IN ?", options.Ids)
	}

	err := query.Find(&projects).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) GetByName(ctx context.Context, name string) (*model.UserProject, error) {
	var project model.UserProject

	err := pr.db.WithContext(ctx).
		Scopes(util.DefaultOrdering()).
		First(&project, "name = ?", name).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &project, nil
}

func (pr *projectRepositoryGorm) First(ctx context.Context, options model.ProjectQueryOptions) (model.UserProject, error) {
	var project model.UserProject

	query := pr.db.WithContext(ctx).Scopes(
		util.FilterByColumns(options.Filters),
		util.PreloadRelations(options.PreloadRelations),
		util.Join(options.Joins),
		util.DefaultOrdering(),
	)

	err := query.First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return project, nil
		}

		return project, eris.Wrap(err, apperror.MsgQueryError)
	}

	return project, nil
}
