# xvalidator

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/go-playground/locales
+ github.com/go-playground/universal-translator
+ github.com/go-playground/validator/v10

## Documents

### Types

+ `type UtTranslator = ut.Translator`
+ `type TranslationRegisterHandler func`

### Variables

+ None

### Constants

+ None

### Functions

+ `func IsValidationError(err error) bool`
+ `func IsRequiredError(err error) bool`
+ `func ParamRegexpValidator() validator.Func`
+ `func RegexpValidator(re *regexp.Regexp) validator.Func`
+ `func DateTimeValidator(layout string) validator.Func`
+ `func And(fns ...validator.Func) validator.Func`
+ `func Or(fns ...validator.Func) validator.Func`
+ `func EqualValidator(p interface{}) validator.Func`
+ `func NotEqualValidator(p interface{}) validator.Func`
+ `func LenValidator(p interface{}) validator.Func`
+ `func GreaterThenValidator(p interface{}) validator.Func`
+ `func LessThenValidator(p interface{}) validator.Func`
+ `func GreaterThenOrEqualValidator(p interface{}) validator.Func`
+ `func LessThenOrEqualValidator(p interface{}) validator.Func`
+ `func LengthInRangeValidator(min, max interface{}) validator.Func`
+ `func LengthOutOfRangeValidator(min, max interface{}) validator.Func`
+ `func OneofValidator(ps ...interface{}) validator.Func`
+ `func ApplyTranslator(validator *validator.Validate, locTranslator locales.Translator, registerFn TranslationRegisterHandler) (ut.Translator, error)`
+ `func AddToTranslatorFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc`
+ `func DefaultTranslateFunc() validator.TranslationFunc`
+ `func EnLocaleTranslator() locales.Translator`
+ `func EnLocaleTranslator() locales.Translator`
+ `func FrLocaleTranslator() locales.Translator`
+ `func IdLocaleTranslator() locales.Translator`
+ `func JaLocaleTranslator() locales.Translator`
+ `func NlLocaleTranslator() locales.Translator`
+ `func PtBrLocaleTranslator() locales.Translator`
+ `func RuLocaleTranslator() locales.Translator`
+ `func TrLocaleTranslator() locales.Translator`
+ `func ZhLocaleTranslator() locales.Translator`
+ `func ZhHantLocaleTranslator() locales.Translator`
+ `func EnTranslationRegisterFunc() TranslationRegisterHandler`
+ `func FrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func IdTranslationRegisterFunc() TranslationRegisterHandler`
+ `func JaTranslationRegisterFunc() TranslationRegisterHandler`
+ `func NlTranslationRegisterFunc() TranslationRegisterHandler`
+ `func PtBrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func RuTranslationRegisterFunc() TranslationRegisterHandler`
+ `func TrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ZhTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ZhTwTranslationRegisterFunc() TranslationRegisterHandler`

### Methods

+ None