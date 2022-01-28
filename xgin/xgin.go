package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"os"
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

// isSpecificHeader checks whether the given param is the same specific header in case-insensitive.
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

// ==========
// wrap pprof
// ==========

var (
	_indexHandler        = func(ctx *gin.Context) { pprof.Index(ctx.Writer, ctx.Request) }
	_heapHandler         = func(ctx *gin.Context) { pprof.Handler("heap").ServeHTTP(ctx.Writer, ctx.Request) }
	_goroutineHandler    = func(ctx *gin.Context) { pprof.Handler("goroutine").ServeHTTP(ctx.Writer, ctx.Request) }
	_allocsHandler       = func(ctx *gin.Context) { pprof.Handler("allocs").ServeHTTP(ctx.Writer, ctx.Request) }
	_blockHandler        = func(ctx *gin.Context) { pprof.Handler("block").ServeHTTP(ctx.Writer, ctx.Request) }
	_threadcreateHandler = func(ctx *gin.Context) { pprof.Handler("threadcreate").ServeHTTP(ctx.Writer, ctx.Request) }
	_cmdlineHandler      = func(ctx *gin.Context) { pprof.Cmdline(ctx.Writer, ctx.Request) }
	_profileHandler      = func(ctx *gin.Context) { pprof.Profile(ctx.Writer, ctx.Request) }
	_symbolHandler       = func(ctx *gin.Context) { pprof.Symbol(ctx.Writer, ctx.Request) }
	_traceHandler        = func(ctx *gin.Context) { pprof.Trace(ctx.Writer, ctx.Request) }
	_mutexHandler        = func(ctx *gin.Context) { pprof.Handler("mutex").ServeHTTP(ctx.Writer, ctx.Request) }
)

// WrapPprof registers several routes from package `net/http/pprof` to gin.Engine. For more, please visit https://github.com/DeanThompson/ginpprof.
func WrapPprof(engine *gin.Engine) {
	for _, r := range []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/debug/pprof/", _indexHandler},
		{"GET", "/debug/pprof/heap", _heapHandler},
		{"GET", "/debug/pprof/goroutine", _goroutineHandler},
		{"GET", "/debug/pprof/allocs", _allocsHandler},
		{"GET", "/debug/pprof/block", _blockHandler},
		{"GET", "/debug/pprof/threadcreate", _threadcreateHandler},
		{"GET", "/debug/pprof/cmdline", _cmdlineHandler},
		{"GET", "/debug/pprof/profile", _profileHandler},
		{"GET", "/debug/pprof/symbol", _symbolHandler},
		{"POST", "/debug/pprof/symbol", _symbolHandler},
		{"GET", "/debug/pprof/trace", _traceHandler},
		{"GET", "/debug/pprof/mutex", _mutexHandler},
	} {
		engine.Handle(r.method, r.path, r.handler) // use path directly
	}
}

// ========================
// mass functions and types
// ========================

// SetPrintRouteFunc sets gin.Engine's debug print route func by modifying gin.DebugPrintRouteFunc directly, defaults to DefaultPrintRouteFunc.
func SetPrintRouteFunc(f func(httpMethod, absolutePath, handlerName string, numHandlers int)) {
	gin.DebugPrintRouteFunc = f
}

// HideDebugPrintRoute hides the gin.DebugPrintRouteFunc logging and returns a function to restore this behavior.
func HideDebugPrintRoute() (restoreFn func()) {
	originFunc := gin.DebugPrintRouteFunc
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	return func() {
		gin.DebugPrintRouteFunc = originFunc
	}
}

func init() {
	// set gin.DebugPrintRouteFunc to DefaultPrintRouteFunc in default.
	SetPrintRouteFunc(DefaultPrintRouteFunc)
}

// DefaultPrintRouteFunc is the default gin.DebugPrintRouteFunc, can be modified by SetPrintRouteFunc.
//
// The default format logs like (just like gin.DebugPrintRouteFunc):
// 	[Gin-debug] GET    /debug/pprof/             --> ... (1 handlers)
// 	[Gin-debug] GET    /debug/pprof/threadcreate --> ... (1 handlers)
// 	[Gin-debug] POST   /debug/pprof/symbol       --> ... (1 handlers)
func DefaultPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int) {
	fmt.Printf("[Gin-debug] %-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, numHandlers)
}

// DefaultColorizedPrintRouteFunc is the DefaultPrintRouteFunc in color.
//
// The default format logs like (just like gin.DebugPrintRouteFunc):
// 	[Gin-debug]
func DefaultColorizedPrintRouteFunc(httpMethod, absolutePath, handlerName string, numHandlers int) {
	// TODO
}

// NewEngineWithoutLogging creates a new gin.Engine without gin.debugPrintWARNINGNew behavior.
func NewEngineWithoutLogging() *gin.Engine {
	originWriter := gin.DefaultWriter
	gin.DefaultWriter = io.Discard
	engine := gin.New()
	gin.DefaultWriter = originWriter
	return engine
}

// RouterDecodeError is an error type for router parameter decoding. At most of the time, the Err field is in strconv.NumError type generated by functions from strconv
// package such as strconv.ParseInt and strconv.Atoi. This type also supports custom translation in TranslateBindingError and WithXginRouterDecodeError.
type RouterDecodeError struct {
	Field   string
	Input   string
	Err     error
	Message string
}

// NewRouterDecodeError creates a new RouterDecodeError by given arguments.
func NewRouterDecodeError(field string, input string, err error, message string) *RouterDecodeError {
	return &RouterDecodeError{Field: field, Input: input, Err: err, Message: message}
}

// Error returns the formatted error message from RouterDecodeError, note that returned value is not RouterDecodeError.Message.
func (r *RouterDecodeError) Error() string {
	// if nErr, ok := r.Err.(*strconv.NumError); ok {
	// 	return nErr.Error()
	// }
	return fmt.Sprintf("parsing %s \"%s\": %v", r.Field, r.Input, r.Err)
}

// Unwrap returns the wrapped error from RouterDecodeError.
func (r *RouterDecodeError) Unwrap() error {
	return r.Err
}

// GetProxyEnv lookups and returns three proxy environments, including http_proxy, https_proxy and socks_proxy.
func GetProxyEnv() (httpProxy string, httpsProxy string, socksProxy string) {
	hp := strings.TrimSpace(os.Getenv("http_proxy"))
	hsp := strings.TrimSpace(os.Getenv("https_proxy"))
	ssp := strings.TrimSpace(os.Getenv("socks_proxy"))
	return hp, hsp, ssp
}
