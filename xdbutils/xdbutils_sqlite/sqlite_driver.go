package xdbutils_sqlite

// Struct types for SQLiteDriverOptions.
type (
	aggregatorRegisterer struct {
		name string
		impl interface{}
		pure bool
	}
	authorizerRegisterer struct {
		callback func(int, string, string, string) int
	}
	collationRegisterer struct {
		name string
		cmp  func(string, string) int
	}
	commitHookRegisterer struct {
		callback func() int
	}
	funcRegisterer struct {
		name string
		impl interface{}
		pure bool
	}
	preUpdateHookRegisterer struct {
		callback func(interface{}) // sqlite3.SQLitePreUpdateData
	}
	rollbackHookRegisterer struct {
		callback func()
	}
	updateHookRegisterer struct {
		callback func(int, string, string, int64)
	}
	sqliteConnectionHooker struct {
		callback func(interface{}) error // sqlite3.SQLiteConn
	}
)

// SQLiteDriverOptions is a type of NewSQLiteDriver's option, each field can be set by SQLiteDriverOption function type.
type SQLiteDriverOptions struct {
	extensions               []string
	aggregatorRegisterers    []aggregatorRegisterer
	authorizerRegisterers    []authorizerRegisterer
	collationRegisterers     []collationRegisterer
	commitHookRegisterers    []commitHookRegisterer
	funcRegisterers          []funcRegisterer
	preUpdateHookRegisterers []preUpdateHookRegisterer
	rollbackHookRegisterers  []rollbackHookRegisterer
	updateHookRegisterers    []updateHookRegisterer
	sqliteConnectionHookers  []sqliteConnectionHooker
}

// SQLiteDriverOption represents an option type for NewSQLiteDriver's option, can be created by WithXXX functions.
type SQLiteDriverOption func(*SQLiteDriverOptions)

// WithExtensions creates an SQLiteDriverOption to specify the extensions for sqlite driver.
func WithExtensions(extensions []string) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.extensions = extensions
	}
}

// WithAggregatorRegisterer creates an SQLiteDriverOption to specify the aggregator registerer for sqlite driver.
func WithAggregatorRegisterer(name string, impl interface{}, pure bool) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.aggregatorRegisterers = append(o.aggregatorRegisterers, aggregatorRegisterer{name, impl, pure})
	}
}

// WithAuthorizerRegisterer creates an SQLiteDriverOption to specify the authorizer registerer for sqlite driver.
func WithAuthorizerRegisterer(callback func(int, string, string, string) int) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.authorizerRegisterers = append(o.authorizerRegisterers, authorizerRegisterer{callback})
	}
}

// WithCollationRegisterer creates an SQLiteDriverOption to specify the collation registerer for sqlite driver.
func WithCollationRegisterer(name string, cmp func(string, string) int) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.collationRegisterers = append(o.collationRegisterers, collationRegisterer{name, cmp})
	}
}

// WithCommitHookRegisterer creates an SQLiteDriverOption to specify the commit hook registerer for sqlite driver.
func WithCommitHookRegisterer(callback func() int) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.commitHookRegisterers = append(o.commitHookRegisterers, commitHookRegisterer{callback})
	}
}

// WithFuncRegisterer creates an SQLiteDriverOption to specify the func registerer for sqlite driver.
func WithFuncRegisterer(name string, impl interface{}, pure bool) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.funcRegisterers = append(o.funcRegisterers, funcRegisterer{name, impl, pure})
	}
}

// WithPreUpdateHookRegister creates an SQLiteDriverOption to specify the pre-update hook registerer for sqlite driver.
func WithPreUpdateHookRegister(callback func(interface{})) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.preUpdateHookRegisterers = append(o.preUpdateHookRegisterers, preUpdateHookRegisterer{callback})
	}
}

// WithRollbackHookRegisterer creates an SQLiteDriverOption to specify the rollback hook registerer for sqlite driver.
func WithRollbackHookRegisterer(callback func()) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.rollbackHookRegisterers = append(o.rollbackHookRegisterers, rollbackHookRegisterer{callback})
	}
}

// WithUpdateHookRegisterer creates an SQLiteDriverOption to specify the update hook registerer for sqlite driver.
func WithUpdateHookRegisterer(callback func(int, string, string, int64)) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.updateHookRegisterers = append(o.updateHookRegisterers, updateHookRegisterer{callback})
	}
}

// WithSQLiteConnectionHooker creates an SQLiteDriverOption to specify the sqlite connection hooker for sqlite driver.
func WithSQLiteConnectionHooker(hooker func(interface{}) error) SQLiteDriverOption {
	return func(o *SQLiteDriverOptions) {
		o.sqliteConnectionHookers = append(o.sqliteConnectionHookers, sqliteConnectionHooker{hooker})
	}
}

// buildSQLiteDriverOptions creates a SQLiteDriverOptions with given SQLiteDriverOption-s.
//
// Remark: this function will be used by //go:linkname in xgorm and xgormv2 package.
func buildSQLiteDriverOptions(options []SQLiteDriverOption) interface{} {
	opt := &SQLiteDriverOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	return opt
}

// getSQLiteDriverOptionExtensions returns extensions from SQLiteDriverOptions, is an internal helper function and
// use interface{} as SQLiteDriverOptions parameter type.
//
// Remark: this function will be used by //go:linkname in xgorm and xgormv2 package.
func getSQLiteDriverOptionExtensions(v interface{}) []string {
	opt := v.(*SQLiteDriverOptions)
	return opt.extensions
}

// applySQLiteDriverOptionRegisterers apply registerers from SQLiteDriverOptions, is an internal helper function and
// use interface{} as SQLiteDriverOptions parameter type.
//
// Remark: this function will be used by //go:linkname in xgorm and xgormv2 package.
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
) error {
	opt := v.(*SQLiteDriverOptions)
	for _, reg := range opt.aggregatorRegisterers {
		err := aggregatorRegisterer(reg.name, reg.impl, reg.pure)
		if err != nil {
			return err
		}
	}
	for _, reg := range opt.authorizerRegisterers {
		authorizerRegisterer(reg.callback)
	}
	for _, reg := range opt.collationRegisterers {
		err := collationRegisterer(reg.name, reg.cmp)
		if err != nil {
			return err
		}
	}
	for _, reg := range opt.commitHookRegisterers {
		commitHookRegisterer(reg.callback)
	}
	for _, reg := range opt.funcRegisterers {
		err := funcRegisterer(reg.name, reg.impl, reg.pure)
		if err != nil {
			return err
		}
	}
	for _, reg := range opt.rollbackHookRegisterers {
		rollbackHookRegisterer(reg.callback)
	}
	for _, reg := range opt.preUpdateHookRegisterers {
		preUpdateHookRegisterer(reg.callback)
	}
	for _, reg := range opt.updateHookRegisterers {
		updateHookRegisterer(reg.callback)
	}
	for _, reg := range opt.sqliteConnectionHookers {
		err := sqliteConnectionHooker(reg.callback)
		if err != nil {
			return err
		}
	}
	return nil
}
