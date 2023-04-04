package xgin

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xpointer"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"strings"
)

// ==========
// new engine
// ==========

var (
	// _defaultEngine represents a global default gin.Engine.
	_defaultEngine *gin.Engine
)

func init() {
	// gin.DebugMode will have no influence to the server's request handling, it
	// only affects: debug logging, template rendering, internal recovery handler.
	gin.SetMode(gin.DebugMode)

	// use xgin.DefaultPrintRouteFunc as default DebugPrintRouteFunc, which has
	// the same effect with gin's default print route func.
	gin.DebugPrintRouteFunc = DefaultPrintRouteFunc

	// set global engine, and hide any logging
	fn := HideDebugLogging()
	_defaultEngine = gin.New()
	fn()
}

// newEngineOptions is a type of NewEngine's option, each field can be set by NewEngineOption function type.
type newEngineOptions struct {
	mode                string
	debugPrintRouteFunc DebugPrintRouteFuncType
	defaultWriter       io.Writer
	defaultErrorWriter  io.Writer

	redirectTrailingSlash  bool
	redirectFixedPath      bool
	handleMethodNotAllowed bool
	forwardedByClientIP    bool
	useRawPath             bool
	unescapePathValues     bool
	removeExtraSlash       bool
	remoteIPHeaders        []string
	trustedPlatform        string
	maxMultipartMemory     int64
	useH2C                 bool
	contextWithFallback    bool

	secureJSONPrefix string
	noRoute          gin.HandlersChain
	noMethod         gin.HandlersChain
	trustedProxies   []string
}

// NewEngineOption represents an option type for NewEngine's option, can be created by WithXXX functions.
type NewEngineOption func(o *newEngineOptions)

// WithMode creates a NewEngineOption to specify gin's global mode, defaults to gin.DebugMode.
func WithMode(mode string) NewEngineOption {
	return func(o *newEngineOptions) {
		o.mode = mode
	}
}

// WithDebugPrintRouteFunc creates a NewEngineOption to specify gin's global debugPrintRouteFunc, defaults to DefaultPrintRouteFunc.
func WithDebugPrintRouteFunc(debugPrintRouteFunc DebugPrintRouteFuncType) NewEngineOption {
	return func(o *newEngineOptions) {
		o.debugPrintRouteFunc = debugPrintRouteFunc
	}
}

// WithDefaultWriter creates a NewEngineOption to specify gin's global defaultWriter, defaults to os.Stdout.
func WithDefaultWriter(defaultWriter io.Writer) NewEngineOption {
	return func(o *newEngineOptions) {
		o.defaultWriter = defaultWriter
	}
}

// WithDefaultErrorWriter creates a NewEngineOption to specify gin's global defaultErrorWriter, defaults to os.Stderr.
func WithDefaultErrorWriter(defaultErrorWriter io.Writer) NewEngineOption {
	return func(o *newEngineOptions) {
		o.defaultErrorWriter = defaultErrorWriter
	}
}

// WithRedirectTrailingSlash creates a NewEngineOption to specify gin engine's redirectTrailingSlash, defaults to true.
func WithRedirectTrailingSlash(redirectTrailingSlash bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.redirectTrailingSlash = redirectTrailingSlash
	}
}

// WithRedirectFixedPath creates a NewEngineOption to specify gin engine's redirectFixedPath, defaults to false.
func WithRedirectFixedPath(redirectFixedPath bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.redirectFixedPath = redirectFixedPath
	}
}

// WithHandleMethodNotAllowed creates a NewEngineOption to specify gin engine's handleMethodNotAllowed, defaults to false.
func WithHandleMethodNotAllowed(handleMethodNotAllowed bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.handleMethodNotAllowed = handleMethodNotAllowed
	}
}

// WithForwardedByClientIP creates a NewEngineOption to specify gin engine's forwardedByClientIP, defaults to true.
func WithForwardedByClientIP(forwardedByClientIP bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.forwardedByClientIP = forwardedByClientIP
	}
}

// WithUseRawPath creates a NewEngineOption to specify gin engine's useRawPath, defaults to false.
func WithUseRawPath(useRawPath bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.useRawPath = useRawPath
	}
}

// WithUnescapePathValues creates a NewEngineOption to specify gin engine's unescapePathValues, defaults to true.
func WithUnescapePathValues(unescapePathValues bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.unescapePathValues = unescapePathValues
	}
}

