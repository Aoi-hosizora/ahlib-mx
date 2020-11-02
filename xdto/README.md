# xdto

### Functions

+ `type ErrorDto struct {}`
+ `BuildBasicErrorDto(err interface{}, requests []string, otherKvs ...interface{}) *ErrorDto`
+ `BuildErrorDto(err interface{}, requests []string, skip int, doPrint bool, otherKvs ...interface{}) *ErrorDto`
