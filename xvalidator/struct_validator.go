package xvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

// ===========================
// ValidateFieldsError related
// ===========================

// WrappedValidateFieldError represents a validator.FieldError with a custom message wrapped used in CustomStructValidator.
type WrappedValidateFieldError struct {
	origin  validator.FieldError
	message string
}

// Origin returns the origin validator.FieldError from WrappedValidateFieldError.
func (v *WrappedValidateFieldError) Origin() validator.FieldError {
	return v.origin
}

// Message returns the wrapped message from WrappedValidateFieldError.
func (v *WrappedValidateFieldError) Message() string {
	return v.message
}

// Error returns the formatted error message from WrappedValidateFieldError, and has the same format with validator.FieldError.Error.
func (v *WrappedValidateFieldError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:%s", v.origin.Namespace(), v.message)
}

// ValidateFieldsError represents the struct fields' validator errors, this error value will be returned by CustomStructValidator.ValidateStruct.
type ValidateFieldsError struct {
	fields []error // validator.FieldError or WrappedValidateFieldError
}

// Fields returns the fields' errors from ValidateFieldsError.
func (v *ValidateFieldsError) Fields() []error {
	return v.fields
}

// Error returns the total formatted error message from ValidateFieldsError.
func (v *ValidateFieldsError) Error() string {
	msgs := make([]string, 0, len(v.fields))
	for _, fe := range v.fields {
		if fe != nil {
			msgs = append(msgs, fe.Error())
		}
	}
	return strings.Join(msgs, "\n")
}

// Translate translates all the field errors using given UtTranslator to a field-message map, and errors of type WrappedValidateFieldError will use
// wrapped message directly. Note that if you set useNamespace flag to true, then the keys from returned map will show in "$struct.$field" format,
// otherwise in "$field" format.
//
// Example:
// 	v.Translate(trans, true)  // => map[s.int:int is a required field, s.str:str cannot be null and empty] (s.int uses validator.FieldError, s.str uses WrappedValidateFieldError)
// 	v.Translate(trans, false) // => map[int:int is a required field, str:str cannot be null and empty]
func (v *ValidateFieldsError) Translate(ut UtTranslator, useNamespace bool) map[string]string {
	if ut == nil {
		panic(panicNilUtTranslator)
	}
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(v.fields))
	for _, err := range v.fields {
		if fe, ok := err.(validator.FieldError); ok {
			result[keyFn(fe)] = fe.Translate(ut)
		} else if we, ok := err.(*WrappedValidateFieldError); ok {
			result[keyFn(we.origin)] = we.message
		} else {
			// skip
		}
	}
	return result
}

// =============================
// CustomStructValidator related
// =============================

// CustomStructValidator represents a custom validator.Validate, which allows some fields to specify their custom error message, and this can set to
// gin's binding.Validator as a binding.StructValidator.
//
// Example:
// 	type User struct {
//		Id   uint64  `json:"id"   form:"id"   binding:"required,gt=1"          validator_message:"required|id is required|gt|id must larger than one"`
//		Name string  `json:"name" form:"name" binding:"required,gt=4,lt=20"    validator_message:"*|name is invalid"`
//		Bio  *string `json:"bio"  form:"bio"  binding:"required,gte=0,lte=255" validator_message:"xxx"`
// 	}
type CustomStructValidator struct {
	validate   *validator.Validate
	messageTag string

	once sync.Once
}

// NewCustomStructValidator creates a new NewCustomStructValidator.
func NewCustomStructValidator() *CustomStructValidator {
	v := &CustomStructValidator{
		validate:   validator.New(),
		messageTag: "validator_message",
	}
	v.validate.SetTagName("binding")
	return v
}

// Engine returns the internal validator.Validate from CustomStructValidator.
func (v *CustomStructValidator) Engine() interface{} {
	return v.validate
}

// ValidateEngine returns the internal validator.Validate from CustomStructValidator.
func (v *CustomStructValidator) ValidateEngine() *validator.Validate {
	return v.validate
}

// SetValidatorTagName sets the validator tag name for CustomStructValidator.
func (v *CustomStructValidator) SetValidatorTagName(name string) {
	v.validate.SetTagName(name)
}

// SetMessageTagName sets the message tag name for CustomStructValidator.
func (v *CustomStructValidator) SetMessageTagName(name string) {
	v.messageTag = name
}

// ValidateStruct validates the given struct and returns the validator error, mostly in xvalidator.ValidateFieldsError type.
func (v *CustomStructValidator) ValidateStruct(obj interface{}) error {
	val, ok := v.extractToStruct(obj)
	if !ok {
		return &validator.InvalidValidationError{Type: reflect.TypeOf(obj)}
	}

	err := v.validate.Struct(val)
	if err == nil {
		return nil
	}
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return err // unreachable
	}

	typ := reflect.TypeOf(val)
	errs := make([]error, 0, len(ve))
	for _, fe := range ve {
		if m, found := v.applyCustomMessage(typ, fe.StructField(), fe.Tag()); found {
			we := &WrappedValidateFieldError{origin: fe, message: m}
			errs = append(errs, we) // WrappedValidateFieldError
		} else {
			errs = append(errs, fe) // validator.FieldError
		}
	}
	return &ValidateFieldsError{fields: errs}
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

// applyCustomMessage checks the struct field and wraps validator.FieldError to WrappedValidateFieldError. Note that "||" is used to represent
// a single "|", such as "*|name ||is|| invalid" means the message for "*" is set to "name |is| invalid".
func (v *CustomStructValidator) applyCustomMessage(typ reflect.Type, fieldName, validateTag string) (string, bool) {
	sf, ok := typ.FieldByName(fieldName)
	if !ok {
		return "", false // unreachable
	}
	msg := sf.Tag.Get(v.messageTag)
	if msg == "" {
		return "", false // no msg
	}
	msg = strings.ReplaceAll(msg, "\\|", "｜")
	sp := strings.Split(msg, "|")
	for i := 0; i < len(sp); i += 2 {
		if i+1 >= len(sp) {
			break
		}
		k, m := strings.TrimSpace(sp[i]), strings.TrimSpace(sp[i+1])
		k, m = strings.ReplaceAll(k, "｜", "|"), strings.ReplaceAll(m, "｜", "|")
		if (k == "*" || k == validateTag) && m != "" {
			return m, true // found
		}
	}
	return "", false // not found
}