// WithRemoveExtraSlash creates a NewEngineOption to specify gin engine's removeExtraSlash, defaults to false.
func WithRemoveExtraSlash(removeExtraSlash bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.removeExtraSlash = removeExtraSlash
	}
}

// WithRemoteIPHeaders creates a NewEngineOption to specify gin engine's remoteIPHeaders, defaults to ["X-Forwarded-For", "X-Real-IP"].
func WithRemoteIPHeaders(remoteIPHeaders []string) NewEngineOption {
	return func(o *newEngineOptions) {
		o.remoteIPHeaders = remoteIPHeaders
	}
}

// WithTrustedPlatform creates a NewEngineOption to specify gin engine's trustedPlatform, defaults to "".
func WithTrustedPlatform(trustedPlatform string) NewEngineOption {
	return func(o *newEngineOptions) {
		o.trustedPlatform = trustedPlatform
	}
}

// WithMaxMultipartMemory creates a NewEngineOption to specify gin engine's maxMultipartMemory, defaults to 32 MB.
func WithMaxMultipartMemory(maxMultipartMemory int64) NewEngineOption {
	return func(o *newEngineOptions) {
		o.maxMultipartMemory = maxMultipartMemory
	}
}

// WithUseH2C creates a NewEngineOption to specify gin engine's useH2C, defaults to false.
func WithUseH2C(useH2C bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.useH2C = useH2C
	}
}

// WithContextWithFallback creates a NewEngineOption to specify gin engine's contextWithFallback, defaults to false.
func WithContextWithFallback(contextWithFallback bool) NewEngineOption {
	return func(o *newEngineOptions) {
		o.contextWithFallback = contextWithFallback
	}
}

// WithSecureJSONPrefix creates a NewEngineOption to specify gin engine's secureJSONPrefix, defaults to "while(1);".
func WithSecureJSONPrefix(secureJSONPrefix string) NewEngineOption {
	return func(o *newEngineOptions) {
		o.secureJSONPrefix = secureJSONPrefix
	}
}

// WithNoRoute creates a NewEngineOption to specify gin engine's noRoute, defaults to nil.
func WithNoRoute(noRoute gin.HandlersChain) NewEngineOption {
	return func(o *newEngineOptions) {
		o.noRoute = noRoute
	}
}

// WithNoMethod creates a NewEngineOption to specify gin engine's noMethod, defaults to nil.
func WithNoMethod(noMethod gin.HandlersChain) NewEngineOption {
	return func(o *newEngineOptions) {
		o.noMethod = noMethod
	}
}

// WithTrustedProxies creates a NewEngineOption to specify gin engine's trustedProxies, defaults to ["0.0.0.0/0", "::/0"].
func WithTrustedProxies(trustedProxies []string) NewEngineOption {
	return func(o *newEngineOptions) {
		o.trustedProxies = trustedProxies
	}
}

