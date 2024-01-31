package database

import (
	"context"
	"database/sql"
)

type Database interface {
	GetConnection() *sql.DB
	Close() error
	ExecTx(ctx context.Context, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error
}
