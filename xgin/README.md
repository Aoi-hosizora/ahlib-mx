# xgin

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/gin-gonic/gin
+ github.com/go-playground/validator/v10
+ github.com/sirupsen/logrus

## Documents

### Types

+ `type DumpRequestOption func`
+ `type RouterDecodeError struct`
+ `type TranslatableError interface`
+ `type TranslateOption func`
+ `type LoggerOption func`
+ `type ResponseLoggerParam struct`
+ `type RecoveryLoggerParam struct`

### Variables

+ `var FormatResponseFunc func(p *ResponseLoggerParam) string`
+ `var FieldifyResponseFunc func(p *ResponseLoggerParam) logrus.Fields`
+ `var FormatRecoveryFunc func(p *RecoveryLoggerParam) string`
+ `var FieldifyRecoveryFunc func(p *RecoveryLoggerParam) logrus.Fields`

### Constants

+ None

### Functions

+ `func WithIgnoreRequestLine(ignore bool) DumpRequestOption`
+ `func WithRetainHeaders(headers ...string) DumpRequestOption`
+ `func WithIgnoreHeaders(headers ...string) DumpRequestOption`
+ `func WithSecretHeaders(headers ...string) DumpRequestOption`
+ `func WithSecretPlaceholder(placeholder string) DumpRequestOption`
+ `func DumpRequest(c *gin.Context, options ...DumpRequestOption) []string`
+ `func DumpHttpRequest(req *http.Request, options ...DumpRequestOption) []string`
+ `func WrapPprof(engine *gin.Engine)`
+ `func GetValidatorEngine() (*validator.Validate, error)`
+ `func GetValidatorTranslator(locale xvalidator.LocaleTranslator, registerFn xvalidator.TranslationRegisterHandler) (xvalidator.UtTranslator, error)`
+ `func SetGlobalTranslator(translator xvalidator.UtTranslator)`
+ `func GetGlobalTranslator() xvalidator.UtTranslator`
+ `func AddBinding(tag string, fn validator.Func) error`
+ `func AddTranslation(translator xvalidator.UtTranslator, tag, message string, override bool) error`
+ `func EnableRegexpBinding() error`
+ `func EnableRegexpBindingTranslator(translator ut.Translator) error`
+ `func EnableRFC3339DateBinding() error`
+ `func EnableRFC3339DateBindingTranslator(translator ut.Translator) error`
+ `func EnableRFC3339DateTimeBinding() error`
+ `func EnableRFC3339DateTimeBindingTranslator(translator ut.Translator) error`
+ `func HideDebugPrintRoute() (restoreFn func())`
+ `func NewRouterDecodeError(routerField string, input string, err error, message string) *RouterDecodeError`
+ `func WithUtTranslator(translator xvalidator.UtTranslator) TranslateOption`
+ `func WithJsonInvalidUnmarshalError(fn func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithJsonUnmarshalTypeError(fn func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithJsonSyntaxError(fn func(*json.SyntaxError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithIoEOFError(fn func(error) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithStrconvNumErrorError(fn func(*strconv.NumError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithXginRouterDecodeError(fn func(*RouterDecodeError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithValidatorInvalidTypeError(fn func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithValidatorFieldsError(fn func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithXvalidatorValidateFieldsError(fn func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithTranslatableError(fn func(TranslatableError) (result map[string]string, need4xx bool)) TranslateOption`
+ `func WithExtraErrorsTranslate(fn func(error) (result map[string]string, need4xx bool)) TranslateOption`
+ `func TranslateBindingError(err error, options ...TranslateOption) (result map[string]string, need4xx bool)`
+ `func WithExtraText(text string) LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) LoggerOption`
+ `func LogResponseToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...LoggerOption)`
+ `func LogResponseToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...LoggerOption)`
+ `func LogRecoveryToLogrus(logger *logrus.Logger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption)`
+ `func LogRecoveryToLogger(logger logrus.StdLogger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption)`

### Methods

+ `func (r *RouterDecodeError) Error() string`
+ `func (r *RouterDecodeError) Unwrap() error`
