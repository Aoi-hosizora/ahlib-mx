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
