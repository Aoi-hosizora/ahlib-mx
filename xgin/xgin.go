package xgin

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"strings"
)

// ============
// dump request
// ============

// dumpRequestOptions represents some options for DumpRequest, set by DumpRequestOption.
type dumpRequestOptions struct {
	ignoreRequestLine bool
	retainHeaders     []string
	ignoreHeaders     []string
	secretHeaders     []string
	secretReplace     string
}

// DumpRequestOption represents an option for DumpRequest, can be created by WithXXX functions.
type DumpRequestOption func(*dumpRequestOptions)

// WithIgnoreRequestLine creates a DumpRequestOption for request line, if set to true, request line such as "GET /xxx HTTP/1.1" will be ignored.
func WithIgnoreRequestLine(ignore bool) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.ignoreRequestLine = ignore
	}
}

// WithRetainHeaders creates a DumpRequestOption for retained header.
func WithRetainHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.retainHeaders = headers
	}
}

// WithIgnoreHeaders creates a DumpRequestOption for ignore headers, this option will be ignored when WithRetainHeaders is used in DumpRequest.
func WithIgnoreHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.ignoreHeaders = headers
	}
}

// WithSecretHeaders creates a DumpRequestOption for secret headers, such as Authorization field.
func WithSecretHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.secretHeaders = headers
	}
}

// WithSecretReplace creates a DumpRequestOption for secret header replace string, defaults to "*".
func WithSecretReplace(secret string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.secretReplace = secret
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
	opt := &dumpRequestOptions{ignoreRequestLine: false, secretReplace: "*"}
	for _, op := range options {
		if op != nil {
			op(opt)
		}
	}

	bs, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil
	}
	params := strings.Split(xstring.FastBtos(bs), "\r\n") // split by \r\n
	result := make([]string, 0, len(params))
	for idx, param := range params {
		if idx == 0 {
			if !opt.ignoreRequestLine {
				result = append(result, param) // request line: METHOD /ENDPOINT HTTP/1.1
			}
			continue
		}
		param = strings.TrimSpace(param)
		if param == "" {
			// after the request line, there is \r\n\r\n, which is the splitter of the request line and the request header
			continue
		}

		// filter headers, use retainHeaders first
		headerList := opt.retainHeaders
		if len(opt.retainHeaders) == 0 {
			headerList = opt.ignoreHeaders
		}
		exist := false
		for _, header := range headerList {
			if isSpecificHeader(param, header) {
				exist = true
				break
			}
		}
		if (len(opt.retainHeaders) != 0 && !exist) || (len(opt.retainHeaders) == 0 && exist) {
			continue
		}

		// rewrite headers that are secret
		for _, header := range opt.secretHeaders {
			if isSpecificHeader(param, header) {
				header = strings.SplitN(param, ":", 2)[0]
				param = header + ": " + opt.secretReplace
				break
			}
		}

		// append to the result slice
		result = append(result, param)
	}

	return result
}

// ==========
// pprof wrap
// ==========