// NewEngine creates a new blank gin.Engine instance with some default settings. Note that WithMode, WithDebugPrintRouteFunc, WithDefaultWriter and WithDefaultErrorWriter options
// will not change gin's global setting if these options are not given, and will directly update global setting if any valid option is passed.
func NewEngine(options ...NewEngineOption) *gin.Engine {
	opt := &newEngineOptions{
		mode:                gin.Mode(),              // gin.DebugMode
		debugPrintRouteFunc: gin.DebugPrintRouteFunc, // DefaultPrintRouteFunc
		defaultWriter:       gin.DefaultWriter,       // os.Stdout
		defaultErrorWriter:  gin.DefaultErrorWriter,  // os.Stderr

		redirectTrailingSlash:  _defaultEngine.RedirectTrailingSlash,  // true
		redirectFixedPath:      _defaultEngine.RedirectFixedPath,      // false
		handleMethodNotAllowed: _defaultEngine.HandleMethodNotAllowed, // false
		forwardedByClientIP:    _defaultEngine.ForwardedByClientIP,    // true
		useRawPath:             _defaultEngine.UseRawPath,             // false
		unescapePathValues:     _defaultEngine.UnescapePathValues,     // true
		removeExtraSlash:       _defaultEngine.RemoveExtraSlash,       // false
		remoteIPHeaders:        _defaultEngine.RemoteIPHeaders,        // ["X-Forwarded-For", "X-Real-IP"]
		trustedPlatform:        _defaultEngine.TrustedPlatform,        // ""
		maxMultipartMemory:     _defaultEngine.MaxMultipartMemory,     // 32 MB
		useH2C:                 _defaultEngine.UseH2C,                 // false
		contextWithFallback:    _defaultEngine.ContextWithFallback,    // false

		secureJSONPrefix: xreflect.GetUnexportedField(xreflect.FieldValueOf(_defaultEngine, "secureJSONPrefix")).Interface().(string),    // "while(1);"
		noRoute:          xreflect.GetUnexportedField(xreflect.FieldValueOf(_defaultEngine, "noRoute")).Interface().(gin.HandlersChain),  // nil
		noMethod:         xreflect.GetUnexportedField(xreflect.FieldValueOf(_defaultEngine, "noMethod")).Interface().(gin.HandlersChain), // nil
		trustedProxies:   xreflect.GetUnexportedField(xreflect.FieldValueOf(_defaultEngine, "trustedProxies")).Interface().([]string),    // ["0.0.0.0/0", "::/0"]
	}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}

	// global setting
	if opt.mode == gin.DebugMode || opt.mode == gin.ReleaseMode || opt.mode == gin.TestMode {
		gin.SetMode(opt.mode)
	}
	if opt.debugPrintRouteFunc != nil {
		gin.DebugPrintRouteFunc = opt.debugPrintRouteFunc
	}
	if opt.defaultWriter != nil {
		gin.DefaultWriter = opt.defaultWriter
	}
	if opt.defaultErrorWriter != nil {
		gin.DefaultErrorWriter = opt.defaultErrorWriter
	}

	// create engine after updating global settings
	engine := gin.New()

	// instance setting
	engine.RedirectTrailingSlash = opt.redirectTrailingSlash
	engine.RedirectFixedPath = opt.redirectFixedPath
	engine.HandleMethodNotAllowed = opt.handleMethodNotAllowed
	engine.ForwardedByClientIP = opt.forwardedByClientIP
	engine.UseRawPath = opt.useRawPath
	engine.UnescapePathValues = opt.unescapePathValues
	engine.RemoveExtraSlash = opt.removeExtraSlash
	if opt.remoteIPHeaders != nil {
		engine.RemoteIPHeaders = opt.remoteIPHeaders
	}
	engine.TrustedPlatform = opt.trustedPlatform
	engine.MaxMultipartMemory = opt.maxMultipartMemory
	engine.UseH2C = opt.useH2C
	engine.ContextWithFallback = opt.contextWithFallback
	engine.SecureJsonPrefix(opt.secureJSONPrefix)
	if opt.noRoute != nil {
		engine.NoRoute(opt.noRoute...)
	}
	if opt.noMethod != nil {
		engine.NoMethod(opt.noMethod...)
	}
	if opt.trustedProxies != nil {
		_ = engine.SetTrustedProxies(opt.trustedProxies)
	}

	return engine
}

// NewEngineSilently creates a new blank gin.Engine instance with some default settings, without any debug logging.
func NewEngineSilently(options ...NewEngineOption) *gin.Engine {
	restore := HideDebugLogging()
	engine := NewEngine(options...)
	restore()

	// keep writer consistent when using WithDefaultWriter and WithDefaultErrorWriter
	opt := &newEngineOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.defaultWriter != nil {
		gin.DefaultWriter = opt.defaultWriter
	}
	if opt.defaultErrorWriter != nil {
		gin.DefaultErrorWriter = opt.defaultErrorWriter
	}

	return engine
}

// ============
// dump request
// ============

// dumpRequestOptions is a type of DumpHttpRequest's option, each field can be set by DumpRequestOption function type.
type dumpRequestOptions struct {
	ignoreRequestLine bool
	retainHeaders     []string
	ignoreHeaders     []string
	secretHeaders     []string
	secretPlaceholder string
}

// DumpRequestOption represents an option type for DumpHttpRequest's option, can be created by WithXXX functions.
type DumpRequestOption func(*dumpRequestOptions)

// WithIgnoreRequestLine creates a DumpRequestOption for request line, if set to true, request line "GET /xxx HTTP/1.1" will be ignored.
func WithIgnoreRequestLine(ignore bool) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.ignoreRequestLine = ignore
	}
}

// WithRetainHeaders creates a DumpRequestOption for headers which are wanted to retain, WithIgnoreHeaders option will be ignored when use with this.
func WithRetainHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.retainHeaders = headers
	}
}

