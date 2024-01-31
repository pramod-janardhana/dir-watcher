package datastore

import (
	"bytes"
	"fmt"

	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/repository/entity"
)

type jobRunsStore struct {
	db    database.Database
	model entity.JobRuns
}

func NewJobRunsStore(db database.Database) entity.JobRunsStore {
	return &jobRunsStore{db: db, model: entity.JobRuns{}}
}

func (s *jobRunsStore) Upset(entity entity.JobRuns) error {
	query := fmt.Sprintf(`INSERT INTO %s(%s, %s, %s, %s, %s, %s) 
		VALUES('%s', %d, '%s', '%s', %d, '%s')
		ON CONFLICT(%s) 
		DO UPDATE SET %s=%d, %s='%s', %s='%s';`,
		entity.TableName(), entity.ColumnRunId(), entity.ColumnJobType(), entity.ColumnStartTime(), entity.ColumnEndTime(), entity.ColumnStatus(), entity.ColumnDetails(),
		entity.RunId, entity.JobType, entity.StartTime.UTC(), entity.EndTime.UTC(), entity.Status, string(entity.Details),
		entity.ColumnRunId(),
		entity.ColumnStatus(), entity.Status, entity.ColumnEndTime(), entity.EndTime.UTC(), entity.ColumnDetails(), string(entity.Details),
	)
	_, err := s.db.GetConnection().Exec(query)
	return err
}

func (s *jobRunsStore) Query(query string) ([]*entity.JobRuns, error) {
	rows, err := s.db.GetConnection().Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobRuns := make([]*entity.JobRuns, 0)
	for rows.Next() {
		var record entity.JobRuns
		var details string
		err := rows.Scan(&record.RunId, &record.JobType, &record.StartTime, &record.EndTime, &record.Status, &details)
		if err != nil {
			return nil, err
		}

		record.Details = bytes.NewBufferString(details).Bytes()
		jobRuns = append(jobRuns, &record)
	}

	return jobRuns, nil
}
