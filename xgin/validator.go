package xgin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"strconv"
)

// ======================
// validator & translator
// ======================

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

// GetValidatorEnglishTranslator applies and returns English xvalidator.UtTranslator for validator.Validate using given parameters, this is a simplified usage of
// GetValidatorTranslator(validator, xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc()). Also see xvalidator.ApplyTranslator.
func GetValidatorEnglishTranslator() (xvalidator.UtTranslator, error) {
	return GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
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

// =================
// binding translate
// =================

// TranslatableError is an interface contains Translate method used in TranslateBindingError, can be used to specify the translation result of your error type.
type TranslatableError interface {
	Error() string
	Translate() (result map[string]string, need4xx bool)
}

// translateOptions is a type of TranslateBindingError's option, each field can be set by TranslateOption function type.
type translateOptions struct {
	utTranslator xvalidator.UtTranslator

	jsonInvalidUnmarshalErrorFn  func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)
	jsonUnmarshalTypeErrorFn     func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)
	jsonSyntaxErrorFn            func(*json.SyntaxError) (result map[string]string, need4xx bool)
	ioEOFErrorFn                 func(error) (result map[string]string, need4xx bool)
	strconvNumErrorFn            func(*strconv.NumError) (result map[string]string, need4xx bool)
	xginRouterDecodeErrorFn      func(*RouterDecodeError) (result map[string]string, need4xx bool)
	validatorInvalidTypeErrorFn  func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)
	validatorFieldsErrorFn       func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)
	xvalidatorMultiFieldsErrorFn func(*xvalidator.MultiFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)
	translatableErrorFn          func(TranslatableError) (result map[string]string, need4xx bool)
	extraErrorsTranslateFn       func(error) (result map[string]string, need4xx bool) // will never be nil
}

// TranslateOption represents an option for TranslateBindingError's options, can be created by WithXXX functions.
type TranslateOption func(*translateOptions)

// WithUtTranslator creates a TranslateOption to specify xvalidator.UtTranslator as the translation of validator.ValidationErrors and xvalidator.MultiFieldsError.
func WithUtTranslator(translator xvalidator.UtTranslator) TranslateOption {
	return func(o *translateOptions) {
		o.utTranslator = translator
	}
}

// WithJsonInvalidUnmarshalError creates a TranslateOption to specify translation function for json.InvalidUnmarshalError.
func WithJsonInvalidUnmarshalError(fn func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonInvalidUnmarshalErrorFn = fn
	}
}

// WithJsonUnmarshalTypeError creates a TranslateOption to specify translation function for json.UnmarshalTypeError.
func WithJsonUnmarshalTypeError(fn func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonUnmarshalTypeErrorFn = fn
	}
}

// WithJsonSyntaxError creates a TranslateOption to specify translation function for json.SyntaxError.
func WithJsonSyntaxError(fn func(*json.SyntaxError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonSyntaxErrorFn = fn
	}
}

// WithIoEOFError creates a TranslateOption to specify translation function for io.EOF and io.ErrUnexpectedEOF.
func WithIoEOFError(fn func(error) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.ioEOFErrorFn = fn
	}
}

// WithStrconvNumErrorError creates a TranslateOption to specify translation function for strconv.NumError.
func WithStrconvNumErrorError(fn func(*strconv.NumError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.strconvNumErrorFn = fn
	}
}

// WithXginRouterDecodeError creates a TranslateOption to specify translation function for xgin.RouterDecodeError.
func WithXginRouterDecodeError(fn func(*RouterDecodeError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginRouterDecodeErrorFn = fn
	}
}

// WithValidatorInvalidTypeError creates a TranslateOption to specify translation function for validator.InvalidValidationError.
func WithValidatorInvalidTypeError(fn func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorInvalidTypeErrorFn = fn
	}
}

// WithValidatorFieldsError creates a TranslateOption to specify translation function for validator.ValidationErrors.
func WithValidatorFieldsError(fn func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorFieldsErrorFn = fn
	}
}

// WithXvalidatorMultiFieldsError creates a TranslateOption to specify translation function for xvalidator.MultiFieldsError.
func WithXvalidatorMultiFieldsError(fn func(*xvalidator.MultiFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xvalidatorMultiFieldsErrorFn = fn
	}
}

// WithTranslatableError creates a TranslateOption to specify translation function for errors that implement xgin.TranslatableError interface.
func WithTranslatableError(fn func(TranslatableError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.translatableErrorFn = fn
	}
}

// WithExtraErrorsTranslate creates a TranslateOption to specify translation function for other errors.
func WithExtraErrorsTranslate(fn func(error) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.extraErrorsTranslateFn = fn
	}
}

