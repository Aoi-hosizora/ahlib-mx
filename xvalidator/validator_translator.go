package xvalidator

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/fr"
	"github.com/go-playground/validator/v10/translations/id"
	"github.com/go-playground/validator/v10/translations/ja"
	"github.com/go-playground/validator/v10/translations/nl"
	"github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/go-playground/validator/v10/translations/ru"
	"github.com/go-playground/validator/v10/translations/tr"
	"github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-playground/validator/v10/translations/zh_tw"
)

type DefaultTranslationFunc func(v *validator.Validate, trans ut.Translator) error

// GetTranslator applies language default translation to ut.Translator.
func GetTranslator(validator *validator.Validate, loc locales.Translator, translatorFunc DefaultTranslationFunc) (ut.Translator, error) {
	uniTranslator := ut.New(loc, loc)
	translator, _ := uniTranslator.GetTranslator(loc.Locale())

	err := translatorFunc(validator, translator)
	if err != nil {
		return nil, err
	}

	return translator, nil
}

// ValidatorRegisterTranslationsFunc is a default validator.RegisterTranslationsFunc for RegisterTranslation.
func ValidatorRegisterTranslationsFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

// ValidatorTranslationFunc is a default validator.TranslationFunc for RegisterTranslation.
func ValidatorTranslationFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
}

// ValidatorTranslationParamFunc is a default validator.TranslationFunc for RegisterTranslation.
func ValidatorTranslationParamFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
}

// From github.com/go-playground/validator/v10/translations/en.
func EnValidatorTranslation() DefaultTranslationFunc {
	return en.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/fr.
func FrValidatorTranslation() DefaultTranslationFunc {
	return fr.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/id.
func IdValidatorTranslation() DefaultTranslationFunc {
	return id.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/ja.
func JaValidatorTranslation() DefaultTranslationFunc {
	return ja.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/nl.
func NlValidatorTranslation() DefaultTranslationFunc {
	return nl.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/pt_BR.
func PtBrValidatorTranslation() DefaultTranslationFunc {
	return pt_BR.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translation/ru.
func RuValidatorTranslation() DefaultTranslationFunc {
	return ru.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/tr.
func TrValidatorTranslation() DefaultTranslationFunc {
	return tr.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/zh.
func ZhValidatorTranslation() DefaultTranslationFunc {
	return zh.RegisterDefaultTranslations
}

// From github.com/go-playground/validator/v10/translations/zh_tw.
func ZhTwValidatorTranslation() DefaultTranslationFunc {
	return zh_tw.RegisterDefaultTranslations
}
