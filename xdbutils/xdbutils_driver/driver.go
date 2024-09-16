package xdbutils_driver

import (
	_ "database/sql"
	"database/sql/driver"
	"sync"
	_ "unsafe"
)

//go:linkname driversMu database/sql.driversMu
//go:linkname drivers database/sql.drivers
var (
	driversMu sync.RWMutex
	drivers   map[string]driver.Driver
)

// GetSQLDriver gets registered driver.Driver by given name, returns nil if unregistered.
func GetSQLDriver(name string) driver.Driver {
	driversMu.Lock()
	defer driversMu.Unlock()
	d, _ := drivers[name]
	return d
}

// ForceRegisterSQLDriver registers given driver.Driver just like sql.Register, but will replace the same-name-registered driver.Driver.
func ForceRegisterSQLDriver(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	drivers[name] = driver
}

// ForceUnregisterSQLDriver unregisters driver.Driver with given name.
func ForceUnregisterSQLDriver(name string) {
	driversMu.Lock()
	defer driversMu.Unlock()
	delete(drivers, name)
}
