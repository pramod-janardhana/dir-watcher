package entity

import (
	"time"
)

type FileEvent struct {
	Path      string
	Event     int8
	Timestamp time.Time
}

func (e *FileEvent) TableName() string {
	return "file_event"
}

func (e *FileEvent) ColumnPath() string {
	return "path"
}

func (e *FileEvent) ColumnEvent() string {
	return "event"
}

func (e *FileEvent) ColumnTimestamp() string {
	return "timestamp"
}

type FileEventStore interface {
	Upset(f FileEvent) error
	Get(pageSize, pageNumber int) ([]*FileEvent, error)
	Truncate() error
}
