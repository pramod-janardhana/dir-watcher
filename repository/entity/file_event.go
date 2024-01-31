package entity

type FileEvent struct {
	Id                int8
	DeletedFilesCount int64
	CreatedFilesCount int64
}

func (e *FileEvent) TableName() string {
	return "file_event"
}

func (e *FileEvent) ColumnId() string {
	return "id"
}

func (e *FileEvent) ColumnDeletedFilesCount() string {
	return "deleted_files_count"
}

func (e *FileEvent) ColumnCreatedFilesCount() string {
	return "created_files_count"
}

type FileEventStore interface {
	Upset(f FileEvent) error
	Get(pageSize, pageNumber int) ([]*FileEvent, error)
	Truncate() error
}
