# xvalidator

### Functions

+ `ValidationRequiredError(err error) bool`
+ `And(fns ...validator.Func) validator.Func`
+ `Or(fns ...validator.Func) validator.Func`
+ `DefaultRegexpValidator() validator.Func`
+ `RegexpValidator(re *regexp.Regexp) validator.Func`
+ `DateTimeValidator(layout string) validator.Func`
+ `EqualValidator(p interface{}) validator.Func`
+ `NotEqualValidator(p interface{}) validator.Func`
+ `LenValidator(p interface{}) validator.Func`
+ `GreaterThenValidator(p interface{}) validator.Func`
+ `LessThenValidator(p interface{}) validator.Func`
+ `GreaterThenOrEqualValidator(p interface{}) validator.Func`
+ `LessThenOrEqualValidator(p interface{}) validator.Func`
+ `LengthRangeValidator(min, max interface{}) validator.Func`
+ `LengthOutOfRangeValidator(min, max interface{}) validator.Func`
+ `OneofValidator(ps ...interface{}) validator.Func`
+ `GetTranslator(validator *validator.Validate, loc locales.Translator, defaultTranslationFunc func(v *validator.Validate, trans ut.Translator) error) (ut.Translator, error)`
+ `DefaultRegisterTranslationsFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc`
+ `DefaultTranslationFunc() validator.TranslationFunc`
