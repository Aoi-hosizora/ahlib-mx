//go:build cgo
// +build cgo

package xgorm

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres" // dummy
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mattn/go-sqlite3"
	"math/rand"
	"testing"
)

func TestMess2(t *testing.T) {
	xtesting.True(t, IsSQLiteUniqueConstraintError(sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(SQLiteUniqueConstraintErrno)}))
	xtesting.True(t, IsSQLiteUniqueConstraintError(&sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(SQLiteUniqueConstraintErrno)}))

	xtesting.NotPanic(t, func() {
		driver := NewSQLiteDriver(WithSQLiteConnectionHooker(func(i interface{}) error {
			return errors.New("test err")
		}))
		err := driver.ConnectHook(&sqlite3.SQLiteConn{})
		xtesting.Equal(t, err.Error(), "test err")
		ForceRegisterSQLDriver(SQLite, driver)
		ForceUnregisterSQLDriver(SQLite)
		xtesting.Nil(t, GetSQLDriver(SQLite))
	})

	xtesting.NotPanic(t, func() {
		NewSQLiteDriver(WithCollationRegisterer("random", func(s string, s2 string) int {
			return rand.Intn(3) - 1
		}))
	})
}

func TestHook(t *testing.T) {
	for _, tc := range []struct {
		giveDialect string
		giveParam   string
	}{
		{MySQL, mysqlDsn},
		{SQLite, sqliteFile},
	} {
		t.Run(tc.giveDialect, func(t *testing.T) {
			testHook(t, tc.giveDialect, tc.giveParam)
		})
	}
}

func TestHelper(t *testing.T) {
	for _, tc := range []struct {
		giveDialect string
		giveParam   string
	}{
		{MySQL, mysqlDsn},
		{SQLite, sqliteFile},
	} {
		t.Run(tc.giveDialect, func(t *testing.T) {
			testHelper(t, tc.giveDialect, tc.giveParam)
		})
	}
}

func TestLogger(t *testing.T) {
	for _, tc := range []struct {
		giveDialect string
		giveParam   string
	}{
		{MySQL, mysqlDsn},
		{SQLite, sqliteFile},
	} {
		t.Run(tc.giveDialect, func(t *testing.T) {
			testLogger(t, tc.giveDialect, tc.giveParam)
		})
	}
}
