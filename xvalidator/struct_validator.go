package xvalidator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// ===========================
// ValidateFieldsError related
// ===========================

// WrappedValidateFieldError represents a validator.FieldError wrapped with a custom message, is used in CustomStructValidator.
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

// Unwrap returns the origin validator.FieldError from WrappedValidateFieldError, and implements the wrapper interface.
func (v *WrappedValidateFieldError) Unwrap() error {
	return v.origin
}

// Error returns the formatted error message from WrappedValidateFieldError, and has the same format with validator.FieldError.Error.
func (v *WrappedValidateFieldError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:%s", v.origin.Namespace(), v.message)
}

// ValidateFieldsError represents the validation error, which is returned by CustomStructValidator.ValidateStruct.
type ValidateFieldsError struct {
	fields []error // validator.FieldError or WrappedValidateFieldError
}

// Errors returns the fields' errors from ValidateFieldsError.
func (v *ValidateFieldsError) Errors() []error {
	return v.fields
}

// Error returns the formatted error message (split by "; ") from ValidateFieldsError.
func (v *ValidateFieldsError) Error() string {
	msgs := make([]string, 0, len(v.fields))
	for _, fe := range v.fields {
		if fe != nil {
			msgs = append(msgs, fe.Error())
		}
	}
	return strings.Join(msgs, "; ")
}

// ================
// translate & flat
// ================

// Translate translates all field errors (include validator.FieldError and WrappedValidateFieldError) using given UtTranslator to a field-message map. 
// Here errors in WrappedValidateFieldError type will use wrapped message directly, also note that if you set useNamespace to true, keys from returned 
// map will be shown in "$struct.$field" format, otherwise in "$field" format.
//
// Example:
// 	type Struct struct {
// 		Int int    `validate:"required"`
// 		Str string `validate:"required" message:"required|str cannot be null and empty"`
// 	}
// 	val := NewCustomStructValidator()
// 	// ...
// 	err := val.ValidateStruct(&Struct{}).(xvalidator.ValidateFieldsError)
// 	err.Translate(trans, true)  // => {Struct.int: int is a required field, Struct.str: str cannot be null and empty}
// 	err.Translate(trans, false) // => {int:        int is a required field, str:        str cannot be null and empty}
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

// TranslateValidationErrors translates all validator.FieldError in validator.ValidationErrors using given UtTranslator to a field-message map. Note that 
// if you set useNamespace to true, keys from returned map will be shown in "$struct.$field" format, that is the same with validator.ValidationErrors' 
// Translate(), otherwise in "$field" format.
//
// Example:
// 	type Struct struct {
// 		Int int    `validate:"required"`
// 		Str string `validate:"required"`
// 	}
// 	val := validator.New()
// 	// ...
// 	err := val.Struct(&Struct{}).(validator.ValidationErrors)
// 	TranslateValidationErrors(err, trans, true)  // => {Struct.int: int is a required field, Struct.str: str is a required field}
// 	TranslateValidationErrors(err, trans, false) // => {int:        int is a required field, str:        str is a required field}
func TranslateValidationErrors(err validator.ValidationErrors, ut UtTranslator, useNamespace bool) map[string]string {
	if ut == nil {
		panic(panicNilUtTranslator)
	}
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(err))
	for _, fe := range err {
		result[keyFn(fe)] = fe.Translate(ut)
	}
	return result
}

// FlatToMap flats all field errors (include validator.FieldError and WrappedValidateFieldError) to a field-message map without using UtTranslator. 
// Here values from returned map will be the error message directly. Also see ValidateFieldsError.Translate for more.
//
// Example:
// 	type Struct struct {
// 		Int int    `validate:"required"`
// 		Str string `validate:"required" message:"required|str cannot be null and empty"`
// 	}
// 	val := NewCustomStructValidator()
// 	// ...
// 	err := validator.ValidateStruct(&Struct{}).(xvalidator.ValidateFieldsError)
// 	err.FlatToMap(true)  // => {Struct.int: Field validation for 'int' failed on the 'required' tag, Struct.str: str cannot be null and empty}
// 	err.FlatToMap(false) // => {int:        Field validation for 'int' failed on the 'required' tag, str:        str cannot be null and empty}
func (v *ValidateFieldsError) FlatToMap(useNamespace bool) map[string]string {
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(v.fields))
	for _, err := range v.fields {
		if fe, ok := err.(validator.FieldError); ok {
			result[keyFn(fe)] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
		} else if we, ok := err.(*WrappedValidateFieldError); ok {
			result[keyFn(we.origin)] = we.message
		} else {
			// skip
		}
	}
	return result
}