// WithIgnoreHeaders creates a DumpRequestOption for headers which are wanted to ignore, this option will be ignored when used with WithRetainHeaders.
func WithIgnoreHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.ignoreHeaders = headers
	}
}

// WithSecretHeaders creates a DumpRequestOption for headers which are secret, such as Authorization field, also see WithSecretPlaceholder.
func WithSecretHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.secretHeaders = headers
	}
}

// WithSecretPlaceholder creates a DumpRequestOption to specify a secret placeholder for secret headers set by WithSecretHeaders, defaults to "*".
func WithSecretPlaceholder(placeholder string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.secretPlaceholder = placeholder
	}
}

// isSpecificHeader checks whether given param is the same specific header in case-insensitive.
func isSpecificHeader(param, header string) bool {
	param = strings.ToLower(param)
	header = strings.ToLower(header)
	return strings.HasPrefix(param, header+": ")
}

// DumpRequest dumps and formats http.Request from gin.Context to string slice using given DumpRequestOption-s.
func DumpRequest(c *gin.Context, options ...DumpRequestOption) []string {
	if c == nil {
		return nil
	}
	return DumpHttpRequest(c.Request, options...)
}

// DumpHttpRequest dumps and formats http.Request to string slice using given DumpRequestOption-s.
func DumpHttpRequest(req *http.Request, options ...DumpRequestOption) []string {
	if req == nil {
		return nil
	}
	opt := &dumpRequestOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.secretPlaceholder == "" {
		opt.secretPlaceholder = "*"
	}

	bs, _ := httputil.DumpRequest(req, false)
	// unreachable error, because body is false, and bytes.Buffer's
	// WriteString method will never return error

	lines := strings.Split(xstring.FastBtos(bs), "\r\n") // split by \r\n
	result := make([]string, 0, len(lines))
	for i, line := range lines {
		if i == 0 {
			if !opt.ignoreRequestLine {
				result = append(result, line) // request line: METHOD /ENDPOINT HTTP/1.1
			}
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			// after the request line, there is \r\n\r\n, which is the splitter between the request line and the request header
			continue
		}

		// I. filter headers, use retainHeaders first
		headerList := opt.retainHeaders
		toIgnore := false
		if len(opt.retainHeaders) == 0 {
			headerList = opt.ignoreHeaders
			toIgnore = true
		}
		exist := false
		for _, header := range headerList {
			if isSpecificHeader(line, header) {
				exist = true
				break
			}
		}
		if (!toIgnore && !exist) || (toIgnore && exist) {
			continue
		}

		// II. rewrite headers that are secret
		for _, header := range opt.secretHeaders {
			if isSpecificHeader(line, header) {
				header = strings.SplitN(line, ":", 2)[0]
				line = header + ": " + opt.secretPlaceholder
				break
			}
		}

		// III. append to the result slice
		result = append(result, line)
	}

	return result
}

// ===================
// handler and wrapper
// ===================

// RedirectHandler creates a gin.HandlerFunc that behaviors a redirection with given code (such as http.StatusMovedPermanently or http.StatusTemporaryRedirect)
// and redirect target location.
func RedirectHandler(code int, location string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(code, location)
	}
}

// WrapPprof registers several pprof routes from package `net/http/pprof` to gin router. For more, please visit https://github.com/DeanThompson/ginpprof.
func WrapPprof(router gin.IRouter) {
	for _, r := range []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/debug/pprof/", pprofIndexHandler},
		{"GET", "/debug/pprof/heap", pprofHeapHandler},
		{"GET", "/debug/pprof/goroutine", pprofGoroutineHandler},
		{"GET", "/debug/pprof/allocs", pprofAllocsHandler},
		{"GET", "/debug/pprof/block", pprofBlockHandler},
		{"GET", "/debug/pprof/threadcreate", pprofThreadHandler},
		{"GET", "/debug/pprof/cmdline", pprofCmdlineHandler},
		{"GET", "/debug/pprof/profile", pprofProfileHandler},
		{"GET", "/debug/pprof/symbol", pprofSymbolHandler},
		{"POST", "/debug/pprof/symbol", pprofSymbolHandler},
		{"GET", "/debug/pprof/trace", pprofTraceHandler},
		{"GET", "/debug/pprof/mutex", pprofMutexHandler},
	} {
		router.Handle(r.method, r.path, r.handler) // use path directly
	}
}