// PprofWrap registers several routes from package `net/http/pprof` to gin.Engine. For more, please visit https://github.com/DeanThompson/ginpprof.
func PprofWrap(engine *gin.Engine, hideDebug bool) {
	if hideDebug {
		temp := gin.DebugPrintRouteFunc
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
		defer func() { gin.DebugPrintRouteFunc = temp }()
	}
	for _, r := range []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/debug/pprof/", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/heap", func(ctx *gin.Context) {
			pprof.Handler("heap").ServeHTTP(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/goroutine", func(ctx *gin.Context) {
			pprof.Handler("goroutine").ServeHTTP(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/allocs", func(ctx *gin.Context) {
			pprof.Handler("allocs").ServeHTTP(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/block", func(ctx *gin.Context) {
			pprof.Handler("block").ServeHTTP(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/threadcreate", func(ctx *gin.Context) {
			pprof.Handler("threadcreate").ServeHTTP(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/cmdline", func(ctx *gin.Context) {
			pprof.Cmdline(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/profile", func(ctx *gin.Context) {
			pprof.Profile(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/symbol", func(ctx *gin.Context) {
			pprof.Symbol(ctx.Writer, ctx.Request)
		}},
		{"POST", "/debug/pprof/symbol", func(ctx *gin.Context) {
			pprof.Symbol(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/trace", func(ctx *gin.Context) {
			pprof.Trace(ctx.Writer, ctx.Request)
		}},
		{"GET", "/debug/pprof/mutex", func(ctx *gin.Context) {
			pprof.Handler("mutex").ServeHTTP(ctx.Writer, ctx.Request)
		}},
	} {
		engine.Handle(r.method, r.path, r.handler) // use path directly
	}
}

// ================================
// validator & translator & binding
// ================================

var (
	errValidatorNotSupported = errors.New("xgin: gin's validator engine is not validator.Validate from github.com/go-playground/validator/v10")
)

// GetValidatorEngine returns gin's binding validator engine, which only supports validator.Validate from github.com/go-playground/validator/v10.
func GetValidatorEngine() (*validator.Validate, error) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errValidatorNotSupported
	}
	return val, nil
}

// GetValidatorTranslator applies and returns xvalidator.UtTranslator for validator.Validate using given parameters. Also see xvalidator.ApplyTranslator.
//
// Example:
// 	translator, _ := xgin.GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
// 	result := err.(validator.ValidationErrors).Translate(translator)
func GetValidatorTranslator(locale xvalidator.LocaleTranslator, registerFn xvalidator.TranslationRegisterHandler) (xvalidator.UtTranslator, error) {
	val, err := GetValidatorEngine()
	if err != nil {
		return nil, err // errValidatorNotSupported
	}
	return xvalidator.ApplyTranslator(val, locale, registerFn) // create translator with locale, register translator to validator
}

// _globalTranslator is a global xvalidator.UtTranslator set by SetGlobalTranslator and can be got by GetGlobalTranslator.
var _globalTranslator xvalidator.UtTranslator

// SetGlobalTranslator stores given xvalidator.UtTranslator to global, it can be got by using GetGlobalTranslator.
func SetGlobalTranslator(translator xvalidator.UtTranslator) {
	_globalTranslator = translator
}

// GetGlobalTranslator gets the stored translator by SetGlobalTranslator, it will return nil if this function is called before SetGlobalTranslator.
func GetGlobalTranslator() xvalidator.UtTranslator {
	return _globalTranslator
}

// AddBinding registers custom validation function to gin's validator engine. You can use your custom validator.Func or functions from xvalidator package
// such as xvalidator.RegexpValidator and xvalidator.DateTimeValidator.
//
// Example:
// 	err := xgin.AddBinding("regexp", xvalidator.ParamRegexpValidator())
// 	err := xgin.AddBinding("xxx", func(fl validator.FieldLevel) bool { /* ... */ })
func AddBinding(tag string, fn validator.Func) error {
	v, err := GetValidatorEngine()
	if err != nil {
		return err
	}
	return v.RegisterValidation(tag, fn)
}

// AddTranslation registers custom validation translation to gin's validator engine, using given tag, message and override flag. Also see xvalidator.DefaultTranslateFunc.
//
// Example:
// 	err := xgin.AddTranslation(translator, "regexp", "{0} must match regexp /{1}/", true)
// 	err := xgin.AddTranslation(translator, "email", "{0} must be an email", true)
func AddTranslation(translator xvalidator.UtTranslator, tag, message string, override bool) error {
	v, err := GetValidatorEngine()
	if err != nil {
		return err
	}
	regisFn := xvalidator.DefaultRegistrationFunc(tag, message, override)
	transFn := xvalidator.DefaultTranslateFunc()
	return v.RegisterTranslation(tag, translator, regisFn, transFn)
}

// EnableParamRegexpBinding enables parameterized regexp validator to `regexp` binding tag, see xvalidator.ParamRegexpValidator.
func EnableParamRegexpBinding() error {
	return AddBinding("regexp", xvalidator.ParamRegexpValidator())
}

// EnableParamRegexpBindingTranslator enables parameterized regexp validator `regexp`'s translation using given xvalidator.UtTranslator.
func EnableParamRegexpBindingTranslator(translator xvalidator.UtTranslator) error {
	return AddTranslation(translator, "regexp", "{0} should match regexp /{1}/", true)
}

// EnableRFC3339DateBinding enables rfc3339 date validator to `date` binding tag, see xvalidator.DateTimeValidator.
func EnableRFC3339DateBinding() error {
	return AddBinding("date", xvalidator.DateTimeValidator(xtime.RFC3339Date))
}

// EnableRFC3339DateBindingTranslator enables rfc3339 date validator `date`'s translation using given xvalidator.UtTranslator.
func EnableRFC3339DateBindingTranslator(translator xvalidator.UtTranslator) error {
	return AddTranslation(translator, "date", "{0} should be an RFC3339 date", true)
}

// EnableRFC3339DateTimeBinding enables rfc3339 datetime validator to `datetime` binding tag, see xvalidator.DateTimeValidator.
func EnableRFC3339DateTimeBinding() error {
	return AddBinding("datetime", xvalidator.DateTimeValidator(xtime.RFC3339DateTime))
}

// EnableRFC3339DateTimeBindingTranslator enables rfc3339 datetime validator `datetime`'s translation using given xvalidator.UtTranslator.
func EnableRFC3339DateTimeBindingTranslator(translator xvalidator.UtTranslator) error {
	return AddTranslation(translator, "datetime", "{0} should be an RFC3339 datetime", true)
}

// ============
// router error
// ============

// RouterDecodeError is an error type for router parameter decoding. At most of the time, the Err field is in strconv.NumError type generated by functions from
// strconv package such as strconv.ParseInt and strconv.Atoi.
type RouterDecodeError struct {
	Field   string
	Input   string
	Err     error
	Message string
}

// NewRouterDecodeError creates a new RouterDecodeError by parameters.
func NewRouterDecodeError(routerField string, input string, err error, message string) *RouterDecodeError {
	return &RouterDecodeError{Field: routerField, Input: input, Err: err, Message: message}
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