// FlatValidationErrors flats all all validator.FieldError in validator.ValidationErrors to a field-message map without using UtTranslator. Here values
// from returned map will be the error message directly. Also see TranslateValidationErrors for more.
//
// Example:
// 	type Struct struct {
// 		Int int    `validate:"required"`
// 		Str string `validate:"required"`
// 	}
// 	val := validator.New()
// 	// ...
// 	err := val.Struct(&Struct{}).(validator.ValidationErrors)
// 	FlatValidationErrors(err, true)  // => {Struct.int: Field validation for 'int' failed on the 'required' tag, Struct.str: Field validation for 'str' failed on the 'required' tag}
// 	FlatValidationErrors(err, false) // => {int:        Field validation for 'int' failed on the 'required' tag, str:        Field validation for 'str' failed on the 'required' tag}
func FlatValidationErrors(err validator.ValidationErrors, useNamespace bool) map[string]string {
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(err))
	for _, fe := range err {
		result[keyFn(fe)] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
	}
	return result
}

// MapToError generates an error from given translated or flatted map to represent the translated or flatted error from validator.ValidationErrors or
// xvalidator.ValidateFieldsError.
func MapToError(result map[string]string) error {
	if len(result) == 0 {
		return nil
	}
	tuples := make([][2]string, 0, len(result)) // map -> [][2]string
	for k, v := range result {
		tuples = append(tuples, [2]string{k, v})
	}
	sort.Slice(tuples, func(i, j int) bool { // [][2]string -> sorted [][2]string
		return tuples[i][0] < tuples[j][0] || (tuples[i][0] == tuples[j][0] && tuples[i][1] < tuples[j][1])
	})
	msgs := make([]string, 0, len(result)) // sorted [][2]string -> []string
	for _, t := range tuples {
		msgs = append(msgs, t[1])
	}
	return errors.New(strings.Join(msgs, "; "))
}

// =============================
// CustomStructValidator related
// =============================

// CustomStructValidator represents a custom validator.Validate, which allows some fields to specify their custom error message, and you can set this to
// gin's binding.Validator as a binding.StructValidator.
//
// Struct example:
// 	type User struct {
//		Id   uint64  `json:"id"   form:"id"   validate:"required,gt=1"          validate_message:"required|id is required|gt|id must larger than one"`
//		Name string  `json:"name" form:"name" validate:"required,gt=4,lt=20"    validate_message:"*|name is invalid"`
//		Bio  *string `json:"bio"  form:"bio"  validate:"required,gte=0,lte=255" validate_message:"xxx"`
// 	}
type CustomStructValidator struct {
	validate   *validator.Validate
	messageTag string

	once sync.Once
}

// NewCustomStructValidator creates a new NewCustomStructValidator, with `validate` validator tag name and `validate_message` message tag name.
func NewCustomStructValidator() *CustomStructValidator {
	v := &CustomStructValidator{
		validate:   validator.New(),
		messageTag: "validate_message",
	}
	v.validate.SetTagName("validate")
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

// SetValidatorTagName sets validator tag name for CustomStructValidator, defaults to `validate`.
func (v *CustomStructValidator) SetValidatorTagName(name string) {
	v.validate.SetTagName(name)
}

// SetMessageTagName sets message tag name for CustomStructValidator, defaults to `validate_message`.
func (v *CustomStructValidator) SetMessageTagName(name string) {
	v.messageTag = name
}

// SetFieldNameTag sets a specific struct tag as field's alternate name, see UseTagAsFieldName for more details.
func (v *CustomStructValidator) SetFieldNameTag(name string) {
	UseTagAsFieldName(v.ValidateEngine(), name)
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

// applyCustomMessage checks the struct field and wraps validator.FieldError to WrappedValidateFieldError. Note that "\|" is used to represent
// a single "|", such as "*|name \|is\| invalid" means the message for "*" is set to "name |is| invalid".
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