// WrapPprofSilently registers several pprof routes from package `net/http/pprof` to gin router without any debug logging.
func WrapPprofSilently(router gin.IRouter) {
	restore := HideDebugPrintRoute()
	WrapPprof(router)
	restore()
}

func pprofIndexHandler(c *gin.Context)     { pprof.Index(c.Writer, c.Request) }                             // GET
func pprofCmdlineHandler(c *gin.Context)   { pprof.Cmdline(c.Writer, c.Request) }                           // GET
func pprofProfileHandler(c *gin.Context)   { pprof.Profile(c.Writer, c.Request) }                           // GET
func pprofSymbolHandler(c *gin.Context)    { pprof.Symbol(c.Writer, c.Request) }                            // GET / POST
func pprofTraceHandler(c *gin.Context)     { pprof.Trace(c.Writer, c.Request) }                             // GET
func pprofAllocsHandler(c *gin.Context)    { pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request) }       // GET
func pprofBlockHandler(c *gin.Context)     { pprof.Handler("block").ServeHTTP(c.Writer, c.Request) }        // GET
func pprofGoroutineHandler(c *gin.Context) { pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request) }    // GET
func pprofHeapHandler(c *gin.Context)      { pprof.Handler("heap").ServeHTTP(c.Writer, c.Request) }         // GET
func pprofMutexHandler(c *gin.Context)     { pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request) }        // GET
func pprofThreadHandler(c *gin.Context)    { pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request) } // GET

const (
	panicWrapSwaggerToEngine = "xgin: wrapping swagger to gin.Engine is not allowed"
	panicNilSwaggerDocGetter = "xgin: using nil swaggerDocGetter function"
)

// SwaggerOptions is a type of WrapSwagger's option, each field can be set by SwaggerOption function type.
type SwaggerOptions struct {
	IndexHtmlRouteName  string `json:"-"`
	DocJsonRouteName    string `json:"-"`
	ConfigJsonRouteName string `json:"-"`
	EnableRedirect      *bool  `json:"-"`

	DeepLinking              *bool   `json:"deepLinking,omitempty"`
	DisplayOperationId       *bool   `json:"displayOperationId,omitempty"`
	DefaultModelsExpandDepth *int    `json:"defaultModelsExpandDepth,omitempty"`
	DefaultModelExpandDepth  *int    `json:"defaultModelExpandDepth,omitempty"`
	DefaultModelRendering    *string `json:"defaultModelRendering,omitempty"`
	DisplayRequestDuration   *bool   `json:"displayRequestDuration,omitempty"`
	DocExpansion             *string `json:"docExpansion,omitempty"`
	MaxDisplayedTags         *int    `json:"maxDisplayedTags,omitempty"`
	OperationsSorter         *string `json:"operationsSorter,omitempty"`
	ShowExtensions           *bool   `json:"showExtensions,omitempty"`
	ShowCommonExtensions     *bool   `json:"showCommonExtensions,omitempty"`
	TagsSorter               *string `json:"tagsSorter,omitempty"`
}

// SwaggerOption represents an option type for WrapSwagger's option, can be created by WithXXX functions.
type SwaggerOption func(o *SwaggerOptions)

// WithSwaggerIndexHtmlRouteName creates a SwaggerOption to specify swagger index.html route name, defaults to "index.html".
func WithSwaggerIndexHtmlRouteName(indexHtmlRouteName string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.IndexHtmlRouteName = indexHtmlRouteName
	}
}

// WithSwaggerDocJsonRouteName creates a SwaggerOption to specify swagger doc.json route name, defaults to "doc.json".
func WithSwaggerDocJsonRouteName(docJsonRouteName string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DocJsonRouteName = docJsonRouteName
	}
}

// WithSwaggerConfigJsonRouteName creates a SwaggerOption to specify swagger config.json route name, defaults to "config.json".
func WithSwaggerConfigJsonRouteName(configJsonRouteName string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.ConfigJsonRouteName = configJsonRouteName
	}
}

// WithSwaggerEnableRedirect creates a SwaggerOption to specify whether to enable redirect route, defaults to true.
func WithSwaggerEnableRedirect(enableRedirect bool) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.EnableRedirect = &enableRedirect
	}
}

// WithSwaggerDeepLinking creates a SwaggerOption to specify swagger deepLinking, defaults to true.
func WithSwaggerDeepLinking(deepLinking bool) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DeepLinking = &deepLinking
	}
}

