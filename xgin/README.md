# xgin

### Functions

+ `DumpRequest(c *gin.Context) []string`
+ `BuildBasicErrorDto(err interface{}, c *gin.Context) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *gin.Context, skip int, print bool) *xdto.ErrorDto`
+ `LoggerWithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context)`
