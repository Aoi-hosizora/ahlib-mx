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
+ `type TranslateOption func`

### Variables

+ None

### Constants

+ None

### Functions

+ `func WithIgnoreRequestLine(ignore bool) DumpRequestOption`
+ `func WithRetainHeaders(headers ...string) DumpRequestOption`
+ `func WithIgnoreHeaders(headers ...string) DumpRequestOption`
+ `func WithSecretHeaders(headers ...string) DumpRequestOption`
+ `func WithSecretReplace(secret string) DumpRequestOption`
+ `func DumpRequest(c *gin.Context, options ...DumpRequestOption) []string`
+ `func DumpHttpRequest(req *http.Request, options ...DumpRequestOption) []string`
+ `func PprofWrap(router *gin.Engine)`
+ `func GetValidatorEngine() (*validator.Validate, error)`
+ `func GetValidatorTranslator(locale xvalidator.LocaleTranslator, registerFn xvalidator.TranslationRegisterHandler) (xvalidator.UtTranslator, error)`
+ `func AddBinding(tag string, fn validator.Func) error`
+ `func AddTranslation(translator xvalidator.UtTranslator, tag, message string, override bool) error`
+ `func EnableRegexpBinding() error`
+ `func EnableRegexpBindingTranslator(translator ut.Translator) error`
+ `func EnableRFC3339DateBinding() error`
+ `func EnableRFC3339DateBindingTranslator(translator ut.Translator) error`
+ `func EnableRFC3339DateTimeBinding() error`
+ `func EnableRFC3339DateTimeBindingTranslator(translator ut.Translator) error`
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
+ `func WithExtraErrorsTranslate(fn func(error) (result map[string]string, need4xx bool)) TranslateOption`
+ `func TranslateBindingError(err error, options ...TranslateOption) (result map[string]string, need4xx bool)`
+ `func WithExtraText(text string) logop.LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) logop.LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption`
+ `func LogToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption)`
+ `func LogToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption)`

### Methods

+ `func (r *RouterDecodeError) Error() string`
+ `func (r *RouterDecodeError) Unwrap() error`
