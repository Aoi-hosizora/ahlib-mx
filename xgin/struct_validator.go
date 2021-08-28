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

type TranslatableValidateFieldError = validator.FieldError

type CustomMessageValidateFieldError struct {
	origin  TranslatableValidateFieldError
	message string
}

func (v *CustomMessageValidateFieldError) Origin() TranslatableValidateFieldError {
	return v.origin
}

func (v *CustomMessageValidateFieldError) Error() string {
	return v.message
}

type FieldsValidateError struct {
	fields []error // TranslatableValidateFieldError or CustomMessageValidateFieldError
}

func (v *FieldsValidateError) Fields() []error {
	return v.fields
}

func (v *FieldsValidateError) Error() string {
	msgs := make([]string, 0, len(v.fields))
	for _, fe := range v.fields {
		msgs = append(msgs, fe.Error())
	}
	return strings.Join(msgs, "\n")
}

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

type CustomStructValidator struct {
	once       sync.Once
	validate   *validator.Validate
	messageTag string
}

var _ binding.StructValidator = (*CustomStructValidator)(nil)

func (v *CustomStructValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.messageTag = "validator_message"
		v.validate.SetTagName("binding")
	})
}

func (v *CustomStructValidator) SetValidatorTagName(name string) {
	v.validate.SetTagName(name)
}

func (v *CustomStructValidator) SetMessageTagName(name string) {
	v.messageTag = name
}

func (v *CustomStructValidator) ValidatorTagName() string {
	return xreflect.GetUnexportedField(reflect.ValueOf(v.validate).Elem().FieldByName("tagName")).Interface().(string)
}

func (v *CustomStructValidator) MessageTagName() string {
	return v.messageTag
}

func (v *CustomStructValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

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
		if ce, found := v.applyCustomMessage(typ, fe); !found {
			errs = append(errs, fe)
		} else {
			errs = append(errs, ce)
		}
	}
	return &FieldsValidateError{fields: errs}
}

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

func (v *CustomStructValidator) applyCustomMessage(typ reflect.Type, fe validator.FieldError) (error, bool) {
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

type RouterDecodeError struct {
	Field     string
	Input     string
	Err       error
	Translate string
}

func (r *RouterDecodeError) Error() string {
	if nErr, ok := r.Err.(*strconv.NumError); ok {
		return nErr.Error()
	}
	return fmt.Sprintf("parsing %s \"%s\": %v", r.Input, r.Input, r.Err)
}

func (r *RouterDecodeError) Unwrap() error {
	return r.Err
}

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

type TranslateOption func(*translateOptions)

func WithUtTranslator(translator ut.Translator) TranslateOption {
	return func(o *translateOptions) {
		o.utTranslator = translator
	}
}

func WithJsonInvalidUnmarshalError(fn func(*json.InvalidUnmarshalError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonInvalidUnmarshalErrorFn = fn
	}
}

func WithJsonUnmarshalTypeError(fn func(*json.UnmarshalTypeError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonUnmarshalTypeErrorFn = fn
	}
}

func WithJsonSyntaxError(fn func(*json.SyntaxError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.jsonSyntaxErrorFn = fn
	}
}

func WithIoEOFError(fn func(error) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.ioEOFErrorFn = fn
	}
}

func WithStrconvNumErrorError(fn func(*strconv.NumError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.strconvNumErrorFn = fn
	}
}

func WithXginRouterDecodeError(fn func(*RouterDecodeError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginRouterDecodeErrorFn = fn
	}
}

func WithValidatorInvalidTypeError(fn func(*validator.InvalidValidationError) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorInvalidTypeErrorFn = fn
	}
}

func WithValidatorFieldsError(fn func(validator.ValidationErrors, ut.Translator) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.validatorFieldsErrorFn = fn
	}
}

func WithXginFieldsValidateError(fn func(*FieldsValidateError, ut.Translator) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.xginFieldsValidateErrorFn = fn
	}
}

func WithExtraErrorsTranslate(fn func(error) (result map[string]string, isUserErr bool)) TranslateOption {
	return func(o *translateOptions) {
		o.extraErrorsTranslateFn = fn
	}
}

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
			reason := e.Translate
			if sErr, ok := e.Err.(*strconv.NumError); ok {
				if errors.Is(sErr.Err, strconv.ErrSyntax) {
					reason = "is not a number"
				} else if errors.Is(sErr.Err, strconv.ErrRange) {
					reason = "is out of range"
				}
			}
			if reason == "" {
				return nil, false
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
			return map[string]string{"_error": fmt.Sprintf("%T: %v", err, err)}, false
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
