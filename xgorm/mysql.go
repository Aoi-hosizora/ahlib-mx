package xgorm

import (
	"github.com/go-sql-driver/mysql"
)

func IsMySqlDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}
