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

// ========================
// MultiFieldsError related
// ========================

// WrappedFieldError represents a validator.FieldError wrapped with a custom message, will be used in MessagedValidator.
type WrappedFieldError struct {
	origin  validator.FieldError
	message string
}

// Origin returns the origin validator.FieldError from WrappedFieldError.
func (w *WrappedFieldError) Origin() validator.FieldError {
	return w.origin
}

// Unwrap returns the origin validator.FieldError from WrappedFieldError, and this implements the error wrapper interface.
func (w *WrappedFieldError) Unwrap() error {
	return w.origin
}

// Message returns the wrapped message from WrappedFieldError.
func (w *WrappedFieldError) Message() string {
	return w.message
}

// Error returns the formatted error message from WrappedFieldError, and this has the same format with validator.FieldError's Error().
func (w *WrappedFieldError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:%s", w.origin.Namespace(), w.message)
}

// MultiFieldsError represents the multiple fields' error in validation that contains errors in validator.FieldError or xvalidator.WrappedFieldError
// type, will be returned by MessagedValidator.ValidateStruct.
type MultiFieldsError struct {
	fields []error // validator.FieldError or xvalidator.WrappedFieldError
}

// Errors returns the fields' errors from MultiFieldsError.
func (m *MultiFieldsError) Errors() []error {
	return m.fields
}

// Error returns the formatted error message (split by "; ") from MultiFieldsError.
func (m *MultiFieldsError) Error() string {
	msgs := make([]string, 0, len(m.fields))
	for _, fe := range m.fields {
		if fe != nil {
			msgs = append(msgs, fe.Error())
		}
	}
	return strings.Join(msgs, "; ")
}

// ================
// translate & flat
// ================

// Translate translates all the field errors (include validator.FieldError and xvalidator.WrappedFieldError) using given UtTranslator to a field-message map.
// Here errors in xvalidator.WrappedFieldError type will use wrapped message directly, also note that if you set useNamespace to true, keys from returned
// map will be shown in "$struct.$field" format, otherwise in "$field" format.
//
// Example:
// 	type Struct struct {
// 		Int int    `validate:"required"`
// 		Str string `validate:"required" message:"required|str cannot be null and empty"`
// 	}
// 	val := NewMessagedValidator()
// 	// ...
// 	err := val.ValidateStruct(&Struct{}).(xvalidator.MultiFieldsError)
// 	err.Translate(trans, true)  // => {Struct.int: int is a required field, Struct.str: str cannot be null and empty}
// 	err.Translate(trans, false) // => {int:        int is a required field, str:        str cannot be null and empty}
func (m *MultiFieldsError) Translate(ut UtTranslator, useNamespace bool) map[string]string {
	if ut == nil {
		panic(panicNilUtTranslator)
	}
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(m.fields))
	for _, err := range m.fields {
		if fe, ok := err.(validator.FieldError); ok {
			result[keyFn(fe)] = fe.Translate(ut)
		} else if we, ok := err.(*WrappedFieldError); ok {
			result[keyFn(we.origin)] = we.message
		} else {
			// skip
		}
	}
	return result
}

// TranslateValidationErrors translates all validator.FieldError in validator.ValidationErrors using given UtTranslator to a field-message map. Note that
// if you set useNamespace to true, keys from returned map will be shown in "$struct.$field" format, the same as validator.ValidationErrors' Translate(),
// otherwise in "$field" format.
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

// FlatToMap flats all the field errors (include validator.FieldError and xvalidator.WrappedFieldError) to a field-message map without using UtTranslator.
// Here values from returned map come from the error message directly. Also see MultiFieldsError.Translate for more.
//
// Example:
// 	type Struct struct {
// 		Int int    `validate:"required"`
// 		Str string `validate:"required" message:"required|str cannot be null and empty"`
// 	}
// 	val := NewMessagedValidator()
// 	// ...
// 	err := validator.ValidateStruct(&Struct{}).(xvalidator.MultiFieldsError)
// 	err.FlatToMap(true)  // => {Struct.int: Field validation for 'int' failed on the 'required' tag, Struct.str: str cannot be null and empty}
// 	err.FlatToMap(false) // => {int:        Field validation for 'int' failed on the 'required' tag, str:        str cannot be null and empty}
func (m *MultiFieldsError) FlatToMap(useNamespace bool) map[string]string {
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(m.fields))
	for _, err := range m.fields {
		if fe, ok := err.(validator.FieldError); ok {
			result[keyFn(fe)] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
		} else if we, ok := err.(*WrappedFieldError); ok {
			result[keyFn(we.origin)] = we.message
		} else {
			// skip
		}
	}
	return result
}

