# xfiber

### Functions

+ `DumpRequest(c *fiber.Ctx) []string`
+ `BuildBasicErrorDto(err interface{}, c *fiber.Ctx) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *fiber.Ctx, skip int, print bool) *xdto.ErrorDto`
+ `WithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx)`
+ `WithLogger(logger *log.Logger, start time.Time, c *fiber.Ctx)`
