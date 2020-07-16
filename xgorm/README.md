# xgorm

### Functions

+ `HookDeleteAtField(db *gorm.DB, defaultDeleteAtTimeStamp string)`
+ `type GormTime struct {}`
+ `type GormTimeWithoutDeletedAt struct {}`
+ `IsMySqlDuplicateEntryError(err error) bool`

### Logger Functions

+ `type GormLogrus struct {}`
+ `NewGormLogrus(logger *logrus.Logger) *GormLogrus`
+ `type GormLogger struct {}`
+ `NewGormLogger(logger *log.Logger) *GormLogger`