// WithSwaggerDisplayOperationId creates a SwaggerOption to specify swagger displayOperationId, defaults to false.
func WithSwaggerDisplayOperationId(displayOperationId bool) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DisplayOperationId = &displayOperationId
	}
}

// WithSwaggerDefaultModelsExpandDepth creates a SwaggerOption to specify swagger defaultModelsExpandDepth, set to -1 completely hide the models, defaults to 1.
func WithSwaggerDefaultModelsExpandDepth(defaultModelsExpandDepth int) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DefaultModelsExpandDepth = &defaultModelsExpandDepth
	}
}

// WithSwaggerDefaultModelExpandDepth creates a SwaggerOption to specify swagger defaultModelExpandDepth, defaults to 1.
func WithSwaggerDefaultModelExpandDepth(defaultModelExpandDepth int) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DefaultModelExpandDepth = &defaultModelExpandDepth
	}
}

// WithSwaggerDefaultModelRendering creates a SwaggerOption to specify swagger defaultModelRendering, "example" and "model" are valid, defaults to "example".
func WithSwaggerDefaultModelRendering(defaultModelRendering string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DefaultModelRendering = &defaultModelRendering
	}
}

// WithSwaggerDisplayRequestDuration creates a SwaggerOption to specify swagger displayRequestDuration, defaults to false.
func WithSwaggerDisplayRequestDuration(displayRequestDuration bool) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DisplayRequestDuration = &displayRequestDuration
	}
}

// WithSwaggerDocExpansion creates a SwaggerOption to specify swagger docExpansion, "list" and "full" and "none" are valid, defaults to "list".
func WithSwaggerDocExpansion(docExpansion string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.DocExpansion = &docExpansion
	}
}

// WithSwaggerMaxDisplayedTags creates a SwaggerOption to specify swagger maxDisplayedTags, defaults to show all operations.
func WithSwaggerMaxDisplayedTags(maxDisplayedTags int) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.MaxDisplayedTags = &maxDisplayedTags
	}
}

// WithSwaggerOperationsSorter creates a SwaggerOption to specify swagger operationsSorter, "" and "alpha" and "method" are valid, defaults to "" which means not to sort.
func WithSwaggerOperationsSorter(operationsSorter string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.OperationsSorter = &operationsSorter
	}
}

// WithSwaggerShowExtensions creates a SwaggerOption to specify swagger showExtensions, defaults to false.
func WithSwaggerShowExtensions(showExtensions bool) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.ShowExtensions = &showExtensions
	}
}

// WithSwaggerShowCommonExtensions creates a SwaggerOption to specify swagger showCommonExtensions (pattern, maxLength, minLength, maximum, minimum), defaults to false.
func WithSwaggerShowCommonExtensions(showCommonExtensions bool) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.ShowCommonExtensions = &showCommonExtensions
	}
}

// WithSwaggerTagsSorter creates a SwaggerOption to specify swagger tagsSorter, "" and "alpha" is valid, default to "" which means not to sort.
func WithSwaggerTagsSorter(tagsSorter string) SwaggerOption {
	return func(o *SwaggerOptions) {
		o.TagsSorter = &tagsSorter
	}
}

