package xvalidator

import (
	"github.com/go-playground/locales"
	loc_en "github.com/go-playground/locales/en"
	loc_fr "github.com/go-playground/locales/fr"
	loc_id "github.com/go-playground/locales/id"
	loc_ja "github.com/go-playground/locales/ja"
	loc_nl "github.com/go-playground/locales/nl"
	loc_pt_BR "github.com/go-playground/locales/pt_BR"
	loc_ru "github.com/go-playground/locales/ru"
	loc_tr "github.com/go-playground/locales/tr"
	loc_zh "github.com/go-playground/locales/zh"
	loc_zh_Hant "github.com/go-playground/locales/zh_Hant"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	trans_en "github.com/go-playground/validator/v10/translations/en"
	trans_fr "github.com/go-playground/validator/v10/translations/fr"
	trans_id "github.com/go-playground/validator/v10/translations/id"
	trans_ja "github.com/go-playground/validator/v10/translations/ja"
	trans_nl "github.com/go-playground/validator/v10/translations/nl"
	trans_pt_BR "github.com/go-playground/validator/v10/translations/pt_BR"
	trans_ru "github.com/go-playground/validator/v10/translations/ru"
	trans_tr "github.com/go-playground/validator/v10/translations/tr"
	trans_zh "github.com/go-playground/validator/v10/translations/zh"
	trans_zh_tw "github.com/go-playground/validator/v10/translations/zh_tw"
	"reflect"
	"strings"
)

// UtTranslator represents an alias type of ut.Translator, and this represent the translator of validator.Validate.
type UtTranslator = ut.Translator

// LocaleTranslator represents an alias type of locales.Translator, which will be used in ApplyTranslator. These kind of values can be get
// from xvalidator.EnLocaleTranslator, xvalidator.ZhLocaleTranslator and so on.
//
// Note: Without using this package, locales.Translator can be get from en.New(), zh.New() and so on, here `en` and `zh` packages are
// from github.com/go-playground/locales.
type LocaleTranslator = locales.Translator

// TranslationRegisterHandler represents a translation register function, which will be used in ApplyTranslator. These kind of values can be get
// from xvalidator.EnTranslationRegisterFunc, xvalidator.ZhTranslationRegisterFunc and so on.
//
// Note: Without using this package, xvalidator.TranslationRegisterHandler can be get from en.RegisterDefaultTranslations, zh.RegisterDefaultTranslations
// and so on, here `en` and `zh` packages are from github.com/go-playground/validator/v10/translations.
type TranslationRegisterHandler func(v *validator.Validate, trans ut.Translator) error

const (
	panicNilValidator             = "xvalidator: nil validator"
	panicNilLocaleTranslator      = "xvalidator: nil locale translator"
	panicNilTranslationRegisterFn = "xvalidator: nil translation register function"
)

// ApplyTranslator applies translator to validator.Validate using given locales.Translator and TranslationRegisterHandler, this will return
// a ut.Translator (universal translator). Also see xvalidator.AddToTranslatorFunc and xvalidator.DefaultTranslateFunc.
//
// Example:
// 	// apply default translation to validator
// 	translator := xvalidator.ApplyTranslator(validator, xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc()) // ut.Translator
//
// 	// register custom translation to validator
// 	fn := xvalidator.AddToTranslatorFunc("tag", "{0} has {1}", false) // validator.RegisterTranslationsFunc
// 	validator.RegisterTranslation("tag", translator, fn, xvalidator.DefaultTranslateFunc())
func ApplyTranslator(validator *validator.Validate, locTranslator locales.Translator, registerFn TranslationRegisterHandler) (ut.Translator, error) {
	if validator == nil {
		panic(panicNilValidator)
	}
	if locTranslator == nil {
		panic(panicNilLocaleTranslator)
	}
	if registerFn == nil {
		panic(panicNilTranslationRegisterFn)
	}

	translator, _ := ut.New(locTranslator, locTranslator).GetTranslator(locTranslator.Locale())
	err := registerFn(validator, translator) // register translator to validator
	if err != nil {
		return nil, err
	}

	return translator, nil
}

