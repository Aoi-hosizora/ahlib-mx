package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func GetValidate() (*validator.Validate, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("gin's validator is not github.com/go-playground/validator/v10")
	}
	return v, nil
}

func GetTranslator(loc locales.Translator, translationFunc xvalidator.DefaultTranslationFunc) (ut.Translator, error) {
	v, err := GetValidate()
	if err != nil {
		return nil, err
	}
	return xvalidator.GetTranslator(v, loc, translationFunc)
}

// AddBinding adds user defined binding.
// Reference see https://godoc.org/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags.
//
// Binding tips:
//
// 1. `required` + non-pointer (Common)
// 	A uint64 `binding:"required"` // cannot be nil and 0
// 	B string `binding:"required"` // cannot be nil and ""
//
// 2. `required` + pointer (Common)
// 	A *uint64 `binding:"required"` // cannot be nil, can be 0
// 	B *string `binding:"required"` // cannot be nil, can be ""
//
// 3. `omitempty` + non-pointer (Common)
// 	A uint64 `binding:"omitempty"` // can be nil and 0
// 	B string `binding:"omitempty"` // can be nil and ""
//
// 4. `omitempty` + pointer => same as 3
// 	A *uint64 `binding:"omitempty"` // can be nil and 0
// 	B *string `binding:"omitempty"` // can be nil and ""
//
// 5. `required` + `omitempty` + non-pointer => same as 1
// 	A uint64 `binding:"required,omitempty"` // cannot be nil and 0
// 	B string `binding:"required,omitempty"` // cannot be nil and ""
//
// 6. `required` + `omitempty` + pointer => same as 2
// 	A *uint64 `binding:"required,omitempty"` // cannot be nil, can be 0
// 	B *string `binding:"required,omitempty"` // cannot be nil, can be ""
func AddBinding(tag string, fn validator.Func) error {
	v, err := GetValidate()
	if err != nil {
		return nil
	}
	return v.RegisterValidation(tag, fn)
}

// AddTranslator adds user defined validation translator to ut.Translator.
func AddTranslator(translator ut.Translator, tag, message string, override, withParam bool) error {
	v, err := GetValidate()
	if err != nil {
		return nil
	}

	transFunc := xvalidator.ValidatorTranslationFunc()
	if withParam {
		transFunc = xvalidator.ValidatorTranslationParamFunc()
	}
	return v.RegisterTranslation(tag, translator, xvalidator.ValidatorRegisterTranslationsFunc(tag, message, override), transFunc)
}

// Enable regexp binding: `regexp`.
func EnableRegexpBinding() error {
	return AddBinding("regexp", xvalidator.DefaultRegexpValidator())
}

// Enable regexp binding translator: `regexp`.
func EnableRegexpBindingTranslator(translator ut.Translator) error {
	return AddTranslator(translator, "regexp", "{0} must matches regexp /{1}/", true, true)
}

// Enable regexp binding and translator: `regexp`.
func EnableRegexpBindingWithTranslator(translator ut.Translator) error {
	err := EnableRegexpBinding()
	if err != nil {
		return err
	}
	return EnableRegexpBindingTranslator(translator)
}

// Enable rfc3339 date binding: `date`.
func EnableRFC3339DateBinding() error {
	return AddBinding("date", xvalidator.DateTimeValidator(xtime.RFC3339Date))
}

// Enable rfc3339 date translator: `date`.
func EnableRFC3339DateBindingTranslator(translator ut.Translator) error {
	return AddTranslator(translator, "date", "{0} must be an RFC3339 date", true, false)
}

// Enable rfc3339 date binding and translator: `date`.
func EnableRFC3339DateBindingWithTranslator(translator ut.Translator) error {
	err := EnableRFC3339DateBinding()
	if err != nil {
		return err
	}
	return EnableRFC3339DateBindingTranslator(translator)
}

// Enable rfc3339 regexp binding: `datetime`.
func EnableRFC3339DateTimeBinding() error {
	return AddBinding("datetime", xvalidator.DateTimeValidator(xtime.RFC3339DateTime))
}

// Enable rfc3339 regexp translator: `datetime`.
func EnableRFC3339DateTimeBindingTranslator(translator ut.Translator) error {
	return AddTranslator(translator, "datetime", "{0} must be an RFC3339 datetime", true, false)
}

// Enable rfc3339 regexp binding and translator: `datetime`.
func EnableRFC3339DateTimeBindingWithTranslator(translator ut.Translator) error {
	err := EnableRFC3339DateTimeBinding()
	if err != nil {
		return err
	}
	return EnableRFC3339DateTimeBindingTranslator(translator)
}
