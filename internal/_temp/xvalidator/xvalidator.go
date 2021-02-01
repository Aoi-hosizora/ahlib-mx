package xvalidator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

// TODO waiting for ahlib-more v1.5.0 published

// ParamRegexpValidator represents regexp validator with param, just like `regexp: xxx`.
func ParamRegexpValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		regexpParam := fl.Param() // param
		i := fl.Field().Interface()
		text, ok := i.(string)
		if !ok {
			return false
		}
		re, err := regexp.Compile(regexpParam)
		if err != nil {
			return false
		}
		return re.MatchString(text)
	}
}

// DateTimeValidator represents datetime validator using given layout.
func DateTimeValidator(layout string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		text, ok := i.(string)
		if !ok {
			return false
		}
		_, err := time.Parse(layout, text)
		return err == nil
	}
}
