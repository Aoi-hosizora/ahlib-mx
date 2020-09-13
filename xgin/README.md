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
+ `type AppRoute struct {}`
+ `NewAppRoute(engine *gin.Engine, router gin.IRouter) *AppRoute`
+ `(a *AppRoute) GET(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) POST(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) DELETE(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) PATCH(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) PUT(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) OPTIONS(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) HEAD(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) Any(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) Do()`
