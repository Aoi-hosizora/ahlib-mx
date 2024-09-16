package xgormv2

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_driver"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_mysql"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_orderby"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_postgres"
	"github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_sqlite"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	_ "unsafe"
)

// ========
// dialects
// ========

const (
	// MySQL represents the "mysql" dialect for gorm. Remember to use github.com/go-gorm/mysql to open a gorm.DB for MySQL.
	MySQL = "mysql"

	// SQLite represents the "sqlite" dialect for gorm. Remember to use github.com/go-gorm/sqlite to open a gorm.DB SQLite.
	SQLite = "sqlite"

	// Postgres represents the "postgres" dialect for gorm. Remember to use github.com/go-gorm/postgres to open a gorm.DB for PostgreSQL.
	Postgres = "postgres"
)

// IsMySQL checks whether the dialect of given gorm.DB is "mysql".
func IsMySQL(db *gorm.DB) bool {
	return db.Dialector.Name() == MySQL
}

// IsSQLite checks whether the dialect of given gorm.DB is "sqlite".
func IsSQLite(db *gorm.DB) bool {
	return db.Dialector.Name() == SQLite
}

// IsPostgreSQL checks whether the dialect of given gorm.DB is "postgres".
func IsPostgreSQL(db *gorm.DB) bool {
	return db.Dialector.Name() == Postgres
}

// MySQLConfig is a configuration for MySQL, can be used to generate DSN by FormatDSN method.
type MySQLConfig = mysql.Config

// MySQLExtraConfig is an extra configuration for MySQL, can be used to generate extends given param by ToParams method.
type MySQLExtraConfig = xdbutils_mysql.MySQLExtraConfig

// SQLiteConfig is a configuration for SQLite, can be used to generate DSN by FormatDSN method.
type SQLiteConfig = xdbutils_sqlite.SQLiteConfig

// PostgreSQLConfig is a configuration for PostgreSQL, can be used to generate DSN by FormatDSN method.
type PostgreSQLConfig = xdbutils_postgres.PostgreSQLConfig

// MySQLDefaultDsn returns the MySQL dsn from given parameters with "utf8mb4" charset and "local" location.
//
// Please visit the follow links for more information:
// - https://github.com/go-sql-driver/mysql#dsn-data-source-name
// - https://dev.mysql.com/doc/refman/8.0/en/connecting-using-uri-or-key-value-pairs.html
func MySQLDefaultDsn(username, password, address, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, database)
}

// SQLiteDefaultDsn returns the SQLite dsn from given parameter (database filename or ":memory:" or empty string).
//
// Please visit the follow links for more information:
// - https://github.com/mattn/go-sqlite3#connection-string
// - https://www.sqlite.org/c3ref/open.html
func SQLiteDefaultDsn(file string) string {
	return file
}

// PostgreSQLDefaultDsn returns the PostgreSQL dsn from given parameters.
//
// Please visit the follow links for more information:
// - https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters
// - https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
// - https://www.postgresql.org/docs/current/runtime-config-client.html
func PostgreSQLDefaultDsn(username, password, host string, port int, database string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, username, password, database)
}

const (
	// MySQLDuplicateEntryErrno is MySQL's DUP_ENTRY errno, referred from https://github.com/VividCortex/mysqlerr/blob/69f897f9a2/mysqlerr.go and
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.htm.
	MySQLDuplicateEntryErrno = uint16(mysqlerr.ER_DUP_ENTRY) // 1062, DUP_ENTRY

	// SQLiteUniqueConstraintErrno is SQLite's CONSTRAINT_UNIQUE extended errno, referred from https://github.com/mattn/go-sqlite3/blob/85a15a7254/error.go
	// and http://www.sqlite.org/c3ref/c_abort_rollback.html.
	SQLiteUniqueConstraintErrno = int(xdbutils_sqlite.ErrConstraintUnique) // 19 | 8<<8, sqlite3.ErrConstraintUnique

	// PostgreSQLUniqueViolationErrno is PostgreSQL's unique_violation errno, referred from https://github.com/lib/pq/blob/89fee89644/error.go and
	// https://www.postgresql.org/docs/10/errcodes-appendix.html
	PostgreSQLUniqueViolationErrno = "23505" // pq.errorCodeNames unique_violation
)

