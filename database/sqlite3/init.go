package sqlite3

import (
	"bytes"
	"io"

	"github.com/pramod-janardhana/dir-watcher/database"
)

func IntiDB(db database.Database, reader io.Reader) error {
	stringBuilter := bytes.Buffer{}
	for {
		buf := make([]byte, 32)
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}

		if _, err := stringBuilter.Write(buf); err != nil {
			return err
		}

	}

	query := stringBuilter.String()
	_, err := db.GetConnection().Exec(query)
	return err
}
