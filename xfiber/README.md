# xfiber

### Functions

+ `DumpRequest(c *fiber.Ctx) []string`
+ `BuildBasicErrorDto(err interface{}, c *fiber.Ctx, others map[string]interface{}) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *fiber.Ctx, others map[string]interface{}, skip int, print bool) *xdto.ErrorDto`
+ `WithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx, other string, otherFields map[string]interface{})`
+ `WithLogger(logger *log.Logger, start time.Time, c *fiber.Ctx, other string)`
+ `PprofHandler() func(*fiber.Ctx)`
+ `AddBinding(tag string, fn validator.Func) error`
+ `EnableRegexpBinding() error`
+ `EnableRFC3339DateBinding() error`
+ `EnableRFC3339DateTimeBinding() error`
+ `Struct(s interface{}) error`
+ `Var(field interface{}, tag string) error`