// UseJsonTagAsFieldName sets json tag as struct field's alternate name using validator.Validate's RegisterTagNameFunc method. This alternate field name will
// be used in validator.FieldError's Namespace() and Field() methods, and will change validator.FieldError's error message and the translator result.
func UseJsonTagAsFieldName(validator *validator.Validate) {
	if validator == nil {
		panic(panicNilValidator)
	}
	validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonTag, ok := field.Tag.Lookup("json")
		if !ok {
			return ""
		}
		name := strings.SplitN(jsonTag, ",", 2)[0]
		if name == "" || name == "-" {
			return ""
		}
		return name // fe.Field() -> json tag
	})
}

// =========================
// register & translate func
// =========================

// AddToTranslatorFunc returns a validator.RegisterTranslationsFunc function, it uses the given tag, translation and override parameters
// to **register translation information** for a ut.Translator and will be used within the validator.TranslationFunc.
//
// This function can be used for validator.Validate RegisterTranslation() method's translationFn (the second) parameter, also see ApplyTranslator.
func AddToTranslatorFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

// DefaultTranslateFunc returns a validator.TranslationFunc function, it uses the field name (validator.FieldError's field) as the first
// parameter ({0}), and the field param (validator.FieldError's param) as the second parameter ({1}) to **create translation for the given tag**.
//
// This function can be used for validator.Validate RegisterTranslation() method's registerFn (the third) parameter, also see ApplyTranslator.
func DefaultTranslateFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param()) // field to {0}, param to {1}
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

// IdLocaleTranslator is a locales.Translator for id.New() from github.com/go-playground/locales/id.
func IdLocaleTranslator() locales.Translator {
	return loc_id.New()
}

// JaLocaleTranslator is a locales.Translator for ja.New() from github.com/go-playground/locales/ja.
func JaLocaleTranslator() locales.Translator {
	return loc_ja.New()
}

// NlLocaleTranslator is a locales.Translator for nl.New() from github.com/go-playground/locales/nl.
func NlLocaleTranslator() locales.Translator {
	return loc_nl.New()
}

// PtBrLocaleTranslator is a locales.Translator for pt_BR.New() from github.com/go-playground/locales/pt_BR.
func PtBrLocaleTranslator() locales.Translator {
	return loc_pt_BR.New()
}

// RuLocaleTranslator is a locales.Translator for ru.New() from github.com/go-playground/locales/ru.
func RuLocaleTranslator() locales.Translator {
	return loc_ru.New()
}

// TrLocaleTranslator is a locales.Translator for tr.New() from github.com/go-playground/locales/tr.
func TrLocaleTranslator() locales.Translator {
	return loc_tr.New()
}

// ZhLocaleTranslator is a locales.Translator for zh.New() from github.com/go-playground/locales/zh.
func ZhLocaleTranslator() locales.Translator {
	return loc_zh.New()
}

// ZhHantLocaleTranslator is a locales.Translator for zh_Hant.New() from github.com/go-playground/locales/zh_Hant.
func ZhHantLocaleTranslator() locales.Translator {
	return loc_zh_Hant.New()
}

// =========================
// translation register func
// =========================

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

// IdTranslationRegisterFunc is a TranslationRegisterHandler for id.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/id.
func IdTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_id.RegisterDefaultTranslations
}

// JaTranslationRegisterFunc is a TranslationRegisterHandler for ja.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/ja.
func JaTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ja.RegisterDefaultTranslations
}

// NlTranslationRegisterFunc is a TranslationRegisterHandler for nl.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/nl.
func NlTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_nl.RegisterDefaultTranslations
}

// PtBrTranslationRegisterFunc is a TranslationRegisterHandler for pt_BR.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/pt_BR.
func PtBrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_pt_BR.RegisterDefaultTranslations
}

// RuTranslationRegisterFunc is a TranslationRegisterHandler for ru.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translation/ru.
func RuTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ru.RegisterDefaultTranslations
}

// TrTranslationRegisterFunc is a TranslationRegisterHandler for tr.RegisterDefaultTranslations.
// From github.com/go-playground/validator/v10/translations/tr.
func TrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_tr.RegisterDefaultTranslations
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
