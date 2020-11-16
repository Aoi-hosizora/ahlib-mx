# xgin

### Functions

+ `DumpRequest(c *gin.Context) []string`
+ `BuildBasicErrorDto(err interface{}, c *gin.Context, otherKvs ...interface{}) *xdto.ErrorDto`
+ `BuildErrorDto(err interface{}, c *gin.Context, skip int, doPrint bool, otherKvs ...interface{}) *xdto.ErrorDto`
+ `WithExtraString(s string) LoggerOption`
+ `WithExtraFields(m map[string]interface{}) LoggerOption`
+ `WithExtraFieldsV(m ...interface{}) LoggerOption`
+ `WithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context, options ...LoggerOption)`
+ `WithLogger(logger *log.Logger, start time.Time, c *gin.Context, options ...LoggerOption)`
+ `PprofWrap(router *gin.Engine)`
+ `func GetValidate() (*validator.Validate, error)`
+ `GetTranslator(loc locales.Translator, translationFunc xvalidator.DefaultTranslationFunc) (ut.Translator, error)`
+ `AddBinding(tag string, fn validator.Func) error`
+ `AddTranslator(translator ut.Translator, tag, message string, override, withParam bool) error`
+ `EnableRegexpBinding() error`
+ `EnableRegexpBindingTranslator(translator ut.Translator) error`
+ `EnableRegexpBindingWithTranslator(translator ut.Translator) error`
+ `EnableRFC3339DateBinding() error`
+ `EnableRFC3339DateBindingTranslator(translator ut.Translator) error`
+ `EnableRFC3339DateBindingWithTranslator(translator ut.Translator) error`
+ `EnableRFC3339DateTimeBinding() error`
+ `EnableRFC3339DateTimeBindingTranslator(translator ut.Translator) error`
+ `EnableRFC3339DateTimeBindingWithTranslator(translator ut.Translator) error`
+ `type AppRoute struct {}`
+ `NewAppRoute(engine *gin.Engine, router gin.IRouter) *AppRoute`
+ `(a *AppRoute) GET(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) POST(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) DELETE(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) PATCH(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) PUT(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) OPTIONS(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) HEAD(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) Any(relativePath string, handlers ...gin.HandlerFunc)`
+ `(a *AppRoute) Do()`
