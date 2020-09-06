# xdto

### Functions

+ `type ErrorDto struct {}`
+ `BuildBasicErrorDto(err interface{}, requests []string, others map[string]interface{}) *ErrorDto`
+ `BuildErrorDto(err interface{}, requests []string, others map[string]interface{}, skip int, print bool) *ErrorDto`
