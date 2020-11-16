package xvalidator

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// GetTranslator applies language default translation to ut.Translator.
func GetTranslator(validator *validator.Validate, loc locales.Translator, defaultTranslationFunc func(v *validator.Validate, trans ut.Translator) error) (ut.Translator, error) {
	uniTranslator := ut.New(loc, loc)
	translator, _ := uniTranslator.GetTranslator(loc.Locale())

	err := defaultTranslationFunc(validator, translator)
	if err != nil {
		return nil, err
	}

	return translator, nil
}

// DefaultRegisterTranslationsFunc is a default validator.RegisterTranslationsFunc for RegisterTranslation.
func DefaultRegisterTranslationsFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

// DefaultTranslationFunc is a default validator.TranslationFunc for RegisterTranslation.
func DefaultTranslationFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
}

// DefaultTranslationWithPramFunc is a default validator.TranslationFunc for RegisterTranslation.
func DefaultTranslationWithPramFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
}
