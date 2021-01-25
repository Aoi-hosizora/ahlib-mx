package xgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http/httputil"
	"net/http/pprof"
	"strings"
)

// ====
// dump
// ====

// DumpRequest dumps and formats request strings from gin.Context, using given ignore headers such as Authorization,
// returns an empty slice when using nil gin.Context.
func DumpRequest(c *gin.Context, ignoreHeaders []string) []string {
	result := make([]string, 0)
	if c == nil {
		return result
	}

	bs, _ := httputil.DumpRequest(c.Request, false) // ignore error
	params := strings.Split(string(bs), "\r\n")
	for _, param := range params {
		if param == "" {
			continue
		}
		for _, header := range ignoreHeaders {
			if strings.HasPrefix(param, header+":") { // is a header && need to ignore
				result = append(result, header+": *")
			} else {
				result = append(result, param)
			}
		}
	}

	return result
}

// =======
// binding
// =======

var (
	errNotV10Validator = errors.New("xgin: validator is not github.com/go-playground/validator/v10")
)

// GetValidatorEngine returns gin's binding validator engine, which only supports validator.Validate from github.com/go-playground/validator/v10.
func GetValidatorEngine() (*validator.Validate, error) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errNotV10Validator
	}
	return val, nil
}

// // GetValidatorTranslator applies and returns ut.Translator for validator.Validate using given locales.Translator and xvalidator.TranslationRegisterHandler.
// // Also see xvalidator.ApplyValidatorTranslator.
// //
// // Example:
// // 	translator, err = xgin.GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
// func GetValidatorTranslator(locTrans locales.Translator, registerFn xvalidator.TranslationRegisterHandler) (ut.Translator, error) {
// 	val, err := GetValidatorEngine()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return xvalidator.ApplyValidatorTranslator(val, locTrans, registerFn) // create translator and do register
// }

// AddBinding adds user defined binding to gin's validator engine. You can use your custom validator.Func or xvalidator's functions
// such as xvalidator.RegexpValidator and xvalidator.DateTimeValidator.
//
// Binding notes:
//
// 1. `required` + non-pointer (common)
// 	A uint64 `binding:"required"` // cannot be nil and 0
// 	B string `binding:"required"` // cannot be nil and ""
//
// 2. `required` + pointer (common)
// 	A *uint64 `binding:"required"` // cannot be nil, can be 0
// 	B *string `binding:"required"` // cannot be nil, can be ""
//
// 3. `omitempty` + non-pointer (common)
// 	A uint64 `binding:"omitempty"` // can be nil and 0
// 	B string `binding:"omitempty"` // can be nil and ""
//
// 4. `omitempty` + pointer => same as 3
// 	A *uint64 `binding:"omitempty"` // can be nil and 0
// 	B *string `binding:"omitempty"` // can be nil and ""
//
// 5. `required` + `omitempty` + non-pointer => same as 1
// 	A uint64 `binding:"required,omitempty"` // cannot be nil and 0
// 	B string `binding:"required,omitempty"` // cannot be nil and ""
//
// 6. `required` + `omitempty` + pointer => same as 2
// 	A *uint64 `binding:"required,omitempty"` // cannot be nil, can be 0
// 	B *string `binding:"required,omitempty"` // cannot be nil, can be ""
//
// Also see https://godoc.org/github.com/go-playground/validator.
func AddBinding(tag string, fn validator.Func) error {
	v, err := GetValidatorEngine()
	if err != nil {
		return nil
	}

	return v.RegisterValidation(tag, fn)
}

// // AddTranslator adds user defined validation translator to ut.Translator, which only supports validator.Validate. Also see
// // xvalidator.AddToTranslatorFunc and xvalidator.DefaultTranslationFunc.
// //
// // Example:
// // 	err = xgin.AddTranslator(translator, "regexp", "{0} must matches regexp /{1}/", true)
// // 	err = xgin.AddTranslator(translator, "email", "{0} must be an email", true)
// func AddTranslator(translator ut.Translator, tag, message string, override bool) error {
// 	v, err := GetValidatorEngine()
// 	if err != nil {
// 		return nil
// 	}
//
// 	fn := xvalidator.AddToTranslatorFunc(tag, message, override)
// 	return v.RegisterTranslation(tag, translator, fn, xvalidator.DefaultTranslationFunc())
// }

