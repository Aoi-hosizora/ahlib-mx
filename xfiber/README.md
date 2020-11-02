# xfiber

### Functions

+ `DumpRequest(c *fiber.Ctx) []string`
+ `BuildBasicErrorDto(err interface{}, c *fiber.Ctx, otherKvs ...interface{}) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *fiber.Ctx, skip int, doPrint bool, otherKvs ...interface{}) *xdto.ErrorDto`
+ `WithExtraString(s string) LoggerOption`
+ `WithExtraFields(m map[string]interface{}) LoggerOption`
+ `WithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context, options ...LoggerOption)`
+ `WithLogger(logger *log.Logger, start time.Time, c *gin.Context, options ...LoggerOption)`
+ `PprofHandler() func(*fiber.Ctx)`
