package service

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type ProjectService interface {
	Create(name string) (*model.Project, error)
	GetAll() ([]model.Project, error)
	GetByID(id uuid.UUID) (*model.Project, error)
	Update(id uuid.UUID, name string) (*model.Project, error)
	Delete(id uuid.UUID) error
}
