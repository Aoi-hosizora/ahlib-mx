# xgin

### Functions

+ `DumpRequest(c *gin.Context) []string`
+ `BuildBasicErrorDto(err interface{}, c *gin.Context) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *gin.Context, skip int, print bool) *xdto.ErrorDto`
+ `LoggerWithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context)`
+ `AddBinding(tag string, fn func(fl validator.FieldLevel) bool)`
+ `AddBindingString(tag string, fn func(str string) bool)`
+ `AddBindingValue(tag string, fn func(i interface{}) bool)`
+ `EnableRegexpBinding()`
+ `SetupRegexpBinding(tag string, re *regexp.Regexp)`
+ `SetupDateTimeBinding(tag string, layout string)`
+ `SetupLengthBinding(tag string, min, max int)`
+ `SetupOneofBinding(tag string, fields ...interface{})`
+ `SetupOneofValueBinding(tag string, fields ...interface{})`
+ `IsValidationFormatError(err error) bool`
+ `PprofWrap(router *gin.Engine)`
