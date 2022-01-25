package xgin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/go-playground/validator/v10"
	"io"
	"strconv"
)

// TranslatableError is an interface contains Translate method used in TranslateBindingError, can be used to specific the translation result of your error type.
type TranslatableError interface {
	error
	Translate() (result map[string]string, need4xx bool)
}

// translateOptions is a type of TranslateBindingError's option, each field can be set by TranslateOption function type.
type translateOptions struct {
	utTranslator xvalidator.UtTranslator

	jsonInvalidUnmarshalErrorFn     func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)
	jsonUnmarshalTypeErrorFn        func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)
	jsonSyntaxErrorFn               func(*json.SyntaxError) (result map[string]string, need4xx bool)
	ioEOFErrorFn                    func(error) (result map[string]string, need4xx bool)
	strconvNumErrorFn               func(*strconv.NumError) (result map[string]string, need4xx bool)
	xginRouterDecodeErrorFn         func(*RouterDecodeError) (result map[string]string, need4xx bool)
	validatorInvalidTypeErrorFn     func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)
	validatorFieldsErrorFn          func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)
	translatableErrorFn             func(TranslatableError) (result map[string]string, need4xx bool)
	xvalidatorValidateFieldsErrorFn func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)

	extraErrorsTranslateFn func(error) (result map[string]string, need4xx bool)
}

// TranslateOption represents an option for TranslateBindingError's options, can be created by WithXXX functions.
type TranslateOption func(*translateOptions)

// WithUtTranslator creates a TranslateOption to specific xvalidator.UtTranslator as the translation of validator.ValidationErrors and xvalidator.ValidateFieldsError.
func WithUtTranslator(translator xvalidator.UtTranslator) TranslateOption {
	return func(o *translateOptions) {
		o.utTranslator = translator
	}
}

// WithJsonInvalidUnmarshalError creates a TranslateOption to specific translation function for json.InvalidUnmarshalError.
func WithJsonInvalidUnmarshalError(fn func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonInvalidUnmarshalErrorFn = fn
	}
}

// WithJsonUnmarshalTypeError creates a TranslateOption to specific translation function for json.UnmarshalTypeError.
func WithJsonUnmarshalTypeError(fn func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonUnmarshalTypeErrorFn = fn
	}
}

// WithJsonSyntaxError creates a TranslateOption to specific translation function for json.SyntaxError.
func WithJsonSyntaxError(fn func(*json.SyntaxError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonSyntaxErrorFn = fn
	}
}

// WithIoEOFError creates a TranslateOption to specific translation function for io.EOF and io.ErrUnexpectedEOF.
func WithIoEOFError(fn func(error) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.ioEOFErrorFn = fn
	}
}

// WithStrconvNumErrorError creates a TranslateOption to specific translation function for strconv.NumError.
func WithStrconvNumErrorError(fn func(*strconv.NumError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.strconvNumErrorFn = fn
	}
}

// WithXginRouterDecodeError creates a TranslateOption to specific translation function for xgin.RouterDecodeError.
func WithXginRouterDecodeError(fn func(*RouterDecodeError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginRouterDecodeErrorFn = fn
	}
}

// WithValidatorInvalidTypeError creates a TranslateOption to specific translation function for validator.InvalidValidationError.
func WithValidatorInvalidTypeError(fn func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorInvalidTypeErrorFn = fn
	}
}

// WithValidatorFieldsError creates a TranslateOption to specific translation function for validator.ValidationErrors.
func WithValidatorFieldsError(fn func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorFieldsErrorFn = fn
	}
}

// WithXvalidatorValidateFieldsError creates a TranslateOption to specific translation function for xvalidator.ValidateFieldsError.
func WithXvalidatorValidateFieldsError(fn func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xvalidatorValidateFieldsErrorFn = fn
	}
}

// WithTranslatableError creates a TranslateOption to specific translation function for errors that implement xgin.TranslatableError interface.
func WithTranslatableError(fn func(TranslatableError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.translatableErrorFn = fn
	}
}

// WithExtraErrorsTranslate creates a TranslateOption to specific translation function for other errors.
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
		if reason != "" {
			return map[string]string{"router parameter": fmt.Sprintf("router parameter %s", reason)}, true
		}
		return nil, false
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
			return xvalidator.FlatValidateErrors(e, false), true
		}
		return xvalidator.TranslateValidationErrors(e, translator, false), true
	}
	_xvalidatorValidateFieldsErrorFn = func(e *xvalidator.ValidateFieldsError, translator xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
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

// TranslateBindingError translates the given error and TranslateOption-s to a field-message map. Note that the returned boolean value means whether the given
// error can be regarded as an HTTP 4xx status code, means a user induced error.
func TranslateBindingError(err error, options ...TranslateOption) (result map[string]string, need4xx bool) {
	if err == nil {
		return nil, false
	}
	opt := &translateOptions{
		jsonInvalidUnmarshalErrorFn:     _jsonInvalidUnmarshalErrorFn,
		jsonUnmarshalTypeErrorFn:        _jsonUnmarshalTypeErrorFn,
		jsonSyntaxErrorFn:               _jsonSyntaxErrorFn,
		ioEOFErrorFn:                    _ioEOFErrorFn,
		strconvNumErrorFn:               _strconvNumErrorFn,
		xginRouterDecodeErrorFn:         _xginRouterDecodeErrorFn,
		validatorInvalidTypeErrorFn:     _validatorInvalidTypeErrorFn,
		validatorFieldsErrorFn:          _validatorFieldsErrorFn,
		xvalidatorValidateFieldsErrorFn: _xvalidatorValidateFieldsErrorFn,
		translatableErrorFn:             _translatableErrorFn,
		extraErrorsTranslateFn:          _extraErrorsTranslateFn,
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
	if fn := opt.xvalidatorValidateFieldsErrorFn; fn != nil {
		if vErr, ok := err.(*xvalidator.ValidateFieldsError); ok {
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
