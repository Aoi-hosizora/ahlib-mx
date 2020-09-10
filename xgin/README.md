# xgin

### Functions

+ `DumpRequest(c *gin.Context) []string`
+ `BuildBasicErrorDto(err interface{}, c *gin.Context, others map[string]interface{}) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *gin.Context, skip int, print bool) *xdto.ErrorDto`
+ `BuildFullErrorDto(err interface{}, c *gin.Context, others map[string]interface{}, skip int, print bool) *xdto.ErrorDto`
+ `type LoggerExtra struct {}`
+ `WithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context, extra *LoggerExtra)`
+ `WithLogger(logger *log.Logger, start time.Time, c *gin.Context, other string)`
+ `PprofWrap(router *gin.Engine)`
+ `AddBinding(tag string, fn validator.Func) error`
+ `EnableRegexpBinding() error`
+ `EnableRFC3339DateBinding() error`
+ `EnableRFC3339DateTimeBinding() error`
