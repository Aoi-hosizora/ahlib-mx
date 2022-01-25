# xvalidator

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/go-playground/validator/v10
+ github.com/go-playground/universal-translator
+ github.com/go-playground/locales

## Documents

### Types

+ `type UtTranslator = ut.Translator`
+ `type LocaleTranslator = locales.Translator`
+ `type TranslationRegisterHandler func`
+ `type WrappedValidateFieldError struct`
+ `type ValidateFieldsError struct`
+ `type CustomStructValidator struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func IsValidationError(err error) bool`
+ `func IsRequiredError(err error) bool`
+ `func UseTagAsFieldName(v *validator.Validate, tagName string)`
+ `func UseDefaultFieldName(v *validator.Validate)`
+ `func ApplyTranslator(validator *validator.Validate, locale LocaleTranslator, registerFn TranslationRegisterHandler) (UtTranslator, error)`
+ `func DefaultRegistrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc`
+ `func DefaultTranslateFunc() validator.TranslationFunc`
+ `func TranslateValidationErrors(err validator.ValidationErrors, ut UtTranslator, useNamespace bool) map[string]string`
+ `func FlatValidateErrors(err validator.ValidationErrors, useNamespace bool) map[string]string`
+ `func FlattedMapToError(kv map[string]string) error`
+ `func EnLocaleTranslator() locales.Translator`
+ `func EsLocaleTranslator() locales.Translator`
+ `func FrLocaleTranslator() locales.Translator`
+ `func IdLocaleTranslator() locales.Translator`
+ `func JaLocaleTranslator() locales.Translator`
+ `func NlLocaleTranslator() locales.Translator`
+ `func PtLocaleTranslator() locales.Translator`
+ `func PtBrLocaleTranslator() locales.Translator`
+ `func RuLocaleTranslator() locales.Translator`
+ `func TrLocaleTranslator() locales.Translator`
+ `func ZhLocaleTranslator() locales.Translator`
+ `func ZhHantLocaleTranslator() locales.Translator`
+ `func EnTranslationRegisterFunc() TranslationRegisterHandler`
+ `func EsTranslationRegisterFunc() TranslationRegisterHandler`
+ `func FrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func IdTranslationRegisterFunc() TranslationRegisterHandler`
+ `func JaTranslationRegisterFunc() TranslationRegisterHandler`
+ `func NlTranslationRegisterFunc() TranslationRegisterHandler`
+ `func PtTranslationRegisterFunc() TranslationRegisterHandler`
+ `func PtBrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func RuTranslationRegisterFunc() TranslationRegisterHandler`
+ `func TrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ZhTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ZhHantTranslationRegisterFunc() TranslationRegisterHandler`
+ `func NewCustomStructValidator() *CustomStructValidator`
+ `func ParamRegexpValidator() validator.Func`
+ `func RegexpValidator(re *regexp.Regexp) validator.Func`
+ `func DateTimeValidator(layout string) validator.Func`
+ `func And(fns ...validator.Func) validator.Func`
+ `func Or(fns ...validator.Func) validator.Func`
+ `func Not(fn validator.Func) validator.Func`
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

### Methods

+ `func (v *WrappedValidateFieldError) Origin() validator.FieldError`
+ `func (v *WrappedValidateFieldError) Message() string`
+ `func (v *WrappedValidateFieldError) Error() string`
+ `func (v *ValidateFieldsError) Fields() []error`
+ `func (v *ValidateFieldsError) Error() string `
+ `func (v *ValidateFieldsError) Translate(translator UtTranslator, useNamespace bool) map[string]string`
+ `func (v *ValidateFieldsError) FlatToMap(useNamespace bool) map[string]string`
+ `func (v *CustomStructValidator) Engine() interface{}`
+ `func (v *CustomStructValidator) ValidateEngine() *validator.Validate`
+ `func (v *CustomStructValidator) SetValidatorTagName(name string)`
+ `func (v *CustomStructValidator) SetMessageTagName(name string)`
+ `func (v *CustomStructValidator) SetFieldNameTag(name string)`
+ `func (v *CustomStructValidator) ValidateStruct(obj interface{}) error`
