//go:build !cgo
// +build !cgo

package sqlite3

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func open(conStr string) (*sql.DB, error) {
	return sql.Open("sqlite", conStr)
}