// Default translation functions, used in TranslateBindingError.
var (
	_jsonInvalidUnmarshalErrorFn = func(e *json.InvalidUnmarshalError) (result map[string]string, need4xx bool) {
		return nil, false
	}
	_jsonUnmarshalTypeErrorFn = func(e *json.UnmarshalTypeError) (result map[string]string, need4xx bool) {
		return map[string]string{"__decode": fmt.Sprintf("type of '%s' in '%s' mismatches with required '%s'", e.Value, e.Field, e.Type.String())}, true
	}
	_jsonSyntaxErrorFn = func(e *json.SyntaxError) (result map[string]string, need4xx bool) {
		return map[string]string{"__decode": fmt.Sprintf("requested json has an invalid syntax at position %d", e.Offset)}, true
	}
	_ioEOFErrorFn = func(error) (result map[string]string, need4xx bool) {
		return map[string]string{"__decode": "requested json has an invalid syntax at position -1"}, true
	}
	_strconvNumErrorFn = func(e *strconv.NumError) (result map[string]string, need4xx bool) {
		reason := ""
		if errors.Is(e.Err, strconv.ErrSyntax) {
			reason = "must be a number"
		} else if errors.Is(e.Err, strconv.ErrRange) {
			reason = "is out of range"
		}
		if reason == "" {
			return nil, false // <<<
		}
		return map[string]string{"router parameter": fmt.Sprintf("router parameter %s", reason)}, true
	}
	_xginRouterDecodeErrorFn = func(e *RouterDecodeError) (result map[string]string, need4xx bool) {
		reason := e.Message
		if nErr, ok := e.Err.(*strconv.NumError); ok && reason == "" {
			if errors.Is(nErr.Err, strconv.ErrSyntax) {
				reason = "must be a number"
			} else if errors.Is(nErr.Err, strconv.ErrRange) {
				reason = "is out of range"
			}
		}
		if reason == "" {
			return nil, false // <<<
		}
		if e.Field == "" {
			return map[string]string{"router parameter": fmt.Sprintf("router parameter %s", reason)}, true
		}
		return map[string]string{e.Field: fmt.Sprintf("router parameter %s %s", e.Field, reason)}, true
	}
	_validatorInvalidTypeErrorFn = func(e *validator.InvalidValidationError) (result map[string]string, need4xx bool) {
		return nil, false
	}
	_validatorFieldsErrorFn = func(e validator.ValidationErrors, translator xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
		if translator == nil {
			return xvalidator.FlatValidationErrors(e, false), true
		}
		return xvalidator.TranslateValidationErrors(e, translator, false), true
	}
	_xvalidatorMultiFieldsErrorFn = func(e *xvalidator.MultiFieldsError, translator xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
		if translator == nil {
			return e.FlatToMap(false), true
		}
		return e.Translate(translator, false), true
	}
	_translatableErrorFn = func(e TranslatableError) (result map[string]string, need4xx bool) {
		return e.Translate()
	}
	_extraErrorsTranslateFn = func(err error) (result map[string]string, need4xx bool) {
		return nil, false // cannot be nil
	}
)

// TranslateBindingError translates given error and TranslateOption-s to a field-message map. Note that returned boolean value means whether given error
// can be regarded as an HTTP 4xx status code, means a user induced error.
func TranslateBindingError(err error, options ...TranslateOption) (result map[string]string, need4xx bool) {
	if err == nil {
		return nil, false
	}
	opt := &translateOptions{
		jsonInvalidUnmarshalErrorFn:  _jsonInvalidUnmarshalErrorFn,
		jsonUnmarshalTypeErrorFn:     _jsonUnmarshalTypeErrorFn,
		jsonSyntaxErrorFn:            _jsonSyntaxErrorFn,
		ioEOFErrorFn:                 _ioEOFErrorFn,
		strconvNumErrorFn:            _strconvNumErrorFn,
		xginRouterDecodeErrorFn:      _xginRouterDecodeErrorFn,
		validatorInvalidTypeErrorFn:  _validatorInvalidTypeErrorFn,
		validatorFieldsErrorFn:       _validatorFieldsErrorFn,
		xvalidatorMultiFieldsErrorFn: _xvalidatorMultiFieldsErrorFn,
		translatableErrorFn:          _translatableErrorFn,
		extraErrorsTranslateFn:       _extraErrorsTranslateFn,
	}
	for _, op := range options {
		if op != nil {
			op(opt)
		}
	}
	if opt.extraErrorsTranslateFn == nil {
		opt.extraErrorsTranslateFn = _extraErrorsTranslateFn
	}

	// 1. body
	if fn := opt.jsonInvalidUnmarshalErrorFn; fn != nil {
		if jErr, ok := err.(*json.InvalidUnmarshalError); ok {
			return fn(jErr)
		}
	}
	if fn := opt.jsonUnmarshalTypeErrorFn; fn != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return fn(jErr)
		}
	}
	if fn := opt.jsonSyntaxErrorFn; fn != nil {
		if jErr, ok := err.(*json.SyntaxError); ok {
			return fn(jErr)
		}
	}
	if fn := opt.ioEOFErrorFn; fn != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return fn(err)
		}
	}

	// 2. router
	if fn := opt.strconvNumErrorFn; fn != nil {
		if sErr, ok := err.(*strconv.NumError); ok {
			return fn(sErr)
		}
	}
	if fn := opt.xginRouterDecodeErrorFn; fn != nil {
		if rErr, ok := err.(*RouterDecodeError); ok {
			return fn(rErr)
		}
	}

	// 3. validate
	if fn := opt.validatorInvalidTypeErrorFn; fn != nil {
		if vErr, ok := err.(*validator.InvalidValidationError); ok {
			return fn(vErr)
		}
	}
	if fn := opt.validatorFieldsErrorFn; fn != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			return fn(vErr, opt.utTranslator)
		}
	}
	if fn := opt.xvalidatorMultiFieldsErrorFn; fn != nil {
		if vErr, ok := err.(*xvalidator.MultiFieldsError); ok {
			return fn(vErr, opt.utTranslator)
		}
	}

	// 4. extra
	if fn := opt.translatableErrorFn; fn != nil {
		if tErr, ok := err.(TranslatableError); ok {
			return fn(tErr)
		}
	}
	return opt.extraErrorsTranslateFn(err)
}
