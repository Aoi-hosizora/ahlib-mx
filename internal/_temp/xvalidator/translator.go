package xvalidator

import (
	"github.com/go-playground/locales"
	loc_en "github.com/go-playground/locales/en"
	loc_fr "github.com/go-playground/locales/fr"
	loc_ja "github.com/go-playground/locales/ja"
	loc_ru "github.com/go-playground/locales/ru"
	loc_zh "github.com/go-playground/locales/zh"
	loc_zh_Hant "github.com/go-playground/locales/zh_Hant"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	trans_en "github.com/go-playground/validator/v10/translations/en"
	trans_fr "github.com/go-playground/validator/v10/translations/fr"
	trans_ja "github.com/go-playground/validator/v10/translations/ja"
	trans_ru "github.com/go-playground/validator/v10/translations/ru"
	trans_zh "github.com/go-playground/validator/v10/translations/zh"
	trans_zh_tw "github.com/go-playground/validator/v10/translations/zh_tw"
)

// TODO waiting for ahlib-more v1.5.0 published

// TranslationRegisterHandler represents a translation register function, which is the type of en.RegisterDefaultTranslations,
// zh.RegisterDefaultTranslations and so on, here `en` and `zh` packages are from github.com/go-playground/validator/v10/translations.
//
// These kind of values can be get from DefaultTranslateFunc, EnTranslationFunc, ZhTranslationFunc and so on.
type TranslationRegisterHandler func(v *validator.Validate, trans ut.Translator) error

const (
	panicNilValidator             = "xvalidator: nil validator"
	panicNilLocaleTranslator      = "xvalidator: nil locale translator"
	panicNilTranslationRegisterFn = "xvalidator: nil translation register function"
)

// ApplyValidatorTranslator applies translation to validator.Validate using given locales.Translator and TranslationRegisterHandler, and
// returns a ut.Translator (universal translator).
//
// Note:
//
// 1. locales.Translator can be get from en.New(), zh.New() and so on, here `en` and `zh` packages are from github.com/go-playground/locales.
// These kind of values can also be get from xvalidator.EnLocaleTranslator, xvalidator.ZhLocaleTranslator and so on;
//
// 2. xvalidator.TranslationRegisterHandler can be get from en.RegisterDefaultTranslations, zh.RegisterDefaultTranslations and so on, here
// `en` and `zh` packages are from github.com/go-playground/validator/v10/translations. These kind of values can also be get from
// xvalidator.DefaultTranslationRegisterFunc, xvalidator.EnTranslationRegisterFunc, xvalidator.ZhTranslationRegisterFunc and so on.
//
// Also see xvalidator.AddToTranslatorFunc and xvalidator.DefaultTranslateFunc.
// Example:
// 	validator := validator.New()
// 	translator := xvalidator.ApplyValidatorTranslator(validator, en.New(), xvalidator.EnTranslationRegisterFunc()) // ut.Translator
// 	translatorFunc := xvalidator.AddToTranslatorFunc("tag", "{0} has {1}", false) // validator.RegisterTranslationsFunc
// 	_ = validator.RegisterTranslation("tag", translator, translatorFunc, xvalidator.DefaultTranslateFunc())
func ApplyValidatorTranslator(validator *validator.Validate, locTrans locales.Translator, registerFn TranslationRegisterHandler) (ut.Translator, error) {
	if validator == nil {
		panic(panicNilValidator)
	}
	if locTrans == nil {
		panic(panicNilLocaleTranslator)
	}
	if registerFn == nil {
		panic(panicNilTranslationRegisterFn)
	}

	uniTranslator := ut.New(locTrans, locTrans)
	translator, _ := uniTranslator.GetTranslator(locTrans.Locale())
	err := registerFn(validator, translator) // register translation to validator
	if err != nil {
		return nil, err
	}

	return translator, nil
}

// ====================
// register & translate
// ====================

// AddToTranslatorFunc returns a validator.RegisterTranslationsFunc function, it uses the given tag, translation and override parameters
// to **register translation information** for a ut.Translator and will be used within the validator.TranslationFunc.
//
// This function can be used for validator.Validate's RegisterTranslation() method's translationFn parameter.
func AddToTranslatorFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

// DefaultTranslateFunc returns a validator.TranslationFunc function, it uses the field name (validator.FieldError's field) as the first
// parameter ({0}), and the field param (validator.FieldError's param) as the second parameter ({1}) to **create translation for the given tag**.
//
// This function can be used for validator.Validate's RegisterTranslation() method's registerFn parameter.
func DefaultTranslateFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param()) // field to {0} & param to {1}
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
}

// =================
// locale translator
// =================

// EnLocaleTranslator is a locales.Translator for en.New() from github.com/go-playground/locales/en.
func EnLocaleTranslator() locales.Translator {
	return loc_en.New()
}

// FrLocaleTranslator is a locales.Translator for fr.New() from github.com/go-playground/locales/fr.
func FrLocaleTranslator() locales.Translator {
	return loc_fr.New()
}

// JaLocaleTranslator is a locales.Translator for ja.New() from github.com/go-playground/locales/ja.
func JaLocaleTranslator() locales.Translator {
	return loc_ja.New()
}

// RuLocaleTranslator is a locales.Translator for ru.New() from github.com/go-playground/locales/ru.
func RuLocaleTranslator() locales.Translator {
	return loc_ru.New()
}

// ZhLocaleTranslator is a locales.Translator for zh.New() from github.com/go-playground/locales/zh.
func ZhLocaleTranslator() locales.Translator {
	return loc_zh.New()
}

// ZhHantLocaleTranslator is a locales.Translator for zh_Hant.New() from github.com/go-playground/locales/zh_Hant.
func ZhHantLocaleTranslator() locales.Translator {
	return loc_zh_Hant.New()
}

// ============================
// default translation register
// ============================

// EnTranslationRegisterFunc is a TranslationRegisterHandler for en.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/en.
func EnTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_en.RegisterDefaultTranslations
}

// FrTranslationRegisterFunc is a TranslationRegisterHandler for fr.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/fr.
func FrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_fr.RegisterDefaultTranslations
}

// JaTranslationRegisterFunc is a TranslationRegisterHandler for ja.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/ja.
func JaTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ja.RegisterDefaultTranslations
}

// RuTranslationRegisterFunc is a TranslationRegisterHandler for ru.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translation/ru.
func RuTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ru.RegisterDefaultTranslations
}

// ZhTranslationRegisterFunc is a TranslationRegisterHandler for zh.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/zh.
func ZhTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_zh.RegisterDefaultTranslations
}

// ZhTwTranslationRegisterFunc is a TranslationRegisterHandler for zh_tw.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/zh_tw.
func ZhTwTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_zh_tw.RegisterDefaultTranslations
}
