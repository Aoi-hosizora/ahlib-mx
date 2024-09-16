package xdbutils_sqlite

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/internal"
	"testing"
)

func TestSQLiteConfig(t *testing.T) {
	// TODO SQLiteConfig
}

type stringError string
type structError struct{ msg string }
type fakeSQLiteError struct{ ExtendedCode ErrNoExtended }

func (n stringError) Error() string     { return string(n) }
func (n structError) Error() string     { return n.msg }
func (f fakeSQLiteError) Error() string { return "x" }

func TestCheckSQLiteErrorExtendedCodeByReflect(t *testing.T) {
	for _, tc := range []struct {
		giveError error
		giveCode  int
		want      bool
	}{
		{nil, 0, false},
		{errors.New("x"), 0, false},
		{stringError("x"), 0, false},
		{(*structError)(nil), 0, false},
		{structError{"x"}, 0, false},
		{&structError{"x"}, 0, false},
		{fakeSQLiteError{0}, 0, true},
		{&fakeSQLiteError{0}, 0, true},
		{fakeSQLiteError{ErrConstraintUnique}, 0, false},
		{&fakeSQLiteError{ErrConstraintUnique}, 0, false},
		{fakeSQLiteError{ErrConstraintUnique}, int(ErrConstraintUnique), true},
		{&fakeSQLiteError{ErrConstraintUnique}, int(ErrConstraintUnique), true},
	} {
		internal.XtestingEqual(t, CheckSQLiteErrorExtendedCodeByReflect(tc.giveError, tc.giveCode), tc.want)
	}
}

