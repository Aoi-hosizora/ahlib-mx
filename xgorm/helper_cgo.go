//go:build cgo
// +build cgo

package xgorm

import (
	"github.com/mattn/go-sqlite3"
)

// NewSQLiteDriver creates sqlite3.SQLiteDriver with extensions and connect hooks, based on given SQLiteDriverOption-s.
func NewSQLiteDriver(options ...SQLiteDriverOption) *sqlite3.SQLiteDriver {
	opt := buildSQLiteDriverOptions(options)
	sqliteDriver := &sqlite3.SQLiteDriver{
		Extensions: getSQLiteDriverOptionExtensions(opt),
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return applySQLiteDriverOptionRegisterers(
				opt,
				conn.RegisterAggregator,
				conn.RegisterAuthorizer,
				conn.RegisterCollation,
				conn.RegisterCommitHook,
				conn.RegisterFunc,
				func(callback func(interface{})) {
					conn.RegisterPreUpdateHook(func(data sqlite3.SQLitePreUpdateData) { callback(data) })
				},
				conn.RegisterRollbackHook,
				conn.RegisterUpdateHook,
				func(callback func(interface{}) error) error { return callback(conn) },
			)
		},
	}
	return sqliteDriver
}

// IsSQLiteUniqueConstraintError checks whether err is SQLite's ErrConstraintUnique error, whose extended code is SQLiteUniqueConstraintErrno.
func IsSQLiteUniqueConstraintError(err error) bool {
	e, ok := err.(sqlite3.Error)
	if ok {
		return int(e.ExtendedCode) == SQLiteUniqueConstraintErrno
	}
	pe, ok := err.(*sqlite3.Error)
	return ok && int(pe.ExtendedCode) == SQLiteUniqueConstraintErrno
}
