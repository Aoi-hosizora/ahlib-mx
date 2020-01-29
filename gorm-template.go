package ahlib_gin_gorm

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

const (
	DefaultDeleteAtTimeStamp = "2000-01-01 00:00:00"
)

// default deleteAt at 2000-01-01 00:00:00
type GormTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"default:'2000-01-01 00:00:00'"`
}

type GormTimeWithoutDeletedAt struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func IsMySqlDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}
