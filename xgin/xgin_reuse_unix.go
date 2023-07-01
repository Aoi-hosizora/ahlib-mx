//go:build unix
// +build unix

package xgin

import (
	"errors"
	"golang.org/x/sys/unix"
	"syscall"
)

// ReuseListenControl can be used as Control field in net.ListenConfig, to reuse port when listening tcp or udp port.
// In *NIX, this function will enable SO_REUSEADDR flag and SO_REUSEPORT flag to current socket fd.
func ReuseListenControl(_, _ string, c syscall.RawConn) error {
	var setoptErr = errors.New("failed to call 'Control' method on syscall.RawConn")
	fn := func(fd uintptr) {
		setoptErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if setoptErr == nil {
			setoptErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		}
	}
	_ = c.Control(fn) // almost no error
	return setoptErr
}