// WrapSwagger registers swagger related routes to gin router, the documentation page will be served in `index.html` or the root of given router.
// Visit https://github.com/swagger-api/swagger-ui/blob/v3.17.2/docs/usage/configuration.md for SwaggerOptions' details. Note: this function is
// not allowed to be used on gin.Engine directly or any other gin.GroupEngine which has empty relativePath.
//
// Example:
// 	func ReadSwaggerDoc() []byte {
// 		bs, _ := os.ReadFile("./api/spec.json")
// 		return bs
// 	}
// 	xgin.WrapSwagger(engine.Group("/v1/swagger"), ReadSwaggerDoc) // => /v1/swagger/index.html
func WrapSwagger(router gin.IRouter, swaggerDocGetter func() []byte, swaggerOptions ...SwaggerOption) {
	if _, ok := router.(*gin.Engine); ok {
		panic(panicWrapSwaggerToEngine)
	}
	if swaggerDocGetter == nil {
		panic(panicNilSwaggerDocGetter)
	}

	opt := &SwaggerOptions{EnableRedirect: xpointer.BoolPtr(true)}
	for _, o := range swaggerOptions {
		if o != nil {
			o(opt)
		}
	}
	indexHtmlName := strings.TrimSpace(opt.IndexHtmlRouteName)
	if indexHtmlName == "" {
		indexHtmlName = fmt.Sprintf("index.html")
	}
	docJsonName := strings.TrimSpace(opt.DocJsonRouteName)
	if docJsonName == "" {
		docJsonName = fmt.Sprintf("doc.json")
	}
	configJsonName := strings.TrimSpace(opt.ConfigJsonRouteName)
	if configJsonName == "" {
		configJsonName = fmt.Sprintf("config.json")
	}

	contentTypeMap := map[string]string{
		".json": "application/json",
		".yaml": "application/yaml",
		".png":  "image/png",
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript",
	}

	if *opt.EnableRedirect {
		router.GET("", func(c *gin.Context) {
			indexUrl := c.FullPath() + "/" + indexHtmlName // => ".../swagger/index.html"
			c.Redirect(http.StatusMovedPermanently, indexUrl)
		})
	}

	router.GET("*file", func(c *gin.Context) {
		var data []byte
		var contentType string
		realUrl := c.Request.URL.String()
		for k, v := range contentTypeMap {
			if strings.HasSuffix(realUrl, k) {
				contentType = v
			}
		}

		pureUrl := strings.TrimSuffix(c.FullPath(), "/*file")
		file := strings.TrimSpace(strings.TrimPrefix(c.Param("file"), "/"))
		if file == "" && *opt.EnableRedirect {
			indexUrl := pureUrl + "/" + indexHtmlName // => ".../swagger/index.html"
			c.Redirect(http.StatusMovedPermanently, indexUrl)
			return
		}

		switch file {
		case docJsonName: // <<<
			data = swaggerDocGetter()
		case configJsonName:
			data = []byte("{}")
			if d, _ := json.Marshal(opt); d != nil {
				data = d
			}
		case indexHtmlName: // <<<
			docUrl := pureUrl + "/" + docJsonName       // => ".../swagger/doc.json"
			configUrl := pureUrl + "/" + configJsonName // => ".../swagger/config.json"
			data = bytes.Replace(_embedded_index_html, []byte("$$URL"), []byte(docUrl), 1)
			data = bytes.Replace(data, []byte("$$CONFIG_URL"), []byte(configUrl), 1)

		// static resources
		case "favicon-16x16.png":
			data = _embedded_favicon_16x16_png
		case "favicon-32x32.png":
			data = _embedded_favicon_32x32_png
		case "oauth2-redirect.html":
			data = _embedded_oauth2_redirect_html
		case "swagger-ui.css":
			data = _embedded_swagger_ui_css
		case "swagger-ui.js":
			data = _embedded_swagger_ui_js
		case "swagger-ui-bundle.js":
			data = _embedded_swagger_ui_bundle_js
		case "swagger-ui-standalone-preset.js":
			data = _embedded_swagger_ui_standalone_preset_js
		}

		if data == nil {
			c.Data(404, "text/plain; charset=utf-8", []byte("404 page not found"))
		} else {
			c.Data(200, contentType, data)
		}
	})
}

// WrapSwaggerSilently registers swagger related routes to gin router without any debug logging.
func WrapSwaggerSilently(router gin.IRouter, swaggerDocGetter func() []byte, swaggerOptions ...SwaggerOption) {
	restore := HideDebugPrintRoute()
	WrapSwagger(router, swaggerDocGetter, swaggerOptions...)
	restore()
}

// The following embedded files are all from `https://github.com/swagger-api/swagger-ui/releases/tag/v3.17.2`, which are used by WrapSwagger.
var (
	//go:embed swaggerDist/favicon-16x16.png
	_embedded_favicon_16x16_png []byte
	//go:embed swaggerDist/favicon-32x32.png
	_embedded_favicon_32x32_png []byte
	//go:embed swaggerDist/index.html
	_embedded_index_html []byte
	//go:embed swaggerDist/oauth2-redirect.html
	_embedded_oauth2_redirect_html []byte
	//go:embed swaggerDist/swagger-ui.css
	_embedded_swagger_ui_css []byte
	//go:embed swaggerDist/swagger-ui.js
	_embedded_swagger_ui_js []byte
	//go:embed swaggerDist/swagger-ui-bundle.js
	_embedded_swagger_ui_bundle_js []byte
	//go:embed swaggerDist/swagger-ui-standalone-preset.js
	_embedded_swagger_ui_standalone_preset_js []byte
)

// ========================
// mass functions and types
// ========================

