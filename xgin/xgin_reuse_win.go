//go:build windows
// +build windows

package xgin

import (
	"errors"
	"golang.org/x/sys/windows"
	"syscall"
)

// ReuseListenControl can be used as Control field in net.ListenConfig, to reuse port when listening tcp or udp port.
// In Windows, this function will enable SO_REUSEADDR flag to current socket fd.
func ReuseListenControl(_, _ string, c syscall.RawConn) error {
	var setoptErr = errors.New("failed to call 'Control' method on syscall.RawConn")
	fn := func(fd uintptr) {
		setoptErr = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
	}
	_ = c.Control(fn) // almost no error
	return setoptErr
}
