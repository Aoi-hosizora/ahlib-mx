# xfiber

### Functions

+ `DumpRequest(c *fiber.Ctx) []string`
+ `BuildBasicErrorDto(err interface{}, c *fiber.Ctx) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *fiber.Ctx, skip int, print bool) *xdto.ErrorDto`
+ `LoggerWithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx)`
