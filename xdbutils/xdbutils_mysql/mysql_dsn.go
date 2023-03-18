package xdbutils_mysql

// TODO add unit test

// DSNFormatter is an interface for types which implement FormatDSN method, such as mysql.Config, xdbutils_sqlite.SQLiteConfig and
// xdbutils_postgres.PostgreSQLConfig.
type DSNFormatter interface {
	FormatDSN() string
}

// MySQLExtraConfig is an extra configuration for MySQL, can be used to generate mysql.Config.Params by ToParams method.
//
// Please visit the follow links for more information:
// - https://github.com/go-sql-driver/mysql#dsn-data-source-name
// - https://dev.mysql.com/doc/refman/8.0/en/connecting-using-uri-or-key-value-pairs.html
type MySQLExtraConfig struct {
	AllowFallbackToPlaintext *bool  // True allowFallbackToPlaintext acts like a --ssl-mode=PREFERRED MySQL client as described
	Charset                  string // Sets the charset used for client-server interaction, prefer "utf8mb4,utf8"
	Autocommit               *bool  // In MySQL, autocommit is default to true
	TimeZone                 string // such as: Europe/Paris, https://dev.mysql.com/doc/refman/8.0/en/time-zone-support.html
	TransactionIsolation     string // such as: REPEATABLE-READ, https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_transaction_isolation
	SQLMode                  string // such as: TRADITIONAL, https://dev.mysql.com/doc/refman/5.7/en/sql-mode.html.
	SysVar                   string // such as: esc@ped

	OtherParams map[string]string // more custom params
}

// ToParams generates a parameter map for mysql.Config.Params from MySQLExtraConfig.
func (m MySQLExtraConfig) ToParams() map[string]string {
	// NOTE: DO NOT set the type of receiver `m` to `*MySQLExtraConfig`, because
	// expression `MySQLExtraConfig{}.ToParams()` is excepted to be used, but if
	// `m` is a pointer value, this expression will have a syntax error (Cannot
	// call a pointer method on '&MySQLExtraConfig{...}').

	result := make(map[string]string)

	if m.AllowFallbackToPlaintext != nil {
		result["allowFallbackToPlaintext"] = boolString(*m.AllowFallbackToPlaintext)
	}
	if m.Charset != "" {
		result["charset"] = m.Charset
	}
	if m.Autocommit != nil {
		result["autocommit"] = boolString(*m.Autocommit)
	}
	if m.TimeZone != "" {
		result["time_zone"] = m.TimeZone
	}
	if m.TransactionIsolation != "" {
		result["transaction_isolation"] = m.TransactionIsolation
	}
	if m.SQLMode != "" {
		result["sql_mode"] = m.SQLMode
	}
	if m.SysVar != "" {
		result["sys_var"] = m.SysVar
	}

	for k, v := range m.OtherParams {
		result[k] = v
	}
	return result
}

// boolString returns string from bool value, in "0" and "1" format.
func boolString(b bool) string {
	if b {
		return "1" // or true
	}
	return "0" // or false
}
