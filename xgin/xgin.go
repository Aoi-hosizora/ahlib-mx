package xgin

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http/httputil"
	"net/http/pprof"
	"strings"
)

// ============
// dump request
// ============

// dumpRequestOptions represents some options for DumpRequest, set by DumpRequestOption.
type dumpRequestOptions struct {
	retainHeaders []string
	ignoreHeaders []string
	secretHeaders []string
	secretReplace string
}

// DumpRequestOption represents an option for DumpRequest, can be created by WithXXX functions.
type DumpRequestOption func(*dumpRequestOptions)

// WithRetainHeaders creates a DumpRequestOption for retained header. Set this option will make DumpRequest ignore the WithIgnoreHeaders option.
func WithRetainHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.retainHeaders = headers
	}
}

// WithIgnoreHeaders creates a DumpRequestOption for ignore headers. This option will be ignored when WithRetainHeaders is used in DumpRequest.
func WithIgnoreHeaders(headers ...string) DumpRequestOption {
	return func(o *dumpRequestOptions) {
		o.ignoreHeaders = headers
	}
}

// WithSecretHeaders creates a DumpRequestOption for secret headers, such as Authorization.
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

// DumpRequest dumps and formats http.Request from gin.Context to string slice, using given DumpRequestOption-s. The first element must be request line
// "METHOD /ENDPOINT HTTP/1.1", and the remaining elements are the request headers "XXX: YYY", returns an empty slice when using nil gin.Context.
func DumpRequest(c *gin.Context, options ...DumpRequestOption) []string {
	if c == nil {
		return make([]string, 0)
	}

	opt := &dumpRequestOptions{secretReplace: "*"}
	for _, op := range options {
		if op != nil {
			op(opt)
		}
	}

	bs, err := httputil.DumpRequest(c.Request, false) // ignore error
	if err != nil {
		return make([]string, 0)
	}
	params := strings.Split(xstring.FastBtos(bs), "\r\n") // split by \r\n
	result := make([]string, 0, len(params))
	for idx, param := range params {
		if idx == 0 {
			result = append(result, param) // METHOD /ENDPOINT HTTP/1.1
			continue
		}
		param = strings.TrimSpace(param)
		if param == "" {
			// after the first line, there is \r\n\r\n, and has a blank line
			continue
		}

		// headers
		if len(opt.retainHeaders) != 0 { // use retainHeaders to filter
			exists := false
			for _, header := range opt.retainHeaders {
				if strings.HasPrefix(param, header+": ") {
					exists = true
					break
				}
			}
			if !exists {
				continue
			}
		} else { // use ignoreHeaders to filter
			exists := false
			for _, header := range opt.ignoreHeaders {
				if strings.HasPrefix(param, header+": ") {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
		}
		for _, header := range opt.secretHeaders { // rewrite header that is secret
			if strings.HasPrefix(param, header+": ") {
				param = header + ": " + opt.secretReplace
				break
			}
		}

		// append
		result = append(result, param)
	}

	return result
}

// =====
// pprof
// =====

// PprofWrap adds several routes from package `net/http/pprof` to gin.Engine. Reference from https://github.com/DeanThompson/ginpprof.
func PprofWrap(router *gin.Engine) {
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
		router.Handle(r.method, r.path, r.handler) // use path directly
	}
}

// ================================
// validator & translator & binding
// ================================

var (
	errValidatorNotSupported = errors.New("xgin: gin's validator engine is not github.com/go-playground/validator/v10")
)

// GetValidatorEngine returns gin's binding validator engine, which only supports validator.Validate from github.com/go-playground/validator/v10.
// Also see binding.Validator.
func GetValidatorEngine() (*validator.Validate, error) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errValidatorNotSupported
	}
	return val, nil
}

// GetValidatorTranslator applies and returns ut.Translator for validator.Validate using given locales.Translator and xvalidator.TranslationRegisterHandler.
// Also see xvalidator.ApplyTranslator.
//
// Example:
// 	translator, _ := xgin.GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
// 	result := validator.New().Struct(&testStruct{}).(validator.ValidationErrors).Translate(translator)
func GetValidatorTranslator(locTranslator locales.Translator, registerFn xvalidator.TranslationRegisterHandler) (ut.Translator, error) {
	val, err := GetValidatorEngine()
	if err != nil {
		return nil, err // errValidatorNotSupported
	}
	return xvalidator.ApplyTranslator(val, locTranslator, registerFn) // create translator and do register
}

// AddBinding adds user defined binding to gin's validator engine. You can use your custom validator.Func or functions provided by xvalidator's such as
// xvalidator.RegexpValidator and xvalidator.DateTimeValidator.
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

// AddTranslator adds user defined validator's translator to given ut.Translator using given tag, message and override. Also see
// xvalidator.AddToTranslatorFunc and xvalidator.DefaultTranslateFunc.
//
// Example:
// 	err := xgin.AddTranslator(translator, "regexp", "{0} must matches regexp /{1}/", true)
// 	err := xgin.AddTranslator(translator, "email", "{0} must be an email", true)
func AddTranslator(translator ut.Translator, tag, message string, override bool) error {
	v, err := GetValidatorEngine()
	if err != nil {
		return err
	}
	fn := xvalidator.AddToTranslatorFunc(tag, message, override)
	return v.RegisterTranslation(tag, translator, fn, xvalidator.DefaultTranslateFunc())
}

// EnableParamRegexpBinding enables parametered regexp validator to `regexp`, see xvalidator.ParamRegexpValidator.
func EnableParamRegexpBinding() error {
	return AddBinding("regexp", xvalidator.ParamRegexpValidator())
}

// EnableParamRegexpBindingTranslator enables parametered regexp validator (`regexp`)'s translator to given ut.Translator.
func EnableParamRegexpBindingTranslator(translator ut.Translator) error {
	return AddTranslator(translator, "regexp", "{0} must match regexp /{1}/", true)
}

// EnableRFC3339DateBinding enables rfc3339 date validator to `date`, see xvalidator.DateTimeValidator.
func EnableRFC3339DateBinding() error {
	return AddBinding("date", xvalidator.DateTimeValidator(xtime.RFC3339Date))
}

// EnableRFC3339DateBindingTranslator enables rfc3339 date validator (`date`)'s translator to given ut.Translator.
func EnableRFC3339DateBindingTranslator(translator ut.Translator) error {
	return AddTranslator(translator, "date", "{0} must be an RFC3339 date", true)
}

// EnableRFC3339DateTimeBinding enables rfc3339 datetime validator to `datetime`, see xvalidator.DateTimeValidator.
func EnableRFC3339DateTimeBinding() error {
	return AddBinding("datetime", xvalidator.DateTimeValidator(xtime.RFC3339DateTime))
}

// EnableRFC3339DateTimeBindingTranslator enables rfc3339 datetime validator (`datetime`)'s translator to given ut.Translator.
func EnableRFC3339DateTimeBindingTranslator(translator ut.Translator) error {
	return AddTranslator(translator, "datetime", "{0} must be an RFC3339 datetime", true)
}
