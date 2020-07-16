package xgorm

import (
	"github.com/go-sql-driver/mysql"
)

// http://go-database-sql.org/errors.html
// Reference from https://github.com/VividCortex/mysqlerr/blob/master/mysqlerr.go
const (
	DuplicateEntryError = 1062
)

// noinspection GoUnusedExportedFunction
func IsMySqlDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == DuplicateEntryError
}
