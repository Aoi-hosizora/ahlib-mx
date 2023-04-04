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
+ `type MessagedValidatorOption func`
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
+ `func EnLocaleTranslator() LocaleTranslator`
+ `func EsLocaleTranslator() LocaleTranslator`
+ `func FrLocaleTranslator() LocaleTranslator`
+ `func IdLocaleTranslator() LocaleTranslator`
+ `func ItLocaleTranslator() LocaleTranslator`
+ `func JaLocaleTranslator() LocaleTranslator`
+ `func NlLocaleTranslator() LocaleTranslator`
+ `func PtLocaleTranslator() LocaleTranslator`
+ `func PtBrLocaleTranslator() LocaleTranslator`
+ `func RuLocaleTranslator() LocaleTranslator`
+ `func TrLocaleTranslator() LocaleTranslator`
+ `func ViLocaleTranslator() LocaleTranslator`
+ `func ZhLocaleTranslator() LocaleTranslator`
+ `func ZhHantLocaleTranslator() LocaleTranslator`
+ `func EnTranslationRegisterFunc() TranslationRegisterHandler`
+ `func EsTranslationRegisterFunc() TranslationRegisterHandler`
+ `func FrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func IdTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ItTranslationRegisterFunc() TranslationRegisterHandler`
+ `func JaTranslationRegisterFunc() TranslationRegisterHandler`
+ `func NlTranslationRegisterFunc() TranslationRegisterHandler`
+ `func PtTranslationRegisterFunc() TranslationRegisterHandler`
+ `func PtBrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func RuTranslationRegisterFunc() TranslationRegisterHandler`
+ `func TrTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ViTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ZhTranslationRegisterFunc() TranslationRegisterHandler`
+ `func ZhHantTranslationRegisterFunc() TranslationRegisterHandler`
+ `func TranslateValidationErrors(err validator.ValidationErrors, ut UtTranslator, useNamespace bool) map[string]string`
+ `func FlattenValidationErrors(err validator.ValidationErrors, useNamespace bool) map[string]string`
+ `func MergeMapToError(result map[string]string) error`
+ `func WithValidateTagName(name string) MessagedValidatorOption`
+ `func WithValidateMessageTagName(name string) MessagedValidatorOption`
+ `func NewMessagedValidator(options ...MessagedValidatorOption) *MessagedValidator`
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
+ `func (m *MultiFieldsError) Error() string`
+ `func (m *MultiFieldsError) Is(target error) bool`
+ `func (m *MultiFieldsError) As(target interface{}) bool`
+ `func (m *MultiFieldsError) Translate(translator UtTranslator, useNamespace bool) map[string]string`
+ `func (m *MultiFieldsError) Flatten(useNamespace bool) map[string]string`
+ `func (m *MessagedValidator) Engine() interface{}`
+ `func (m *MessagedValidator) ValidateEngine() *validator.Validate`
+ `func (m *MessagedValidator) SetValidateTagName(name string)`
+ `func (m *MessagedValidator) SetMessageTagName(name string)`
+ `func (m *MessagedValidator) UseTagAsFieldName(name ...string)`
+ `func (m *MessagedValidator) ValidateStruct(obj interface{}) error`
