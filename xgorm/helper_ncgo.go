//go:build !cgo
// +build !cgo

package xgorm

import (
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_sqlite"
)

// IsSQLiteUniqueConstraintError checks whether err is SQLite's ErrConstraintUnique error, whose extended code is SQLiteUniqueConstraintErrno.
func IsSQLiteUniqueConstraintError(err error) bool {
	return xdbutils_sqlite.CheckSQLiteErrorExtendedCodeByReflect(err, SQLiteUniqueConstraintErrno)
}