package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"strings"
)

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

// WithSecretPlaceholder creates a DumpRequestOption to specific a secret placeholder for secret headers set by WithSecretHeaders, defaults to "*".
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
	for _, op := range options {
		if op != nil {
			op(opt)
		}
	}
	if opt.secretPlaceholder == "" {
		opt.secretPlaceholder = "*"
	}

	bs, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil // unreachable
	}
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

// =============
// route & pprof
// =============

// RedirectHandler creates a gin.HandlerFunc that behaviors a redirection with given code (such as http.StatusMovedPermanently or http.StatusTemporaryRedirect)
// and redirect target location.
func RedirectHandler(code int, location string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(code, location)
	}
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

// WrapPprof registers several routes from package `net/http/pprof` to gin.Engine. For more, please visit https://github.com/DeanThompson/ginpprof.
func WrapPprof(engine *gin.Engine) {
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
		engine.Handle(r.method, r.path, r.handler) // use path directly
	}
}

// ========================
// mass functions and types
// ========================

// HideDebugLogging hides gin's all loggings and returns a function to restore this behavior.
func HideDebugLogging() (restoreFn func()) {
	originWriter := gin.DefaultWriter
	gin.DefaultWriter = io.Discard
	return func() {
		gin.DefaultWriter = originWriter
	}
}

// HideDebugPrintRoute hides gin.DebugPrintRouteFunc logging and returns a function to restore this behavior.
func HideDebugPrintRoute() (restoreFn func()) {
	originFunc := gin.DebugPrintRouteFunc
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	return func() {
		gin.DebugPrintRouteFunc = originFunc
	}
}

// SetPrintRouteFunc sets gin.Engine's debug print route func by modifying gin.DebugPrintRouteFunc directly, defaults to DefaultPrintRouteFunc.
func SetPrintRouteFunc(f func(httpMethod, absolutePath, handlerName string, numHandlers int)) {
	gin.DebugPrintRouteFunc = f
}

func init() {
	// set gin's mode to gin.DebugMode, and set debugPrintRouteFunc to DefaultPrintRouteFunc in default
	gin.SetMode(gin.DebugMode)
	SetPrintRouteFunc(DefaultPrintRouteFunc)
}

// DefaultPrintRouteFunc is the default gin.DebugPrintRouteFunc, can be modified by SetPrintRouteFunc.
//
// The default format logs like (just like gin.DebugPrintRouteFunc):
// 	[Gin-debug] GET    /debug/pprof/             --> ... (1 handlers)
// 	[Gin-debug] GET    /debug/pprof/threadcreate --> ... (1 handlers)
// 	[Gin-debug] POST   /debug/pprof/symbol       --> ... (1 handlers)
// 	           |------|-------------------------|   |---|
// 	              6               25                 ...
func DefaultPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int) {
	fmt.Printf("[Gin-debug] %-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, numHandlers)
}

// DefaultColorizedPrintRouteFunc is the DefaultPrintRouteFunc in color.
//
// The default format logs like (just like gin.DebugPrintRouteFunc):
// 	[Gin-debug] GET    /debug/pprof/             --> ... (1 handlers)
// 	[Gin-debug] GET    /debug/pprof/threadcreate --> ... (1 handlers)
// 	[Gin-debug] POST   /debug/pprof/symbol       --> ... (1 handlers)
// 	           |------|-------------------------|   |---|
// 	           6 (blue)       25 (blue)              ...
func DefaultColorizedPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int) {
	fmt.Printf("[Gin-debug] %s --> %s (%d handlers)\n", xcolor.Blue.Sprintf("%-6s %-25s", httpMethod, absolutePath), handlerName, numHandlers)
}

const (
	panicNilError = "xgin: nil error for RouterDecodeError"
)

// RouterDecodeError is an error type for router parameter decoding. At most of the time, the Err field is in strconv.NumError type generated by functions from
// strconv package such as strconv.ParseInt and strconv.Atoi. This type also supports custom translation in TranslateBindingError with WithXginRouterDecodeError.
type RouterDecodeError struct {
	Field   string
	Input   string
	Err     error
	Message string
}

// NewRouterDecodeError creates a new RouterDecodeError by given arguments, panics when using nil error.
func NewRouterDecodeError(field string, input string, err error, message string) *RouterDecodeError {
	if err == nil {
		panic(panicNilError)
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
