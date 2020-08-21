# xgorm

### Functions

#### Normal

+ `HookDeleteAtField(db *gorm.DB, defaultDeleteAtTimeStamp string)`
+ `type GormTime struct {}`
+ `type GormTimeWithoutDeletedAt struct {}`
+ `IsMySqlDuplicateEntryError(err error) bool`
+ `type GormLogrus struct {}`
+ `NewGormLogrus(logger *logrus.Logger) *GormLogrus`
+ `type GormLogger struct {}`
+ `NewGormLogger(logger *log.Logger) *GormLogger`

#### Helper

+ `type Helper struct {}`
+ `WithDB(db *gorm.DB) *Helper`
+ `(h *Helper) Pagination(limit int32, page int32) *gorm.DB`
+ `(h *Helper) Count(model interface{}, where interface{}) uint64`
+ `(h *Helper) Exist(model interface{}, where interface{}) bool`
+ `(h *Helper) Create(model interface{}, object interface{}) xstatus.DbStatus`
+ `(h *Helper) Update(model interface{}, where interface{}, object interface{}) xstatus.DbStatus`
+ `(h *Helper) Delete(model interface{}, where interface{}, object interface{}) xstatus.DbStatus`
+ `CreateDB(rdb *gorm.DB) (xstatus.DbStatus, error)`
+ `UpdateDB(rdb *gorm.DB) (xstatus.DbStatus, error)`
+ `DeleteDB(rdb *gorm.DB) (xstatus.DbStatus, error)`