// IsMySQLDuplicateEntryError checks whether err is MySQL's ER_DUP_ENTRY error, whose error code is MySQLDuplicateEntryErrno.
func IsMySQLDuplicateEntryError(err error) bool {
	e, ok := err.(*mysql.MySQLError)
	return ok && e.Number == MySQLDuplicateEntryErrno
}

// IsPostgreSQLUniqueViolationError is a variable that used to check whether err is PostgreSQL's unique_violation error, whose error code is PostgreSQLUniqueViolationErrno.
//
// You have to import "github.com/jackc/pgconn" and set this variable manually in order to make it be usable in CreateErr and UpdateErr.
//
//	xgormv2.IsPostgreSQLUniqueViolationError = func(err error) bool {
//		perr, ok := err.(*pgconn.PgError)
//		return ok && perr.Code == xgormv2.PostgreSQLUniqueViolationErrno
//	}
var IsPostgreSQLUniqueViolationError func(err error) bool

// ===============
// CRUD and others
// ===============

// IsRecordNotFound checks whether given error from gorm.DB is gorm.ErrRecordNotFound.
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// QueryErr checks gorm.DB after query operated, will only return xstatus.DbSuccess, xstatus.DbNotFound and xstatus.DbFailed.
func QueryErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	switch {
	case IsRecordNotFound(rdb.Error):
		return xstatus.DbNotFound, nil // not found
	case rdb.Error != nil:
		return xstatus.DbFailed, rdb.Error // failed
	}
	return xstatus.DbSuccess, nil
}

// CreateErr checks gorm.DB after create operated, will only return xstatus.DbSuccess, xstatus.DbExisted and xstatus.DbFailed.
func CreateErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	switch {
	case IsMySQL(rdb) && IsMySQLDuplicateEntryError(rdb.Error),
		IsSQLite(rdb) && IsSQLiteUniqueConstraintError(rdb.Error),
		IsPostgreSQLUniqueViolationError != nil && IsPostgreSQL(rdb) && IsPostgreSQLUniqueViolationError(rdb.Error):
		return xstatus.DbExisted, rdb.Error // duplicate
	case rdb.Error != nil:
		return xstatus.DbFailed, rdb.Error // failed
	}
	return xstatus.DbSuccess, nil
}

// UpdateErr checks gorm.DB after update operated, will only return xstatus.DbSuccess, xstatus.DbNotFound, xstatus.DbExisted and xstatus.DbFailed.
func UpdateErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	switch {
	case IsMySQL(rdb) && IsMySQLDuplicateEntryError(rdb.Error),
		IsSQLite(rdb) && IsSQLiteUniqueConstraintError(rdb.Error),
		IsPostgreSQLUniqueViolationError != nil && IsPostgreSQL(rdb) && IsPostgreSQLUniqueViolationError(rdb.Error):
		return xstatus.DbExisted, rdb.Error // duplicate
	case rdb.Error != nil:
		return xstatus.DbFailed, rdb.Error // failed
	case rdb.RowsAffected == 0:
		return xstatus.DbNotFound, nil // not found
	}
	return xstatus.DbSuccess, nil
}

// DeleteErr checks gorm.DB after delete operated, will only return xstatus.DbSuccess, xstatus.DbNotFound and xstatus.DbFailed.
func DeleteErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	switch {
	case rdb.Error != nil:
		return xstatus.DbFailed, rdb.Error // failed
	case rdb.RowsAffected == 0:
		return xstatus.DbNotFound, nil // not found
	}
	return xstatus.DbSuccess, nil
}

// PropertyValue represents database single entity's property mapping rule, is used in GenerateOrderByExpr.
type PropertyValue = xdbutils_orderby.PropertyValue

// PropertyDict is used to store PropertyValue-s for data transfer object (dto) to entity's property mapping rule, is used in GenerateOrderByExpr.
type PropertyDict = xdbutils_orderby.PropertyDict

