package xdbutils_postgres

// TODO add unit test

// PostgreSQLConfig is a configuration for PostgreSQL, can be used to generate DSN by FormatDSN method.
//
// Please visit the follow links for more information:
// - https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters
// - https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
// - https://www.postgresql.org/docs/current/runtime-config-client.html
type PostgreSQLConfig struct {
	dbname          string // The name of the database to connect to
	user            string // The user to sign in as
	password        string // The user's password
	host            string // The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
	port            string // The port to bind to. (default is 5432)
	sslmode         string // Whether to use SSL (default is require, this is not the default for libpq)
	connect_timeout string // Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.
	sslcert         string // Cert file location. The file must contain PEM encoded data.
	sslkey          string // Key file location. The file must contain PEM encoded data.
	sslrootcert     string // The location of the root certificate file. The file must contain PEM encoded data.
	// TODO PostgreSQLConfig fields
}

// FormatDSN generates formatted DSN string from PostgreSQLConfig.
func (p *PostgreSQLConfig) FormatDSN() string {
	return ""
}
