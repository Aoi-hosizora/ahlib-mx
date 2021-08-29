package xgin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// TranslatableValidateFieldError represents an alias type of validator.FieldError interface, and this is the error for normal single field's validator error.
type TranslatableValidateFieldError = validator.FieldError

// CustomMessageValidateFieldError represents a message wrapped TranslatableValidateFieldError, and this is the error for those fields which can have specific message.
type CustomMessageValidateFieldError struct {
	origin  TranslatableValidateFieldError
	message string
}

// Origin returns the origin TranslatableValidateFieldError from CustomMessageValidateFieldError.
func (v *CustomMessageValidateFieldError) Origin() TranslatableValidateFieldError {
	return v.origin
}

// Error returns the error message from CustomMessageValidateFieldError.
func (v *CustomMessageValidateFieldError) Error() string {
	return v.message
}

// FieldsValidateError represents the struct fields' validator errors slice, and this error may be returned by CustomStructValidator.ValidateStruct.
type FieldsValidateError struct {
	fields []error // TranslatableValidateFieldError or CustomMessageValidateFieldError
}

// Fields returns the fields' errors from FieldsValidateError.
func (v *FieldsValidateError) Fields() []error {
	return v.fields
}

// Error returns the error message from FieldsValidateError.
func (v *FieldsValidateError) Error() string {
	msgs := make([]string, 0, len(v.fields))
	for _, fe := range v.fields {
		msgs = append(msgs, fe.Error())
	}
	return strings.Join(msgs, "\n")
}

// Translate translates FieldsValidateError to a field-message map, using given ut.Translator and two kinds of error.
func (v *FieldsValidateError) Translate(translator ut.Translator) map[string]string {
	result := make(map[string]string, len(v.fields))
	for _, err := range v.fields {
		if te, ok := err.(TranslatableValidateFieldError); ok {
			result[te.Field()] = te.Translate(translator)
		} else if ce, ok := err.(*CustomMessageValidateFieldError); ok {
			result[ce.origin.Field()] = ce.message
		}
	}
	return result
}

// CustomStructValidator represents a custom validator.Validate as binding.StructValidator, which allows fields to specify their custom error messages.
//
// Example:
// 	type User struct {
//		Id   uint64  `json:"id"   form:"id"   binding:"required,gt=1"          validator_message:"required|id is required|gt|id must larger than one"`
//		Name string  `json:"name" form:"name" binding:"required,gt=4,lt=20"    validator_message:"*|name is invalid"`
//		Bio  *string `json:"bio"  form:"bio"  binding:"required,gte=0,lte=255" validator_message:"xxx"`
// 	}
type CustomStructValidator struct {
	once       sync.Once
	validate   *validator.Validate
	messageTag string
}

var _ binding.StructValidator = (*CustomStructValidator)(nil)

// lazyinit initializes CustomStructValidator.validate in lazy.
func (v *CustomStructValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		v.messageTag = "validator_message"
	})
}

// SetValidatorTagName sets the validator tag name for CustomStructValidator.
func (v *CustomStructValidator) SetValidatorTagName(name string) {
	v.validate.SetTagName(name)
}

// SetMessageTagName sets the message tag name for CustomStructValidator.
func (v *CustomStructValidator) SetMessageTagName(name string) {
	v.messageTag = name
}

// ValidatorTagName returns the validator tag name of CustomStructValidator.
func (v *CustomStructValidator) ValidatorTagName() string {
	val := reflect.ValueOf(v.validate).Elem()
	return xreflect.GetUnexportedField(val.FieldByName("tagName")).Interface().(string)
}

// MessageTagName returns the validator tag name of CustomStructValidator.
func (v *CustomStructValidator) MessageTagName() string {
	return v.messageTag
}

// Engine returns the internal validator.Validate engine from CustomStructValidator.
func (v *CustomStructValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

// ValidateStruct validates the given struct and returns the error, the type of error can be validator.InvalidValidationError or xgin.FieldsValidateError.
func (v *CustomStructValidator) ValidateStruct(obj interface{}) error {
	val, ok := v.extractToStruct(obj)
	if !ok {
		return &validator.InvalidValidationError{Type: reflect.TypeOf(obj)}
	}

	v.lazyinit()
	err := v.validate.Struct(val) // InvalidValidationError or ValidationErrors
	if err == nil {
		return nil
	}
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	typ := reflect.TypeOf(val)
	errs := make([]error, 0, len(ve))
	for _, fe := range ve {
		if ce, found := v.applyCustomMessage(typ, fe); found {
			errs = append(errs, ce) // TranslatableValidateFieldError -> CustomMessageValidateFieldError
		} else {
			errs = append(errs, fe) // TranslatableValidateFieldError (validator.FieldError)
		}
	}
	return &FieldsValidateError{fields: errs}
}

// extractToStruct checks and extracts the given interface to struct type.
func (v *CustomStructValidator) extractToStruct(obj interface{}) (interface{}, bool) {
	if obj == nil {
		return nil, false
	}
	val := reflect.ValueOf(obj)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, false
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, false
	}
	return val.Interface(), true
}

