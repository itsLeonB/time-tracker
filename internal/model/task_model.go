package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/constant"
)

type Task struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID  uuid.UUID  `json:"projectId"`
	Number     string     `json:"number"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	InProgress bool       `json:"inProgress" gorm:"-"`
	Points     float64    `json:"points" gorm:"-"`
	TimeSpent  *TimeSpent `json:"timeSpent" gorm:"-"`
	Logs       []TaskLog  `json:"logs,omitempty" gorm:"foreignKey:TaskID"`
}

func (task *Task) DetermineProgress() {
	if len(task.Logs) == 0 {
		task.InProgress = false
		return
	}

	task.InProgress = task.Logs[len(task.Logs)-1].Action == constant.LogAction.Start
}

func (task *Task) CalculateTotalTime() {
	var totalTime time.Duration

	for i := 0; i < len(task.Logs); i += 2 {
		if i+1 < len(task.Logs) {
			startLog := task.Logs[i]
			stopLog := task.Logs[i+1]

			if startLog.Action == constant.LogAction.Start && stopLog.Action == constant.LogAction.Stop {
				duration := stopLog.CreatedAt.Sub(startLog.CreatedAt)
				totalTime += duration
			}
		}
	}

	task.TimeSpent = &TimeSpent{
		Duration: totalTime,
		Minutes:  totalTime.Minutes(),
		Hours:    totalTime.Hours(),
		String:   totalTime.String(),
	}
}

func (t *Task) TableName() string {
	return "tasks"
}

type NewTaskRequest struct {
	ProjectID uuid.UUID `json:"projectId" binding:"required"`
	Number    string    `json:"number" binding:"required"`
	Name      string    `json:"name" binding:"required"`
}

type TaskLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID    uuid.UUID `json:"taskId"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (tl *TaskLog) TableName() string {
	return "task_logs"
}

type NewLogRequest struct {
	Action string `json:"action" binding:"required,oneof=START STOP"`
}

type LogByNumberRequest struct {
	Number string `json:"number" binding:"required"`
	Action string `json:"action" binding:"required,oneof=START STOP"`
}

type TimeSpent struct {
	Duration time.Duration `json:"duration"`
	Minutes  float64       `json:"minutes"`
	Hours    float64       `json:"hours"`
	String   string        `json:"string"`
}
