# xvalidator

### Functions

+ `SetupRegexBinding()`
+ `SetupSpecificRegexpBinding(tag string, re string)`
+ `SetupDateTimeLocBinding(tag string, layout string, loc *time.Location)`
+ `SetupDateTimeBinding(tag string, layout string)`
+ `SetupBinding(tag string, valFunc func(fl validator.FieldLevel) bool)`
+ `IsValidationFormatError(err error) bool`
