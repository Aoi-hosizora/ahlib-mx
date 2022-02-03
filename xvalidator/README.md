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
+ `type WrappedFieldError struct`
+ `type MultiFieldsError struct`
+ `type MessagedStructValidator struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func IsValidationError(err error) bool`
+ `func IsRequiredError(err error) bool`
+ `func UseTagAsFieldName(v *validator.Validate, tagName ...string)`
+ `func ApplyTranslator(validator *validator.Validate, locale LocaleTranslator, registerFn TranslationRegisterHandler) (UtTranslator, error)`
+ `func ApplyEnglishTranslator(validator *validator.Validate) (UtTranslator, error)`
+ `func DefaultRegistrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc`
+ `func DefaultTranslateFunc() validator.TranslationFunc`
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
+ `func TranslateValidationErrors(err validator.ValidationErrors, ut UtTranslator, useNamespace bool) map[string]string`
+ `func FlatValidationErrors(err validator.ValidationErrors, useNamespace bool) map[string]string`
+ `func MapToError(result map[string]string) error`
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

+ `func (w *WrappedFieldError) Origin() validator.FieldError`
+ `func (w *WrappedFieldError) Message() string`
+ `func (w *WrappedFieldError) Error() string`
+ `func (w *WrappedFieldError) Unwrap() error`
+ `func (m *MultiFieldsError) Errors() []error`
+ `func (m *MultiFieldsError) Error() string `
+ `func (m *MultiFieldsError) Translate(translator UtTranslator, useNamespace bool) map[string]string`
+ `func (m *MultiFieldsError) FlatToMap(useNamespace bool) map[string]string`
+ `func (m *MessagedValidator) Engine() interface{}`
+ `func (m *MessagedValidator) ValidateEngine() *validator.Validate`
+ `func (m *MessagedValidator) SetValidateTagName(name string)`
+ `func (m *MessagedValidator) SetMessageTagName(name string)`
+ `func (m *MessagedValidator) UseTagAsFieldName(name ...string)`
+ `func (m *MessagedValidator) ValidateStruct(obj interface{}) error`