// GetTrustedProxies returns trusted proxies string slice from given gin.Engine, returns nil if given nil engine.
func GetTrustedProxies(engine *gin.Engine) []string {
	if engine == nil {
		return nil
	}
	var val = xreflect.GetUnexportedField(xreflect.FieldValueOf(engine, "trustedProxies"))
	if val.IsNil() {
		return nil
	}
	return val.Interface().([]string)
}

// HideDebugLogging hides gin's all logging (gin.DefaultWriter and gin.DefaultErrorWriter) and returns a function to restore this behavior.
func HideDebugLogging() (restoreFn func()) {
	originWriter := gin.DefaultWriter
	originErrorWriter := gin.DefaultErrorWriter
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	return func() {
		gin.DefaultWriter = originWriter
		gin.DefaultErrorWriter = originErrorWriter
	}
}

// HideDebugPrintRoute hides gin.DebugPrintRouteFunc logging and returns a function to restore this behavior.
func HideDebugPrintRoute() (restoreFn func()) {
	originFunc := gin.DebugPrintRouteFunc
	gin.DebugPrintRouteFunc = SilentPrintRouteFunc
	return func() {
		gin.DebugPrintRouteFunc = originFunc
	}
}

// DebugPrintRouteFuncType represents the function type of gin.DebugPrintRouteFunc.
type DebugPrintRouteFuncType func(httpMethod, absolutePath, handlerName string, numHandlers int)

// SilentPrintRouteFunc is a global gin.DebugPrintRouteFunc, which do nothing for logging, means silent.
func SilentPrintRouteFunc(_, _, _ string, _ int) {
	// pass
}

// DefaultPrintRouteFunc is the default gin.DebugPrintRouteFunc, can be modified by overwriting value to gin.DebugPrintRouteFunc.
//
// The default format logs like (just like gin.DebugPrintRouteFunc except [Gin] prefix):
// 	[Gin] GET    /debug/pprof/             --> ... (1 handlers)
// 	[Gin] GET    /debug/pprof/threadcreate --> ... (1 handlers)
// 	[Gin] POST   /debug/pprof/symbol       --> ... (1 handlers)
// 	     |------|-------------------------|   |---|
// 	        6               25                 ...
func DefaultPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int) {
	fmt.Printf("[Gin] %-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, numHandlers)
}

// DefaultColorizedPrintRouteFunc is the colorized version of DefaultPrintRouteFunc.
//
// The default format logs like (just like gin.DebugPrintRouteFunc except [Gin] prefix):
// 	[Gin] GET    /debug/pprof/             --> ... (1 handlers)
// 	[Gin] GET    /debug/pprof/threadcreate --> ... (1 handlers)
// 	[Gin] POST   /debug/pprof/symbol       --> ... (1 handlers)
// 	     |------|-------------------------|   |---|
// 	     6 (blue)       25 (blue)              ...
func DefaultColorizedPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int) {
	fmt.Printf("[Gin] %s --> %s (%d handlers)\n", xcolor.Blue.Sprintf("%-6s %-25s", httpMethod, absolutePath), handlerName, numHandlers)
}

const (
	panicNilErrorForRouterDecodeError = "xgin: nil error for RouterDecodeError"
)

// RouterDecodeError is an error type for router parameter decoding. At most of the time, the Err field is expected to be strconv.NumError type generated by strconv
// functions such as strconv.ParseInt and strconv.Atoi. This type also supports custom translation behavior in TranslateBindingError with WithXginRouterDecodeError.
type RouterDecodeError struct {
	Field   string
	Input   string
	Err     error
	Message string
}

// NewRouterDecodeError creates a new RouterDecodeError by given arguments, panics when using nil error.
func NewRouterDecodeError(field string, input string, err error, message string) *RouterDecodeError {
	if err == nil {
		panic(panicNilErrorForRouterDecodeError)
	}
	return &RouterDecodeError{Field: field, Input: input, Err: err, Message: message}
}

// Error returns the formatted error message from RouterDecodeError. Note that returned value does not contain custom message.
func (r *RouterDecodeError) Error() string {
	if r.Field == "" {
		return fmt.Sprintf("parsing \"%s\": %v", r.Input, r.Err)
	}
	return fmt.Sprintf("parsing %s \"%s\": %v", r.Field, r.Input, r.Err)
}

// Unwrap returns the wrapped error from RouterDecodeError.
func (r *RouterDecodeError) Unwrap() error {
	return r.Err
}