func TestErrNo(t *testing.T) {
	// https://github.com/mattn/go-sqlite3/blob/85a15a7254/error.go#L37
	errnos := []ErrNo{
		ErrError, ErrInternal, ErrPerm, ErrAbort, ErrBusy, ErrLocked, ErrNomem, ErrReadonly, ErrInterrupt, ErrIoErr, ErrCorrupt, ErrNotFound, ErrFull, ErrCantOpen,
		ErrProtocol, ErrEmpty, ErrSchema, ErrTooBig, ErrConstraint, ErrMismatch, ErrMisuse, ErrNoLFS, ErrAuth, ErrFormat, ErrRange, ErrNotADB, ErrNotice, ErrWarning,
	}
	for i := 1; i <= 28; i++ {
		internal.XtestingEqual(t, errnos[i-1], ErrNo(i))
	}

	// https://github.com/mattn/go-sqlite3/blob/85a15a7254/error.go#L97
	for _, tc := range []struct {
		giveExtend ErrNoExtended
		wantErrNo  ErrNo
		wantBy     int
	}{
		// ErrIoErr = ErrNo(10)
		{ErrIoErrRead, ErrIoErr, 1},
		{ErrIoErrShortRead, ErrIoErr, 2},
		{ErrIoErrWrite, ErrIoErr, 3},
		{ErrIoErrFsync, ErrIoErr, 4},
		{ErrIoErrDirFsync, ErrIoErr, 5},
		{ErrIoErrTruncate, ErrIoErr, 6},
		{ErrIoErrFstat, ErrIoErr, 7},
		{ErrIoErrUnlock, ErrIoErr, 8},
		{ErrIoErrRDlock, ErrIoErr, 9},
		{ErrIoErrDelete, ErrIoErr, 10},
		{ErrIoErrBlocked, ErrIoErr, 11},
		{ErrIoErrNoMem, ErrIoErr, 12},
		{ErrIoErrAccess, ErrIoErr, 13},
		{ErrIoErrCheckReservedLock, ErrIoErr, 14},
		{ErrIoErrLock, ErrIoErr, 15},
		{ErrIoErrClose, ErrIoErr, 16},
		{ErrIoErrDirClose, ErrIoErr, 17},
		{ErrIoErrSHMOpen, ErrIoErr, 18},
		{ErrIoErrSHMSize, ErrIoErr, 19},
		{ErrIoErrSHMLock, ErrIoErr, 20},
		{ErrIoErrSHMMap, ErrIoErr, 21},
		{ErrIoErrSeek, ErrIoErr, 22},
		{ErrIoErrDeleteNoent, ErrIoErr, 23},
		{ErrIoErrMMap, ErrIoErr, 24},
		{ErrIoErrGetTempPath, ErrIoErr, 25},
		{ErrIoErrConvPath, ErrIoErr, 26},

		// ErrLocked = ErrNo(6)
		{ErrLockedSharedCache, ErrLocked, 1},

		// ErrBusy = ErrNo(5)
		{ErrBusyRecovery, ErrBusy, 1},
		{ErrBusySnapshot, ErrBusy, 2},

		// ErrCantOpen = ErrNo(14)
		{ErrCantOpenNoTempDir, ErrCantOpen, 1},
		{ErrCantOpenIsDir, ErrCantOpen, 2},
		{ErrCantOpenFullPath, ErrCantOpen, 3},
		{ErrCantOpenConvPath, ErrCantOpen, 4},

		// ErrCorrupt = ErrNo(11)
		{ErrCorruptVTab, ErrCorrupt, 1},

		// ErrReadonly = ErrNo(8)
		{ErrReadonlyRecovery, ErrReadonly, 1},
		{ErrReadonlyCantLock, ErrReadonly, 2},
		{ErrReadonlyRollback, ErrReadonly, 3},
		{ErrReadonlyDbMoved, ErrReadonly, 4},

		// ErrAbort = ErrNo(4)
		{ErrAbortRollback, ErrAbort, 2},

		// ErrConstraint = ErrNo(19)
		{ErrConstraintCheck, ErrConstraint, 1},
		{ErrConstraintCommitHook, ErrConstraint, 2},
		{ErrConstraintForeignKey, ErrConstraint, 3},
		{ErrConstraintFunction, ErrConstraint, 4},
		{ErrConstraintNotNull, ErrConstraint, 5},
		{ErrConstraintPrimaryKey, ErrConstraint, 6},
		{ErrConstraintTrigger, ErrConstraint, 7},
		{ErrConstraintUnique, ErrConstraint, 8},
		{ErrConstraintVTab, ErrConstraint, 9},
		{ErrConstraintRowID, ErrConstraint, 10},

		// ErrNotice = ErrNo(27)
		{ErrNoticeRecoverWAL, ErrNotice, 1},
		{ErrNoticeRecoverRollback, ErrNotice, 2},

		// ErrWarning = ErrNo(28)
		{ErrWarningAutoIndex, ErrWarning, 1},
	} {
		internal.XtestingEqual(t, int(tc.giveExtend), int(tc.wantErrNo.Extend(tc.wantBy)))
	}
}

