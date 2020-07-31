# xgin

### Functions

+ `DumpRequest(c *gin.Context) []string`
+ `BuildBasicErrorDto(err interface{}, c *gin.Context) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *gin.Context, skip int, print bool) *xdto.ErrorDto`
+ `LoggerWithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context)`

### Validator Functions

+ `SetupRegexBinding()`
+ `SetupSpecificRegexpBinding(tag string, re string)`
+ `SetupDateTimeLocBinding(tag string, layout string, loc *time.Location)`
+ `SetupDateTimeBinding(tag string, layout string)`
+ `SetupBinding(tag string, valFunc func(fl validator.FieldLevel) bool)`
+ `IsValidationFormatError(err error) bool`
