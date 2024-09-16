//go:build !cgo
// +build !cgo

package xgorm

import (
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_sqlite"
)

// NewSQLiteDriver creates sqlite3.SQLiteDriver with extensions and connect hooks, based on given SQLiteDriverOption-s.
func NewSQLiteDriver(options ...SQLiteDriverOption) interface{} {
	panic("xgorm: NewSQLiteDriver is not supported without cgo")
}

// IsSQLiteUniqueConstraintError checks whether err is SQLite's ErrConstraintUnique error, whose extended code is SQLiteUniqueConstraintErrno.
func IsSQLiteUniqueConstraintError(err error) bool {
	return xdbutils_sqlite.CheckSQLiteErrorExtendedCodeByReflect(err, SQLiteUniqueConstraintErrno)
}
