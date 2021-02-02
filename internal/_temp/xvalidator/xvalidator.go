package xvalidator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

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
			return false // return false
		}
		return re.MatchString(text)
	}
}

func RegexpValidator(re *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		text, ok := i.(string)
		if !ok {
			return false
		}
		return re.MatchString(text)
	}
}

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
