package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

// Setup userDefined bind.
//
// Binding tips:
//
// 1. If you want to allow empty value but not allow nil, then use `required` with pointer type of field. (most common)
//
// Example:
// 	Gender   *uint8  `binding:"required,o_gender"`  // can be 0, but not nil
// 	Profile  *string `binding:"required,l_profile"` // can be "", but not nul
//
// 2. If you don't require the field to appear, but allow empty, then use `omitempty` without pointer type of field. (most common)
//
// Example:
// 	Email    string `binding:"omitempty,l_email,email"` // can be nil, but also be ""
// 	Username string `binding:"omitempty,l_name,r_name"` // can be nil, but also be ""
//
// 3. If you have a `required` (ignore `omitempty`) without pointer, the when you used a zero value, it will throw error.
//
// Example:
// 	Gender   uint8  `binding:"required,o_gender"`            // could not be 0, and nil
// 	Profile  string `binding:"required,omitempty,l_profile"` // could not be "", and nil
//
// 4. Just no `required` and `omitempty`, the field will depend on other remain tag.
// Example:
// 	Email    string `binding:"l_email,email"` // can not be nil, and ""
// 	Username string `binding:""`              // can be nil, but also be ""
//
func SetupBinding(tag string, fn func(fl validator.FieldLevel) bool) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, fn)
	}
}

// Setup userDefined bind (using string).
func SetupBindingWithString(tag string, fn func(str string) bool) {
	SetupBinding(tag, func(fl validator.FieldLevel) bool {
		f := fmt.Sprintf("%v", fl.Field().Interface())
		return fn(f)
	})
}

// Setup binding tag: `regexp=xxx`.
func EnableRegexpBinding() {
	SetupBinding("regexp", func(fl validator.FieldLevel) bool {
		f := fmt.Sprintf("%v", fl.Field().Interface())
		re := regexp.MustCompile(fl.Param())
		return re.MatchString(f)
	})
}

// Setup binding tag: `$tag·.
func SetupRegexpBinding(tag string, re *regexp.Regexp) {
	SetupBindingWithString(tag, func(f string) bool {
		return re.MatchString(f)
	})
}

// Setup binging tag for datetime.
func SetupDateTimeBinding(tag string, layout string) {
	SetupBindingWithString(tag, func(f string) bool {
		_, err := time.Parse(layout, f)
		if err != nil {
			return false
		}
		return true
	})
}

// Setup length binding tag: $tag.
func SetupLengthBinding(tag string, min, max int) {
	SetupBindingWithString(tag, func(f string) bool {
		return len(f) >= min && len(f) <= max
	})
}

// Setup oneof binding tag: $tag.
func SetupOneofBinding(tag string, fields ...interface{}) {
	SetupBindingWithString(tag, func(f string) bool {
		for _, field := range fields {
			if f == fmt.Sprintf("%v", field) {
				return true
			}
		}
		return false
	})
}

// Check is err is validator.ValidationErrors and is required error.
func IsValidationFormatError(err error) bool {
	// validator.ValidationErrors
	// *errors.errorString
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}

	for _, field := range errs {
		if field.Tag() == "required" {
			return false
		}
	}
	return true
}
