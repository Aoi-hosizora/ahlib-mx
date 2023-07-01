//go:build !windows && !unix
// +build !windows,!unix

package xgin

import (
	"syscall"
)

// ReuseListenControl does nothing in current OS.
func ReuseListenControl(_, _ string, c syscall.RawConn) error {
	return nil
}
