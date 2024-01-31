package cronjobs

import (
	"time"

	"github.com/pramod-janardhana/dir-watcher/repository/entity"
)

type JobRunDetails struct {
	RunId        string                `json:"runId"`
	JobType      int                   `json:"jobType"`
	StartTime    time.Time             `json:"startTime"`
	EndTime      time.Time             `json:"endTime"`
	TotalRunTime string                `json:"totalRunTime"`
	Status       int                   `json:"status"`
	Details      entity.CronJobDetails `json:"details"`
}

type GetByStartTimeRes struct {
	JobRunDetails
}

type UpdateConfigRes string
