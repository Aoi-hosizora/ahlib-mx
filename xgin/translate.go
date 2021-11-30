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

// translateOptions is a type of TranslateBindingError function's option, each field can be set by TranslateOption.
type translateOptions struct {
	translator xvalidator.UtTranslator

	jsonInvalidUnmarshalErrorFn func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)
	jsonUnmarshalTypeErrorFn    func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)
	jsonSyntaxErrorFn           func(*json.SyntaxError) (result map[string]string, need4xx bool)
	ioEOFErrorFn                func(error) (result map[string]string, need4xx bool)
	strconvNumErrorFn           func(*strconv.NumError) (result map[string]string, need4xx bool)
	xginRouterDecodeErrorFn     func(*RouterDecodeError) (result map[string]string, need4xx bool)
	validatorInvalidTypeErrorFn func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)
	validatorFieldsErrorFn      func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)
	xginFieldsValidateErrorFn   func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)
	extraErrorsTranslateFn      func(error) (result map[string]string, need4xx bool)
}

// TranslateOption represents an option for translateOptions, can be created by WithXXX functions.
type TranslateOption func(*translateOptions)

// WithTranslator creates a TranslateOption to specific UtTranslator as the translation of validator.ValidationErrors and xgin.ValidateFieldsError.
func WithTranslator(translator xvalidator.UtTranslator) TranslateOption {
	return func(o *translateOptions) {
		o.translator = translator
	}
}

// WithJsonInvalidUnmarshalError creates a translation function for json.InvalidUnmarshalError.
func WithJsonInvalidUnmarshalError(fn func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonInvalidUnmarshalErrorFn = fn
	}
}

// WithJsonUnmarshalTypeError creates a translation function for json.UnmarshalTypeError.
func WithJsonUnmarshalTypeError(fn func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonUnmarshalTypeErrorFn = fn
	}
}

// WithJsonSyntaxError creates a translation function for json.SyntaxError.
func WithJsonSyntaxError(fn func(*json.SyntaxError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonSyntaxErrorFn = fn
	}
}

// WithIoEOFError creates a translation function for io.EOF and io.ErrUnexpectedEOF.
func WithIoEOFError(fn func(error) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.ioEOFErrorFn = fn
	}
}

// WithStrconvNumErrorError creates a translation function for strconv.NumError.
func WithStrconvNumErrorError(fn func(*strconv.NumError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.strconvNumErrorFn = fn
	}
}

// WithXginRouterDecodeError creates a translation function for xgin.RouterDecodeError.
func WithXginRouterDecodeError(fn func(*RouterDecodeError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginRouterDecodeErrorFn = fn
	}
}

// WithValidatorInvalidTypeError creates a translation function for validator.InvalidValidationError.
func WithValidatorInvalidTypeError(fn func(*validator.InvalidValidationError) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorInvalidTypeErrorFn = fn
	}
}

// WithValidatorFieldsError creates a translation function for validator.ValidationErrors.
func WithValidatorFieldsError(fn func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorFieldsErrorFn = fn
	}
}

// WithXginFieldsValidateError creates a translation function for xgin.ValidateFieldsError.
func WithXginFieldsValidateError(fn func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginFieldsValidateErrorFn = fn
	}
}

// WithExtraErrorsTranslate creates a translation function for other errors, and it is the default translation function for extra error types.
func WithExtraErrorsTranslate(fn func(error) (result map[string]string, need4xx bool)) TranslateOption {
	return func(o *translateOptions) {
		o.extraErrorsTranslateFn = fn
	}
}

// TranslateBindingError translates the given error and TranslateOption-s to a field-message map, note that the returned boolean value means whether the given
// error can be regarded as an HTTP 4xx status code (user induced error).
func TranslateBindingError(err error, options ...TranslateOption) (result map[string]string, need4xx bool) {
	if err == nil {
		return nil, false
	}

	opt := &translateOptions{
		jsonInvalidUnmarshalErrorFn: func(e *json.InvalidUnmarshalError) (result map[string]string, need4xx bool) {
			return nil, false
		},
		jsonUnmarshalTypeErrorFn: func(e *json.UnmarshalTypeError) (result map[string]string, need4xx bool) {
			return map[string]string{"__decode": fmt.Sprintf("type of '%s' in '%s' mismatches with required '%s'", e.Value, e.Field, e.Type.String())}, true
		},
		jsonSyntaxErrorFn: func(e *json.SyntaxError) (result map[string]string, need4xx bool) {
			return map[string]string{"__decode": fmt.Sprintf("requested json has an invalid syntax at position %d", e.Offset)}, true
		},
		ioEOFErrorFn: func(error) (result map[string]string, need4xx bool) {
			return map[string]string{"__decode": "requested json has an invalid syntax at position -1"}, true
		},
		strconvNumErrorFn: func(e *strconv.NumError) (result map[string]string, need4xx bool) {
			reason := ""
			if errors.Is(e.Err, strconv.ErrSyntax) {
				reason = "is not a number"
			} else if errors.Is(e.Err, strconv.ErrRange) {
				reason = "is out of range"
			}
			if reason != "" {
				return map[string]string{"router parameter": fmt.Sprintf("router parameter %s", reason)}, true
			}
			return nil, false
		},
		xginRouterDecodeErrorFn: func(e *RouterDecodeError) (result map[string]string, need4xx bool) {
			reason := e.Translation
			if nErr, ok := e.Err.(*strconv.NumError); ok && reason == "" {
				if errors.Is(nErr.Err, strconv.ErrSyntax) {
					reason = "is not a number"
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
		},
		validatorInvalidTypeErrorFn: func(e *validator.InvalidValidationError) (result map[string]string, need4xx bool) {
			return nil, false
		},
		validatorFieldsErrorFn: func(e validator.ValidationErrors, trans xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
			return xvalidator.TranslateValidationErrors(e, trans, false), true
		},
		xginFieldsValidateErrorFn: func(e *xvalidator.ValidateFieldsError, trans xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
			return e.Translate(trans, false), true
		},
		extraErrorsTranslateFn: func(e error) (result map[string]string, need4xx bool) {
			return nil, false
		},
	}
	for _, op := range options {
		if op != nil {
			op(opt)
		}
	}

	// 1. body
	if jErr, ok := err.(*json.InvalidUnmarshalError); ok {
		return opt.jsonInvalidUnmarshalErrorFn(jErr)
	}
	if jErr, ok := err.(*json.UnmarshalTypeError); ok {
		return opt.jsonUnmarshalTypeErrorFn(jErr)
	}
	if jErr, ok := err.(*json.SyntaxError); ok {
		return opt.jsonSyntaxErrorFn(jErr)
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return opt.ioEOFErrorFn(err)
	}

	// 2. router
	if sErr, ok := err.(*strconv.NumError); ok {
		return opt.strconvNumErrorFn(sErr)
	}
	if rErr, ok := err.(*RouterDecodeError); ok {
		return opt.xginRouterDecodeErrorFn(rErr)
	}

	// 3. validate
	if vErr, ok := err.(*validator.InvalidValidationError); ok {
		return opt.validatorInvalidTypeErrorFn(vErr)
	}
	if opt.translator != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			return opt.validatorFieldsErrorFn(vErr, opt.translator)
		}
		if vErr, ok := err.(*xvalidator.ValidateFieldsError); ok {
			return opt.xginFieldsValidateErrorFn(vErr, opt.translator)
		}
	}

	// 4. extra & default
	return opt.extraErrorsTranslateFn(err)
}
