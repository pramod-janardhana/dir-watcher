package datastore

import (
	"fmt"

	"github.com/pramod-janardhana/dir-watcher/database"
	"github.com/pramod-janardhana/dir-watcher/repository/entity"
)

type fileEventStore struct {
	db    database.Database
	model entity.FileEvent
}

func NewFileEventStore(db database.Database) entity.FileEventStore {
	return &fileEventStore{db: db, model: entity.FileEvent{}}
}

func (s *fileEventStore) Upset(entity entity.FileEvent) error {
	query := fmt.Sprintf(`INSERT INTO %s(%s, %s, %s) 
		VALUES('%s', %d, '%s')
		ON CONFLICT(%s) 
		DO UPDATE SET %s=%d, %s='%s';`,
		entity.TableName(), entity.ColumnPath(), entity.ColumnEvent(), entity.ColumnTimestamp(),
		entity.Path, entity.Event, entity.Timestamp.UTC(),
		entity.ColumnPath(),
		entity.ColumnEvent(), entity.Event, entity.ColumnTimestamp(), entity.Timestamp.UTC(),
	)
	_, err := s.db.GetConnection().Exec(query)
	return err
}

func (s *fileEventStore) Get(limit, offset int) ([]*entity.FileEvent, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d",
		s.model.TableName(), limit, offset)

	rows, err := s.db.GetConnection().Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fileEvents := make([]*entity.FileEvent, 0)
	for rows.Next() {
		var record entity.FileEvent
		err := rows.Scan(&record.Path, &record.Event, &record.Timestamp)
		if err != nil {
			return nil, err
		}

		fileEvents = append(fileEvents, &record)
	}

	return fileEvents, nil
}

func (s fileEventStore) Truncate() error {
	_, err := s.db.GetConnection().Exec(fmt.Sprintf("DELETE FROM %s;", s.model.TableName()))
	return err
}
