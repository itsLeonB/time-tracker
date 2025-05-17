package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
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

func (userTask *UserTask) GetStartTime() string {
	if len(userTask.Logs) > 0 {
		if userTask.Logs[0].Action == constant.LogAction.Start {
			return userTask.Logs[0].CreatedAt.Format(time.RFC3339)
		}
	}

	return ""
}

func (userTask *UserTask) IsActive() bool {
	return len(userTask.Logs) > 0 && userTask.Logs[0].Action == "START"
}

func (userTask *UserTask) IsInProgress() bool {
	if len(userTask.Logs) == 0 {
		return false
	}

	return userTask.Logs[len(userTask.Logs)-1].Action == constant.LogAction.Start
}

func (userTask *UserTask) GetTotalTime() time.Duration {
	var totalTime time.Duration

	for i := 0; i < len(userTask.Logs); i += 2 {
		if i+1 < len(userTask.Logs) {
			startLog := userTask.Logs[i]
			stopLog := userTask.Logs[i+1]

			if startLog.Action == constant.LogAction.Start && stopLog.Action == constant.LogAction.Stop {
				duration := stopLog.CreatedAt.Sub(startLog.CreatedAt)
				totalTime += duration
			}
		}
	}

	return totalTime
}

func (ut *UserTask) TableName() string {
	return "user_tasks"
}

type UserTaskQueryOptions struct {
	Filters          map[string]any
	PreloadRelations []string
}
