package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/constant"
)

type UserTask struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserProjectId uuid.UUID
	TaskId        uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Logs          []UserTaskLog `gorm:"foreignKey:UserTaskId"`
	Task          Task          `gorm:"foreignKey:TaskId;references:ID"`
}

func (userTask *UserTask) IsInProgress() bool {
	if len(userTask.Logs) == 0 {
		return false
	}

	return userTask.Logs[len(userTask.Logs)-1].Action == constant.LogAction.Start
}

func (ut *UserTask) TableName() string {
	return "user_tasks"
}

type UserTaskQueryOptions struct {
	Filters          map[string]any
	PreloadRelations []string
}