func TestSQLiteDriverOption(t *testing.T) {
	internal.XtestingPanic(t, false, func() { buildSQLiteDriverOptions([]SQLiteDriverOption{}) })
	internal.XtestingPanic(t, false, func() { buildSQLiteDriverOptions(nil) })

	internal.XtestingEqual(t, getSQLiteDriverOptionExtensions(buildSQLiteDriverOptions(nil)), []string(nil))
	internal.XtestingEqual(t, getSQLiteDriverOptionExtensions(buildSQLiteDriverOptions(
		[]SQLiteDriverOption{WithExtensions([]string{"1", "2"})})), []string{"1", "2"})
	internal.XtestingEqual(t, getSQLiteDriverOptionExtensions(buildSQLiteDriverOptions(
		[]SQLiteDriverOption{WithExtensions([]string{"1", "2"}), WithExtensions([]string{})})), []string{})

	internal.XtestingEqual(t, applySQLiteDriverOptionRegisterers(buildSQLiteDriverOptions(nil),
		nil, nil, nil, nil, nil, nil,
		nil, nil, nil), error(nil))

	callback := 0
	opt := buildSQLiteDriverOptions(
		[]SQLiteDriverOption{
			WithAggregatorRegisterer("name", "impl", true),
			WithAggregatorRegisterer("name2", "impl2", false),
			WithAuthorizerRegisterer(func(i int, s string, s2 string, s3 string) int { callback += 1; return i + 1 }),
			WithAuthorizerRegisterer(func(i int, s string, s2 string, s3 string) int { callback += 1; return i + 2 }),
			WithCollationRegisterer("name", func(s string, s2 string) int { callback += 1; return len(s) + 3 }),
			WithCollationRegisterer("name2", func(s string, s2 string) int { callback += 1; return len(s) + 4 }),
			WithCommitHookRegisterer(func() int { callback += 1; return 5 }),
			WithCommitHookRegisterer(func() int { callback += 1; return 6 }),
			WithFuncRegisterer("name", "impl", true),
			WithFuncRegisterer("name2", "impl2", false),
			WithPreUpdateHookRegister(func(i interface{}) { callback += 1; *i.(*int) = 7 }),
			WithPreUpdateHookRegister(func(i interface{}) { callback += 1; *i.(*int) = 8 }),
			WithRollbackHookRegisterer(func() { callback += 1 }),
			WithRollbackHookRegisterer(func() { callback += 1 }),
			WithUpdateHookRegisterer(func(i int, s string, s2 string, i2 int64) { callback += 1 }),
			WithUpdateHookRegisterer(func(i int, s string, s2 string, i2 int64) { callback += 1 }),
			WithSQLiteConnectionHooker(func(i interface{}) error { callback += 1; *i.(*int) = 9; return nil }),
			WithSQLiteConnectionHooker(func(i interface{}) error { callback += 1; *i.(*int) = 10; return nil }),
		})
	for _, tc := range []struct {
		name    string
		err1    error
		err2    error
		err3    error
		err4    error
		wantErr bool
		wantCnt int
	}{
		{"normal", nil, nil, nil, nil, false, 14},
		{"err1", errors.New("err1"), nil, nil, nil, true, 0},
		{"err2", nil, errors.New("err2"), nil, nil, true, 3},
		{"err3", nil, nil, errors.New("err3"), nil, true, 6},
		{"err4", nil, nil, nil, errors.New("err4"), true, 13},
	} {
		t.Run(tc.name, func(t *testing.T) {
			callback = 0
			err := applySQLiteDriverOptionRegisterers(
				opt,
				func(name string, impl interface{}, pure bool) error {
					internal.XtestingEqual(t, (name == "name" && impl == "impl" && pure) || (name == "name2" && impl == "impl2" && !pure), true)
					return tc.err1
				},
				func(callback func(int, string, string, string) int) {
					res := callback(1, "", "", "")
					internal.XtestingEqual(t, res == 2 || res == 3, true)
				},
				func(name string, cmp func(string, string) int) error {
					internal.XtestingEqual(t, name == "name" || name == "name2", true)
					res := cmp(name, name)
					internal.XtestingEqual(t, res == 7 || res == 9, true)
					return tc.err2
				},
				func(callback func() int) {
					res := callback()
					internal.XtestingEqual(t, res == 5 || res == 6, true)
				},
				func(name string, impl interface{}, pure bool) error {
					internal.XtestingEqual(t, (name == "name" && impl == "impl" && pure) || (name == "name2" && impl == "impl2" && !pure), true)
					return tc.err3
				},
				func(callback func(interface{})) {
					ptr := new(int)
					callback(ptr)
					internal.XtestingEqual(t, *ptr == 7 || *ptr == 8, true)
				},
				func(callback func()) { callback() },
				func(callback func(int, string, string, int64)) { callback(0, "", "", 0) },
				func(callback func(interface{}) error) error {
					ptr := new(int)
					internal.XtestingEqual(t, callback(ptr), error(nil))
					internal.XtestingEqual(t, *ptr == 9 || *ptr == 10, true)
					return tc.err4
				},
			)
			internal.XtestingEqual(t, err != nil, tc.wantErr)
			internal.XtestingEqual(t, callback, tc.wantCnt)
		})
	}
}
