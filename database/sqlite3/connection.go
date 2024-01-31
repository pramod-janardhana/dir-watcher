package sqlite3

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pramod-janardhana/dir-watcher/database"
)

type sqliteDB struct {
	db *sql.DB
}

/*
Open opens a sqlite database specified by its database connection info and returns a Database.
It returns an error if the connection could not be established.
*/
func Open(conStr string) (database.Database, error) {
	db, err := open(conStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	sqliteDB := &sqliteDB{
		db: db,
	}

	return sqliteDB, nil
}

func (d *sqliteDB) GetConnection() *sql.DB {
	return d.db
}

/*
ExecTx begins a transaction and excutes the callback function(fn) followed by committing the transaction.
In case of any errors during transaction, the transaction is rolled back.
*/
func (d *sqliteDB) ExecTx(ctx context.Context, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	tx, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return nil
}

func (d *sqliteDB) Close() error {
	return d.db.Close()
}
