package xdbutils_driver

import (
	"database/sql/driver"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/internal"
	"testing"
)

type mockDriver struct{}

var _ = driver.Driver(&mockDriver{})

func (m *mockDriver) Open(string) (driver.Conn, error) {
	return nil, nil
}

const MockDriverName = "mockDriver"

func TestSQLDriver(t *testing.T) {
	internal.XtestingEqual(t, GetSQLDriver(MockDriverName), nil)
	internal.XtestingPanic(t, true, func() {
		ForceRegisterSQLDriver("", nil)
	})

	internal.XtestingPanic(t, false, func() {
		ForceUnregisterSQLDriver(MockDriverName)
	})
	internal.XtestingEqual(t, GetSQLDriver(MockDriverName), nil)

	mocked := &mockDriver{}
	internal.XtestingPanic(t, false, func() {
		ForceRegisterSQLDriver(MockDriverName, mocked)
	})
	internal.XtestingEqual(t, GetSQLDriver(MockDriverName), mocked)

	driver2_ := &mockDriver{}
	internal.XtestingPanic(t, false, func() {
		ForceRegisterSQLDriver(MockDriverName, driver2_)
	})
	internal.XtestingEqual(t, GetSQLDriver(MockDriverName), driver2_)

	internal.XtestingPanic(t, false, func() {
		ForceUnregisterSQLDriver(MockDriverName)
	})
	internal.XtestingEqual(t, GetSQLDriver(MockDriverName), nil)
}
