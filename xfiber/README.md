# xfiber

### Functions

+ `DumpRequest(c *fiber.Ctx) []string`
+ `BuildBasicErrorDto(err interface{}, c *fiber.Ctx, others map[string]interface{}) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *fiber.Ctx, skip int, print bool) *xdto.ErrorDto`
+ `BuildFullErrorDto(err interface{}, c *fiber.Ctx, others map[string]interface{}, skip int, print bool) *xdto.ErrorDto`
+ `type LoggerExtra struct {}`
+ `WithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context, extra *LoggerExtra)`
+ `WithLogger(logger *log.Logger, start time.Time, c *fiber.Ctx, other string)`
+ `PprofHandler() func(*fiber.Ctx)`
