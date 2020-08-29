package xfiber

import (
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func AddBinding(tag string, fn validator.Func) error {
	return Validator.RegisterValidation(tag, fn)
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

func Struct(s interface{}) error {
	return Validator.Struct(s)
}

func Var(field interface{}, tag string) error {
	return Validator.Var(field, tag)
}
