# xgin

### Functions

+ `DumpRequest(c *gin.Context) []string`
+ `BuildBasicErrorDto(err interface{}, c *gin.Context) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *gin.Context, skip int, print bool) *xdto.ErrorDto`
+ `LoggerWithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context)`
+ `PprofWrap(router *gin.Engine)`
+ `AddBinding(tag string, fn validator.Func) error`
+ `EnableRegexpBinding() error`
+ `EnableRFC3339DateBinding() error`
+ `EnableRFC3339DateTimeBinding() error`
+ `type HandlerFuncW func(c *gin.Context) (int, interface{})`
+ `JsonW(fn HandlerFuncW) gin.HandlerFunc`
+ `XmlW(fn HandlerFuncW) gin.HandlerFunc`
