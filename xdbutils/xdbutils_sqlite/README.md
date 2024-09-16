# xdbutils_sqlite

## Dependencies

+ None

## Documents

### Types

+ `type ErrNo int`
+ `type ErrNoExtended int`
+ `type SQLiteConfig struct`
+ `type SQLiteDriverOptions struct`

### Variables

+ None

### Constants

+ ...

### Functions

+ `func CheckSQLiteErrorExtendedCodeByReflect(err error, code int) bool`
+ `func WithExtensions(extensions []string) SQLiteDriverOption`
+ `func WithAggregatorRegisterer(name string, impl interface{}, pure bool) SQLiteDriverOption`
+ `func WithAuthorizerRegisterer(callback func(int, string, string, string) int) SQLiteDriverOption`
+ `func WithCollationRegisterer(name string, cmp func(string, string) int) SQLiteDriverOption`
+ `func WithCommitHookRegisterer(callback func() int) SQLiteDriverOption`
+ `func WithFuncRegisterer(name string, impl interface{}, pure bool) SQLiteDriverOption`
+ `func WithPreUpdateHookRegister(callback func(interface{})) SQLiteDriverOption`
+ `func WithRollbackHookRegisterer(callback func()) SQLiteDriverOption`
+ `func WithUpdateHookRegisterer(callback func(int, string, string, int64)) SQLiteDriverOption`
+ `func WithSQLiteConnectionHooker(hooker func(interface{}) error) SQLiteDriverOption`

### Methods

+ `func (err ErrNo) Extend(by int) ErrNoExtended`
+ `func (s *SQLiteConfig) FormatDSN() string`
