package repository

import (
	"context"
	"time"

	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/entity"
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

func (pr *projectRepositoryGorm) Insert(ctx context.Context, project entity.Project) (entity.Project, error) {
	db, err := pr.getGormInstance(ctx)
	if err != nil {
		return entity.Project{}, err
	}

	err = db.Create(&project).Error
	if err != nil {
		return entity.Project{}, eris.Wrap(err, appconstant.MsgInsertError)
	}

	return project, nil
}

func (pr *projectRepositoryGorm) FindAll(ctx context.Context, spec entity.ProjectSpecification) ([]entity.Project, error) {
	db, err := pr.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	var projects []entity.Project

	err = db.
		Scopes(util.DefaultOrdering(), ezutil.WhereBySpec(spec.Project)).
		Find(&projects).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return projects, nil
}

func (pr *projectRepositoryGorm) FindFirst(ctx context.Context, spec entity.ProjectSpecification) (entity.Project, error) {
	db, err := pr.getGormInstance(ctx)
	if err != nil {
		return entity.Project{}, err
	}

	var project entity.Project

	err = db.
		Scopes(
			ezutil.WhereBySpec(spec.Project),
			ezutil.PreloadRelations(spec.PreloadRelations),
		).
		First(&project).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Project{}, nil
		}

		return entity.Project{}, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return project, nil
}

func loadTasksWithLogs(start, end time.Time) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).Preload("Tasks.Logs", func(db *gorm.DB) *gorm.DB {
			query := db.Order("created_at DESC")

			whereClause, args := ezutil.GetTimeRangeClause("created_at", start, end)
			if whereClause != "" {
				query = query.Where(whereClause, args...)
			}

			return query
		})
	}
}

func (pr *projectRepositoryGorm) FindFirstPopulated(ctx context.Context, spec entity.ProjectSpecification) (entity.Project, error) {
	db, err := pr.getGormInstance(ctx)
	if err != nil {
		return entity.Project{}, err
	}

	var project entity.Project

	err = db.
		Scopes(
			ezutil.WhereBySpec(spec.Project),
			loadTasksWithLogs(spec.Start, spec.End),
		).
		First(&project).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Project{}, nil
		}

		return entity.Project{}, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return project, nil
}

func (pr *projectRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return pr.db.WithContext(ctx), nil
}