// OrderByOption represents an option type for GenerateOrderByExpr's option, can be created by WithXXX functions.
type OrderByOption = xdbutils_orderby.OrderByOption

// WithOrderBySourceSeparator creates an OrderByOption to specify the source order-by expression fields separator, defaults to ",".
func WithOrderBySourceSeparator(separator string) OrderByOption {
	return xdbutils_orderby.WithSourceSeparator(separator)
}

// WithOrderByTargetSeparator creates an OrderByOption to specify the target order-by expression fields separator, defaults to ", ".
func WithOrderByTargetSeparator(separator string) OrderByOption {
	return xdbutils_orderby.WithTargetSeparator(separator)
}

// WithOrderBySourceProcessor creates an OrderByOption to specify the source processor for extracting field name and ascending flag from given source,
// defaults to use the "field asc" or "field desc" format (case-insensitive) to extract information.
func WithOrderBySourceProcessor(processor func(source string) (field string, asc bool)) OrderByOption {
	return xdbutils_orderby.WithSourceProcessor(processor)
}

// WithOrderByTargetProcessor creates an OrderByOption to specify the target processor for combining field name and ascending flag to target expression,
// defaults to generate the target with "destination ASC" or "destination DESC" format.
func WithOrderByTargetProcessor(processor func(destination string, asc bool) (target string)) OrderByOption {
	return xdbutils_orderby.WithTargetProcessor(processor)
}

// WithDefaultExpression creates an OrderByOption to specify the default order expression when generated result is empty,
// defaults to empty.
func WithDefaultExpression(defaultExpression string) OrderByOption {
	return xdbutils_orderby.WithDefaultExpression(defaultExpression)
}

// NewPropertyValue creates a PropertyValue by given reverse and destinations, is used to describe database single entity's property mapping rule.
//
// Here:
// 1. `destinations` represents mapping property destination list, use `property_name` directly for sql, use `returned_name.property_name` for cypher.
// 2. `reverse` represents the flag whether you need to revert the order or not.
func NewPropertyValue(reverse bool, destinations ...string) *PropertyValue {
	return xdbutils_orderby.NewPropertyValue(reverse, destinations...)
}

// GenerateOrderByExpr returns a generated order-by expression by given order-by query source string (such as "name desc, age asc") and PropertyDict,
// with some OrderByOption-s. The generated expression will be in mysql-sql (such as "xxx ASC") or neo4j-cypher style (such as "xxx.yyy DESC").
//
// Example:
//
//	dict := PropertyDict{
//		"uid":  NewPropertyValue(false, "uid"),
//		"name": NewPropertyValue(false, "firstname", "lastname"),
//		"age":  NewPropertyValue(true, "birthday"),
//	}
//	_ = GenerateOrderByExpr(`uid, age desc`, dict) // => uid ASC, birthday ASC
//	_ = GenerateOrderByExpr(`age, username desc`, dict) // => birthday DESC, firstname DESC, lastname DESC
func GenerateOrderByExpr(querySource string, dict PropertyDict, options ...OrderByOption) string {
	return xdbutils_orderby.GenerateOrderByExpr(querySource, dict, options...)
}

// =======
// drivers
// =======

// GetSQLDriver gets registered driver.Driver by given name, returns nil if unregistered.
func GetSQLDriver(name string) driver.Driver {
	return xdbutils_driver.GetSQLDriver(name)
}

// ForceRegisterSQLDriver registers given driver.Driver just like sql.Register, but will replace the same-name-registered driver.Driver.
func ForceRegisterSQLDriver(name string, driver driver.Driver) {
	xdbutils_driver.ForceRegisterSQLDriver(name, driver)
}

// ForceUnregisterSQLDriver unregisters driver.Driver with given name.
func ForceUnregisterSQLDriver(name string) {
	xdbutils_driver.ForceUnregisterSQLDriver(name)
}

// SQLiteDriverOption represents an option type for NewSQLiteDriver's option, can be created by WithXXX functions.
type SQLiteDriverOption = xdbutils_sqlite.SQLiteDriverOption

//go:linkname buildSQLiteDriverOptions github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_sqlite.buildSQLiteDriverOptions
//goland:noinspection GoUnusedParameter
func buildSQLiteDriverOptions(options []SQLiteDriverOption) interface{}