// applyCustomMessage checks the struct field and wraps TranslatableValidateFieldError to CustomMessageValidateFieldError. Note that "||" is used to
// represent a single "|", such as "*|name ||is|| invalid" means the message for "*" is set to "name |is| invalid".
func (v *CustomStructValidator) applyCustomMessage(typ reflect.Type, fe TranslatableValidateFieldError) (error, bool) {
	sf, ok := typ.FieldByName(fe.StructField())
	if !ok {
		return nil, false // unreachable
	}
	msg := sf.Tag.Get("validator_message")
	if msg == "" {
		return nil, false // no msg
	}
	msg = strings.ReplaceAll(msg, "||", "｜")
	sp := strings.Split(msg, "|")
	for i := 0; i < len(sp); i += 2 {
		if i+1 >= len(sp) {
			break
		}
		k, m := strings.TrimSpace(sp[i]), strings.TrimSpace(sp[i+1])
		k, m = strings.ReplaceAll(k, "｜", "|"), strings.ReplaceAll(m, "｜", "|")
		if k == "*" || k == fe.Tag() {
			return &CustomMessageValidateFieldError{origin: fe, message: m}, true // found
		}
	}
	return nil, false // not found
}

// RouterDecodeError represents an error type for router parameter decoding, and this error can specify custom translate message and error string.
// At most of the time, Err field is in strconv.NumError type generated by strconv.ParseInt, etc.
type RouterDecodeError struct {
	Field       string
	Input       string
	Err         error
	Translation string

	CustomErrorFn func(error) string
}

// Error returns the error message from RouterDecodeError.
func (r *RouterDecodeError) Error() string {
	if r.CustomErrorFn != nil {
		return r.CustomErrorFn(r.Err)
	}
	if nErr, ok := r.Err.(*strconv.NumError); ok {
		return nErr.Error()
	}
	return fmt.Sprintf("parsing %s \"%s\": %v", r.Input, r.Input, r.Err)
}

// Unwrap returns the wrapped error from RouterDecodeError.
func (r *RouterDecodeError) Unwrap() error {
	return r.Err
}

// translateOptions represents some TranslateBindingError function options, set by TranslateOption.
type translateOptions struct {
	utTranslator                ut.Translator
	jsonInvalidUnmarshalErrorFn func(*json.InvalidUnmarshalError) (result map[string]string, isUserErr bool)
	jsonUnmarshalTypeErrorFn    func(*json.UnmarshalTypeError) (result map[string]string, isUserErr bool)
	jsonSyntaxErrorFn           func(*json.SyntaxError) (result map[string]string, isUserErr bool)
	ioEOFErrorFn                func(error) (result map[string]string, isUserErr bool)
	strconvNumErrorFn           func(*strconv.NumError) (result map[string]string, isUserErr bool)
	xginRouterDecodeErrorFn     func(*RouterDecodeError) (result map[string]string, isUserErr bool)
	validatorInvalidTypeErrorFn func(*validator.InvalidValidationError) (result map[string]string, isUserErr bool)
	validatorFieldsErrorFn      func(validator.ValidationErrors, ut.Translator) (result map[string]string, isUserErr bool)
	xginFieldsValidateErrorFn   func(*FieldsValidateError, ut.Translator) (result map[string]string, isUserErr bool)
	extraErrorsTranslateFn      func(error) (result map[string]string, isUserErr bool)
}

// TranslateOption represents an option for translateOptions, created by WithXXX functions.
type TranslateOption func(*translateOptions)

// WithUtTranslator creates a TranslateOption to specific ut.Translator as the translation of validator.ValidationErrors and xgin.FieldsValidateError.
func WithUtTranslator(translator ut.Translator) TranslateOption {
	return func(o *translateOptions) {
		o.utTranslator = translator
	}
}

// WithJsonInvalidUnmarshalError creates a translation function for json.InvalidUnmarshalError.
func WithJsonInvalidUnmarshalError(fn func(*json.InvalidUnmarshalError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonInvalidUnmarshalErrorFn = fn
	}
}

// WithJsonUnmarshalTypeError creates a translation function for json.UnmarshalTypeError.
func WithJsonUnmarshalTypeError(fn func(*json.UnmarshalTypeError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonUnmarshalTypeErrorFn = fn
	}
}

// WithJsonSyntaxError creates a translation function for json.SyntaxError.
func WithJsonSyntaxError(fn func(*json.SyntaxError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonSyntaxErrorFn = fn
	}
}

