package entity

import (
	"encoding/json"
	"time"
)

type Type int8

const (
	_ Type = iota
	CRON_JOB
)

type Status int8

const (
	_ Status = iota
	IN_PROGRESS
	FAILED
	COMPLETED
)

type CronJobDetails struct {
	MagicStringCount  int64 `json:"magicStringCount"`
	DeletedFilesCount int64 `json:"deletedFilesCount"`
	CreatedFilesCount int64 `json:"createdFilesCount"`
}

type JobRuns struct {
	RunId     string
	JobType   Type
	StartTime time.Time
	EndTime   time.Time
	Status    Status
	Details   json.RawMessage
}

func (e *JobRuns) TableName() string {
	return "job_runs"
}

func (e *JobRuns) ColumnRunId() string {
	return "run_id"
}

func (e *JobRuns) ColumnJobType() string {
	return "job_type"
}

func (e *JobRuns) ColumnDetails() string {
	return "details"
}

func (e *JobRuns) ColumnStartTime() string {
	return "start_time"
}

func (e *JobRuns) ColumnEndTime() string {
	return "end_time"
}

func (e *JobRuns) ColumnStatus() string {
	return "status"
}

type JobRunsStore interface {
	Upset(f JobRuns) error
	Query(query string) ([]*JobRuns, error)
}