// FlatValidationErrors flats all the validator.FieldError in validator.ValidationErrors to a field-message map without using UtTranslator. Here values
// from returned map come from the error message directly. Also see TranslateValidationErrors for more.
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
		// has the same format as fe.Error()
		result[keyFn(fe)] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
	}
	return result
}

// MapToError generates a single error from given map generated by translate or flat, to represent the translated or flatted error from
// validator.ValidationErrors or xvalidator.MultiFieldsError.
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

// =========================
// MessagedValidator related
// =========================

// MessagedValidator represents a messaged validator.Validate, which allows some fields to specify their custom error message, and you can set this to
// gin's binding.Validator as a binding.StructValidator.
//
// Struct example:
// 	type User struct {
//		Id   uint64  `json:"id"   form:"id"   validate:"required,gt=1"          validate_message:"required|id is required|gt|id must larger than one"`
//		Name string  `json:"name" form:"name" validate:"required,gt=4,lt=20"    validate_message:"*|name is invalid"`
//		Bio  *string `json:"bio"  form:"bio"  validate:"required,gte=0,lte=255" validate_message:"xxx"`
// 	}
type MessagedValidator struct {
	validate   *validator.Validate
	messageTag string

	once sync.Once
}

// NewMessagedValidator creates a new NewMessagedValidator, with `validate` validator tag name and `validate_message` message tag name.
func NewMessagedValidator() *MessagedValidator {
	m := &MessagedValidator{
		validate:   validator.New(),
		messageTag: "validate_message",
	}
	m.validate.SetTagName("validate")
	return m
}

// Engine returns the internal validator.Validate from MessagedValidator.
func (m *MessagedValidator) Engine() interface{} {
	return m.validate
}

// ValidateEngine returns the internal validator.Validate from MessagedValidator.
func (m *MessagedValidator) ValidateEngine() *validator.Validate {
	return m.validate
}

// SetValidateTagName sets validate tag name for MessagedValidator, defaults to `validate`.
func (m *MessagedValidator) SetValidateTagName(name string) {
	m.validate.SetTagName(name)
}

// SetMessageTagName sets message tag name for MessagedValidator, defaults to `validate_message`.
func (m *MessagedValidator) SetMessageTagName(name string) {
	m.messageTag = name
}

// UseTagAsFieldName sets a specific struct tag as field's alternate name, see UseTagAsFieldName for more details.
func (m *MessagedValidator) UseTagAsFieldName(name ...string) {
	UseTagAsFieldName(m.ValidateEngine(), name...)
}

// ValidateStruct validates given struct and returns the validator error, mostly in xvalidator.MultiFieldsError type.
func (m *MessagedValidator) ValidateStruct(obj interface{}) error {
	itf, ok := m.extractToStruct(obj)
	if !ok {
		return &validator.InvalidValidationError{Type: reflect.TypeOf(obj)}
	}

	err := m.validate.Struct(itf)
	if err == nil {
		return nil
	}
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return err // unreachable
	}

	typ := reflect.TypeOf(itf)
	errs := make([]error, 0, len(ve))
	for _, fe := range ve {
		if m, found := m.applyCustomMessage(typ, fe.StructField(), fe.Tag()); found {
			we := &WrappedFieldError{origin: fe, message: m}
			errs = append(errs, we) // xvalidator.WrappedFieldError
		} else {
			errs = append(errs, fe) // validator.FieldError
		}
	}
	return &MultiFieldsError{fields: errs} // MultiFieldsError
}

// extractToStruct checks and extracts given interface to struct type.
func (m *MessagedValidator) extractToStruct(obj interface{}) (interface{}, bool) {
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

// applyCustomMessage checks the struct field and wraps validator.FieldError to xvalidator.WrappedFieldError. Note that "\|" represents a single "|",
// such as "*|name \|is\| invalid" means the validation message for "*" (all tags) is "name |is| invalid".
func (m *MessagedValidator) applyCustomMessage(typ reflect.Type, fieldName, validateTag string) (string, bool) {
	sf, ok := typ.FieldByName(fieldName)
	if !ok {
		return "", false // unreachable
	}
	msg := sf.Tag.Get(m.messageTag)
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
