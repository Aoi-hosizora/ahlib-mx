package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Add user defined binding.
// Reference see https://godoc.org/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags.
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
func AddBinding(tag string, fn validator.Func) error {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("gin's validator is not validator.Validate")
	}
	return val.RegisterValidation(tag, fn)
}

// Enable regexp binding: `regexp`.
func EnableRegexpBinding() error {
	return AddBinding("regexp", xvalidator.DefaultRegexpValidator())
}

// Enable rfc3339 date binding: `date`.
func EnableRFC3339DateBinding() error {
	return AddBinding("date", xvalidator.DateTimeValidator(xtime.RFC3339Date))
}

// Enable rfc3339 regexp binding: `datetime`.
func EnableRFC3339DateTimeBinding() error {
	return AddBinding("datetime", xvalidator.DateTimeValidator(xtime.RFC3339DateTime))
}