// // EnableRegexpBinding enables parametered regexp validator to `regexp`, see xvalidator.ParamRegexpValidator.
// func EnableRegexpBinding() error {
// 	return AddBinding("regexp", xvalidator.ParamRegexpValidator())
// }
//
// // EnableRegexpBindingTranslator enables parametered regexp validator's translator to given ut.Translator.
// func EnableRegexpBindingTranslator(translator ut.Translator) error {
// 	return AddTranslator(translator, "regexp", "{0} must matches regexp /{1}/", true)
// }
//
// // EnableRFC3339DateBinding enables rfc3339 date validator to `date`, see xvalidator.DateTimeValidator.
// func EnableRFC3339DateBinding() error {
// 	return AddBinding("date", xvalidator.DateTimeValidator(xtime.RFC3339Date))
// }
//
// // EnableRFC3339DateBindingTranslator enables rfc3339 date validator's translator to given ut.Translator.
// func EnableRFC3339DateBindingTranslator(translator ut.Translator) error {
// 	return AddTranslator(translator, "date", "{0} must be an RFC3339 date", true)
// }
//
// // EnableRFC3339DateTimeBinding enables rfc3339 datetime validator to `datetime`, see xvalidator.DateTimeValidator.
// func EnableRFC3339DateTimeBinding() error {
// 	return AddBinding("datetime", xvalidator.DateTimeValidator(xtime.RFC3339DateTime))
// }
//
// // EnableRFC3339DateTimeBindingTranslator enables rfc3339 datetime validator's translator to given ut.Translator.
// func EnableRFC3339DateTimeBindingTranslator(translator ut.Translator) error {
// 	return AddTranslator(translator, "datetime", "{0} must be an RFC3339 datetime", true)
// }

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
		{"GET", "/debug/pprof/", indexHandler()},
		{"GET", "/debug/pprof/heap", heapHandler()},
		{"GET", "/debug/pprof/goroutine", goroutineHandler()},
		{"GET", "/debug/pprof/allocs", allocsHandler()},
		{"GET", "/debug/pprof/block", blockHandler()},
		{"GET", "/debug/pprof/threadcreate", threadCreateHandler()},
		{"GET", "/debug/pprof/cmdline", cmdlineHandler()},
		{"GET", "/debug/pprof/profile", profileHandler()},
		{"GET", "/debug/pprof/symbol", symbolHandler()},
		{"POST", "/debug/pprof/symbol", symbolHandler()},
		{"GET", "/debug/pprof/trace", traceHandler()},
		{"GET", "/debug/pprof/mutex", mutexHandler()},
	} {
		router.Handle(r.method, r.path, r.handler) // use path directly
	}
}

// indexHandler is used for GET /debug/pprof to pprof.
func indexHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Index(ctx.Writer, ctx.Request)
	}
}

// heapHandler is for GET /debug/pprof/heap to pprof.
func heapHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("heap").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// goroutineHandler is used for GET /debug/pprof/goroutine to pprof.
func goroutineHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("goroutine").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// allocsHandler is used for GET /debug/pprof/allocs to pprof.
func allocsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("allocs").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// blockHandler is used for GET /debug/pprof/block to pprof.
func blockHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("block").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// threadCreateHandler is used for GET /debug/pprof/threadcreate to pprof.
func threadCreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// cmdlineHandler is used for GET /debug/pprof/cmdline to pprof.
func cmdlineHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Cmdline(ctx.Writer, ctx.Request)
	}
}

// profileHandler is used for GET /debug/pprof/profile to pprof.
func profileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Profile(ctx.Writer, ctx.Request)
	}
}

// symbolHandler is used for GET, POST /debug/pprof/symbol to pprof.
func symbolHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Symbol(ctx.Writer, ctx.Request)
	}
}

// traceHandler is used for GET /debug/pprof/trace to pprof.
func traceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Trace(ctx.Writer, ctx.Request)
	}
}

// mutexHandler is used for GET /debug/pprof/mutex to pprof.
func mutexHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("mutex").ServeHTTP(ctx.Writer, ctx.Request)
	}
}