//go:linkname getSQLiteDriverOptionExtensions github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_sqlite.getSQLiteDriverOptionExtensions
//goland:noinspection GoUnusedParameter
func getSQLiteDriverOptionExtensions(v interface{}) []string

//go:linkname applySQLiteDriverOptionRegisterers github.com/Aoi-hosizora/ahlib-mx/xdbutils/xdbutils_sqlite.applySQLiteDriverOptionRegisterers
//goland:noinspection GoUnusedParameter
func applySQLiteDriverOptionRegisterers(
	v interface{},
	aggregatorRegisterer func(name string, impl interface{}, pure bool) error,
	authorizerRegisterer func(callback func(int, string, string, string) int),
	collationRegisterer func(name string, cmp func(string, string) int) error,
	commitHookRegisterer func(callback func() int),
	funcRegisterer func(name string, impl interface{}, pure bool) error,
	preUpdateHookRegisterer func(callback func(interface{})),
	rollbackHookRegisterer func(callback func()),
	updateHookRegisterer func(callback func(int, string, string, int64)),
	sqliteConnectionHooker func(callback func(interface{}) error) error,
) error

// WithExtensions creates an SQLiteDriverOption to specify the extensions for sqlite driver.
func WithExtensions(extensions []string) SQLiteDriverOption {
	return xdbutils_sqlite.WithExtensions(extensions)
}

// WithAggregatorRegisterer creates an SQLiteDriverOption to specify the aggregator registerer for sqlite driver.
func WithAggregatorRegisterer(name string, impl interface{}, pure bool) SQLiteDriverOption {
	return xdbutils_sqlite.WithAggregatorRegisterer(name, impl, pure)
}

// WithAuthorizerRegisterer creates an SQLiteDriverOption to specify the authorizer registerer for sqlite driver.
func WithAuthorizerRegisterer(callback func(int, string, string, string) int) SQLiteDriverOption {
	return xdbutils_sqlite.WithAuthorizerRegisterer(callback)
}

// WithCollationRegisterer creates an SQLiteDriverOption to specify the collation registerer for sqlite driver.
func WithCollationRegisterer(name string, cmp func(string, string) int) SQLiteDriverOption {
	return xdbutils_sqlite.WithCollationRegisterer(name, cmp)
}

// WithCommitHookRegisterer creates an SQLiteDriverOption to specify the commit hook registerer for sqlite driver.
func WithCommitHookRegisterer(callback func() int) SQLiteDriverOption {
	return xdbutils_sqlite.WithCommitHookRegisterer(callback)
}

// WithFuncRegisterer creates an SQLiteDriverOption to specify the func registerer for sqlite driver.
func WithFuncRegisterer(name string, impl interface{}, pure bool) SQLiteDriverOption {
	return xdbutils_sqlite.WithFuncRegisterer(name, impl, pure)
}

// WithPreUpdateHookRegister creates an SQLiteDriverOption to specify the pre-update hook registerer for sqlite driver.
func WithPreUpdateHookRegister(callback func(interface{})) SQLiteDriverOption {
	return xdbutils_sqlite.WithPreUpdateHookRegister(callback)
}

// WithRollbackHookRegisterer creates an SQLiteDriverOption to specify the rollback hook registerer for sqlite driver.
func WithRollbackHookRegisterer(callback func()) SQLiteDriverOption {
	return xdbutils_sqlite.WithRollbackHookRegisterer(callback)
}

// WithUpdateHookRegisterer creates an SQLiteDriverOption to specify the update hook registerer for sqlite driver.
func WithUpdateHookRegisterer(callback func(int, string, string, int64)) SQLiteDriverOption {
	return xdbutils_sqlite.WithUpdateHookRegisterer(callback)
}

// WithSQLiteConnectionHooker creates an SQLiteDriverOption to specify the sqlite connection hooker for sqlite driver.
func WithSQLiteConnectionHooker(hooker func(interface{}) error) SQLiteDriverOption {
	return xdbutils_sqlite.WithSQLiteConnectionHooker(hooker)
}
