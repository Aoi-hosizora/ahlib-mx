# xgin

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/gin-gonic/gin (Note: DO NOT update to v1.9.0)
+ github.com/go-playground/validator/v10
+ github.com/sirupsen/logrus
+ golang.org/x/sys

## Documents

### Types

+ `type NewEngineOption func`
+ `type DumpRequestOption func`
+ `type SwaggerOptions struct`
+ `type SwaggerOption func`
+ `type DebugPrintRouteFuncType func`
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

+ `func WithMode(mode string) NewEngineOption`
+ `func WithDebugPrintRouteFunc(debugPrintRouteFunc DebugPrintRouteFuncType) NewEngineOption`
+ `func WithDefaultWriter(defaultWriter io.Writer) NewEngineOption`
+ `func WithDefaultErrorWriter(defaultErrorWriter io.Writer) NewEngineOption`
+ `func WithRedirectTrailingSlash(redirectTrailingSlash bool) NewEngineOption`
+ `func WithRedirectFixedPath(redirectFixedPath bool) NewEngineOption`
+ `func WithHandleMethodNotAllowed(handleMethodNotAllowed bool) NewEngineOption`
+ `func WithForwardedByClientIP(forwardedByClientIP bool) NewEngineOption`
+ `func WithUseRawPath(useRawPath bool) NewEngineOption`
+ `func WithUnescapePathValues(unescapePathValues bool) NewEngineOption`
+ `func WithRemoveExtraSlash(removeExtraSlash bool) NewEngineOption`
+ `func WithRemoteIPHeaders(remoteIPHeaders []string) NewEngineOption`
+ `func WithTrustedPlatform(trustedPlatform string) NewEngineOption`
+ `func WithMaxMultipartMemory(maxMultipartMemory int64) NewEngineOption`
+ `func WithUseH2C(useH2C bool) NewEngineOption`
+ `func WithContextWithFallback(contextWithFallback bool) NewEngineOption`
+ `func WithSecureJSONPrefix(secureJSONPrefix string) NewEngineOption`
+ `func WithNoRoute(noRoute gin.HandlersChain) NewEngineOption`
+ `func WithNoMethod(noMethod gin.HandlersChain) NewEngineOption`
+ `func WithTrustedProxies(trustedProxies []string) NewEngineOption`
+ `func NewEngine(options ...NewEngineOption) *gin.Engine`
+ `func NewEngineSilently(options ...NewEngineOption) *gin.Engine`
+ `func WithIgnoreRequestLine(ignore bool) DumpRequestOption`
+ `func WithRetainHeaders(headers ...string) DumpRequestOption`
+ `func WithIgnoreHeaders(headers ...string) DumpRequestOption`
+ `func WithSecretHeaders(headers ...string) DumpRequestOption`
+ `func WithSecretPlaceholder(placeholder string) DumpRequestOption`
+ `func DumpRequest(c *gin.Context, options ...DumpRequestOption) []string`
+ `func DumpHttpRequest(req *http.Request, options ...DumpRequestOption) []string`
+ `func RedirectHandler(code int, location string) gin.HandlerFunc`
+ `func WrapPprof(router gin.IRouter)`
+ `func WrapPprofSilently(router gin.IRouter)`
+ `func WithSwaggerIndexHtmlRouteName(indexHtmlRouteName string) SwaggerOption`
+ `func WithSwaggerDocJsonRouteName(docJsonRouteName string) SwaggerOption`
+ `func WithSwaggerConfigJsonRouteName(configJsonRouteName string) SwaggerOption`
+ `func WithSwaggerEnableRedirect(enableRedirect bool) SwaggerOption`
+ `func WithSwaggerDeepLinking(deepLinking bool) SwaggerOption`
+ `func WithSwaggerDisplayOperationId(displayOperationId bool) SwaggerOption`
+ `func WithSwaggerDefaultModelsExpandDepth(defaultModelsExpandDepth int) SwaggerOption`
+ `func WithSwaggerDefaultModelExpandDepth(defaultModelExpandDepth int) SwaggerOption`
+ `func WithSwaggerDefaultModelRendering(defaultModelRendering string) SwaggerOption`
+ `func WithSwaggerDisplayRequestDuration(displayRequestDuration bool) SwaggerOption`
+ `func WithSwaggerDocExpansion(docExpansion string) SwaggerOption`
+ `func WithSwaggerMaxDisplayedTags(maxDisplayedTags int) SwaggerOption`
+ `func WithSwaggerOperationsSorter(operationsSorter string) SwaggerOption`
+ `func WithSwaggerShowExtensions(showExtensions bool) SwaggerOption`
+ `func WithSwaggerShowCommonExtensions(showCommonExtensions bool) SwaggerOption`
+ `func WithSwaggerTagsSorter(tagsSorter string) SwaggerOption`
+ `func ReadSwaggerDoc() []byte`
+ `func WrapSwagger(router gin.IRouter, swaggerDocGetter func() []byte, swaggerOptions ...SwaggerOption)`
+ `func WrapSwaggerSilently(router gin.IRouter, swaggerDocGetter func() []byte, swaggerOptions ...SwaggerOption)`
+ `func GetTrustedProxies(engine *gin.Engine) []string`
+ `func HideDebugLogging() (restoreFn func())`
+ `func HideDebugPrintRoute() (restoreFn func())`
+ `func SilentPrintRouteFunc(_, _, _ string, _ int)`
+ `func DefaultPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int)`
+ `func DefaultColorizedPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int)`
+ `func NewRouterDecodeError(field string, input string, err error, message string) *RouterDecodeError`
+ `func ListenAndServeWithReuse(ctx context.Context, server *http.Server) error`
+ `func ReuseListenControl(_, _ string, c syscall.RawConn) error`
+ `func GetValidatorEngine() (*validator.Validate, error)`
+ `func GetValidatorTranslator(locale xvalidator.LocaleTranslator, registerFn xvalidator.TranslationRegisterHandler) (xvalidator.UtTranslator, error)`
+ `func GetValidatorEnglishTranslator() (xvalidator.UtTranslator, error)`
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
+ `func WithMoreExtraText(text string) LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) LoggerOption`
+ `func WithMoreExtraFields(fields map[string]interface{}) LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) LoggerOption`
+ `func WithMoreExtraFieldsV(fields ...interface{}) LoggerOption`
+ `func LogResponseToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...LoggerOption)`
+ `func LogResponseToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...LoggerOption)`
+ `func LogRecoveryToLogrus(logger *logrus.Logger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption)`
+ `func LogRecoveryToLogger(logger logrus.StdLogger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption)`

### Methods

+ `func (r *RouterDecodeError) Error() string`
+ `func (r *RouterDecodeError) Unwrap() error`