// WithIOEOFError creates a translation function for io.EOF and io.ErrUnexpectedEOF.
func WithIOEOFError(fn func(error) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.ioEOFErrorFn = fn
	}
}

// WithStrconvNumErrorError creates a translation function for strconv.NumError.
func WithStrconvNumErrorError(fn func(*strconv.NumError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.strconvNumErrorFn = fn
	}
}

// WithXginRouterDecodeError creates a translation function for xgin.RouterDecodeError.
func WithXginRouterDecodeError(fn func(*RouterDecodeError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginRouterDecodeErrorFn = fn
	}
}

// WithValidatorInvalidTypeError creates a translation function for validator.InvalidValidationError.
func WithValidatorInvalidTypeError(fn func(*validator.InvalidValidationError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorInvalidTypeErrorFn = fn
	}
}

// WithValidatorFieldsError creates a translation function for validator.ValidationErrors.
func WithValidatorFieldsError(fn func(validator.ValidationErrors, ut.Translator) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorFieldsErrorFn = fn
	}
}

// WithXginFieldsValidateError creates a translation function for xgin.FieldsValidateError.
func WithXginFieldsValidateError(fn func(*FieldsValidateError, ut.Translator) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginFieldsValidateErrorFn = fn
	}
}

// WithExtraErrorsTranslate creates a translation function for other errors, and it is the default translation function for extra error types.
func WithExtraErrorsTranslate(fn func(error) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.extraErrorsTranslateFn = fn
	}
}

// TranslateBindingError translates the given error and TranslateOption-s to a field-message map, note that the returned boolean value means that if the given
// error can be regarded as an HTTP 4xx status code (user induced error).
func TranslateBindingError(err error, options ...TranslateOption) (result map[string]string, isUserErr bool) {
	if err == nil {
		return nil, false
	}

	opt := &translateOptions{
		jsonInvalidUnmarshalErrorFn: func(e *json.InvalidUnmarshalError) (result map[string]string, isUserErr bool) {
			return nil, false
		},
		jsonUnmarshalTypeErrorFn: func(e *json.UnmarshalTypeError) (result map[string]string, isUserErr bool) {
			return map[string]string{"_decode": fmt.Sprintf("type of '%s' in '%s' mismatches with required '%s'", e.Value, e.Field, e.Type.String())}, true
		},
		jsonSyntaxErrorFn: func(e *json.SyntaxError) (result map[string]string, isUserErr bool) {
			return map[string]string{"_decode": fmt.Sprintf("requested json has an invalid syntax at position %d", e.Offset)}, true
		},
		ioEOFErrorFn: func(error) (result map[string]string, isUserErr bool) {
			return map[string]string{"_decode": "requested json has an invalid syntax at position -1"}, true
		},
		strconvNumErrorFn: func(e *strconv.NumError) (result map[string]string, isUserErr bool) {
			if errors.Is(e.Err, strconv.ErrSyntax) {
				return map[string]string{"router parameter": "router parameter is not a number"}, true
			}
			if errors.Is(e.Err, strconv.ErrRange) {
				return map[string]string{"router parameter": "router parameter is out of range"}, true
			}
			return nil, false
		},
		xginRouterDecodeErrorFn: func(e *RouterDecodeError) (result map[string]string, isUserErr bool) {
			reason := e.Translation
			if sErr, ok := e.Err.(*strconv.NumError); ok && reason == "" {
				if errors.Is(sErr.Err, strconv.ErrSyntax) {
					reason = "is not a number"
				} else if errors.Is(sErr.Err, strconv.ErrRange) {
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
		validatorInvalidTypeErrorFn: func(e *validator.InvalidValidationError) (result map[string]string, isUserErr bool) {
			return nil, false
		},
		validatorFieldsErrorFn: func(e validator.ValidationErrors, trans ut.Translator) (result map[string]string, isUserErr bool) {
			result = make(map[string]string, len(e))
			for _, fe := range e {
				result[fe.Field()] = fe.Translate(trans)
			}
			return result, true
		},
		xginFieldsValidateErrorFn: func(e *FieldsValidateError, trans ut.Translator) (result map[string]string, isUserErr bool) {
			return e.Translate(trans), true
		},
		extraErrorsTranslateFn: func(e error) (result map[string]string, isUserErr bool) {
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
	if trans := opt.utTranslator; trans != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			return opt.validatorFieldsErrorFn(vErr, trans)
		}
		if vErr, ok := err.(*FieldsValidateError); ok {
			return opt.xginFieldsValidateErrorFn(vErr, trans)
		}
	}

	// 4. extra & default
	return opt.extraErrorsTranslateFn(err)
}
