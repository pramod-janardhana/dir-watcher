//go:build cgo
// +build cgo

package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func open(conStr string) (*sql.DB, error) {
	return sql.Open("sqlite3", conStr)
}
