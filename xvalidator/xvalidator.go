package xvalidator

import (
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/go-playground/locales"
	loc_en "github.com/go-playground/locales/en"
	loc_es "github.com/go-playground/locales/es"
	loc_fr "github.com/go-playground/locales/fr"
	loc_id "github.com/go-playground/locales/id"
	loc_ja "github.com/go-playground/locales/ja"
	loc_nl "github.com/go-playground/locales/nl"
	loc_pt "github.com/go-playground/locales/pt"
	loc_pt_BR "github.com/go-playground/locales/pt_BR"
	loc_ru "github.com/go-playground/locales/ru"
	loc_tr "github.com/go-playground/locales/tr"
	loc_zh "github.com/go-playground/locales/zh"
	loc_zh_Hant "github.com/go-playground/locales/zh_Hant"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	trans_en "github.com/go-playground/validator/v10/translations/en"
	trans_es "github.com/go-playground/validator/v10/translations/es"
	trans_fr "github.com/go-playground/validator/v10/translations/fr"
	trans_id "github.com/go-playground/validator/v10/translations/id"
	trans_ja "github.com/go-playground/validator/v10/translations/ja"
	trans_nl "github.com/go-playground/validator/v10/translations/nl"
	trans_pt "github.com/go-playground/validator/v10/translations/pt"
	trans_pt_BR "github.com/go-playground/validator/v10/translations/pt_BR"
	trans_ru "github.com/go-playground/validator/v10/translations/ru"
	trans_tr "github.com/go-playground/validator/v10/translations/tr"
	trans_zh "github.com/go-playground/validator/v10/translations/zh"
	trans_zh_tw "github.com/go-playground/validator/v10/translations/zh_tw"
	"log"
	"reflect"
	"strings"
)

// =================
// validator related
// =================

// IsValidationError returns true if the error is validator.ValidationErrors.
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

// IsRequiredError returns true if the error is validator.ValidationErrors which contains "required" tag.
func IsRequiredError(err error) bool {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}

	for _, fe := range ve {
		if fe.Tag() == "required" || fe.ActualTag() == "required" {
			return true
		}
	}
	return false
}

const (
	panicNilValidator = "xvalidator: nil validator"
	panicEmptyTagName = "xvalidator: empty tag name"
)

// UseTagAsFieldName sets a specific struct tag as field's alternate name, this name will be used and returned by validator.FieldError's Namespace() and Field()
// methods, and the origin name will be returned by StructField() and StructNamespace() methods. Note that this will change validator.FieldError's the return
// value of Error() and the translated result.
//
// Example:
// 	v := validator.New()
// 	xvalidator.UseTagAsFieldName(v, "json")
// 	type s struct {
// 		Str string `validate:"required,gt=2,lt=10" json:"sss"`
// 	}
// 	errs := v.Struct(&s{"1234567890"}).(validator.ValidationErrors)
// 	// errs[0].Field()           => sss
// 	// errs[0].Namespace()       => s.sss
// 	// errs[0].StructField()     => Str
// 	// errs[0].StructNamespace() => s.Str
// 	// errs[0].Error()           => Key: 's.sss' Error:Field validation for 'sss' failed on the 'lt' tag
func UseTagAsFieldName(v *validator.Validate, tagName string) {
	if v == nil {
		panic(panicNilValidator)
	}
	if tagName = strings.TrimSpace(tagName); tagName == "" {
		panic(panicEmptyTagName)
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tagMsg, ok := field.Tag.Lookup(tagName)
		if !ok {
			return "" // use origin name
		}
		name := strings.SplitN(tagMsg, ",", 2)[0]
		if name == "" || name == "-" {
			return ""
		}
		return name // fe.Field() -> tag msg
	})
}

// UseDefaultFieldName unregisters tagNameFunc for validator.Validate and remove field's alternate name, also see UseTagAsFieldName.
func UseDefaultFieldName(v *validator.Validate) {
	if v == nil {
		panic(panicNilValidator)
	}
	var nilFunc validator.TagNameFunc = nil
	xreflect.SetUnexportedField(reflect.ValueOf(v).Elem().FieldByName("tagNameFunc"), reflect.ValueOf(nilFunc))
	xreflect.SetUnexportedField(reflect.ValueOf(v).Elem().FieldByName("hasTagNameFunc"), reflect.ValueOf(false))
}

// ==================
// translator related
// ==================

// UtTranslator represents an alias type of UtTranslator, and this is the translator of validator.Validate.
type UtTranslator = ut.Translator

// LocaleTranslator represents an alias type of LocaleTranslator, which will be used in ApplyTranslator. These kinds of values can be got
// from xvalidator.EnLocaleTranslator, xvalidator.ZhLocaleTranslator and so on.
type LocaleTranslator = locales.Translator

// TranslationRegisterHandler represents a translation register function, which will be used in ApplyTranslator. These kinds of values can be got
// from xvalidator.EnTranslationRegisterFunc, xvalidator.ZhTranslationRegisterFunc and so on.
type TranslationRegisterHandler func(v *validator.Validate, translator UtTranslator) error

const (
	panicNilLocaleTranslator      = "xvalidator: nil locale translator"
	panicNilTranslationRegisterFn = "xvalidator: nil translation register function"
)

// ApplyTranslator applies translator to validator.Validate using given LocaleTranslator (locales.Translator) and TranslationRegisterHandler, this function
// will return a UtTranslator (ut.Translator, universal translator). Also see xvalidator.DefaultRegistrationFunc and xvalidator.DefaultTranslateFunc.
//
// Example:
// 	// apply default translation to validator
// 	translator := xvalidator.ApplyTranslator(validator, xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc()) // UtTranslator
//
// 	// register custom translation to validator
// 	regisFn := xvalidator.DefaultRegistrationFunc("tag", "{0} has {1}", false) // validator.RegisterTranslationsFunc
// 	transFn := xvalidator.DefaultTranslateFunc() // validator.TranslationFunc
// 	validator.RegisterTranslation("tag", translator, regisFn, transFn)
func ApplyTranslator(validator *validator.Validate, locale LocaleTranslator, registerFn TranslationRegisterHandler) (UtTranslator, error) {
	if validator == nil {
		panic(panicNilValidator)
	}
	if locale == nil {
		panic(panicNilLocaleTranslator)
	}
	if registerFn == nil {
		panic(panicNilTranslationRegisterFn)
	}

	translator, _ := ut.New(locale, locale).GetTranslator(locale.Locale())
	err := registerFn(validator, translator) // register translator to validator (by validator.RegisterTranslation)
	if err != nil {
		return nil, err
	}

	return translator, nil
}

// DefaultRegistrationFunc returns a validator.RegisterTranslationsFunc function, it uses the given tag, translation and override flag to register
// normal translation information for a UtTranslator, {#} is the only replacement type accepted and will be set by validator.TranslationFunc.
//
// This function can be used for validator.Validate RegisterTranslation() method's second parameter translationFn, also see ApplyTranslator.
func DefaultRegistrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut UtTranslator) error {
		return ut.Add(tag, translation, override)
		// ignore ut.AddCardinal (ut.C), ut.AddOrdinal (ut.O), ut.AddRange (ut.R)
	}
}

// DefaultTranslateFunc returns a validator.TranslationFunc function, it uses the struct field name as {0} and the validator tag param as {1} from
// validator.FieldError's methods to create the translation for the given tag. Note that if the tag is not found, it will log a warning message.
//
// This function can be used for validator.Validate RegisterTranslation() method's third parameter registerFn, also see ApplyTranslator.
func DefaultTranslateFunc() validator.TranslationFunc {
	return func(ut UtTranslator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param()) // {0} => fe.Field(), {1} => fe.Param()
		if err != nil {
			// ut.ErrUnknownTranslation
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}
}

// TranslateValidationErrors translates all the field errors using given UtTranslator to a field-message map. Note that if you set useNamespace
// flag to true, the keys from returned map will be shown in "$struct.$field" format, that is the same with validator.ValidationErrors' Translate(),
// otherwise those will be shown in "$field" format.
//
// Example:
// 	err := validator.Struct(&Struct{}).(validator.ValidationErrors)
// 	TranslateValidationErrors(err, trans, true)  // => map[Struct.int:int is a required field, Struct.str:str is a required field]
// 	TranslateValidationErrors(err, trans, false) // => map[int:int is a required field, str:str is a required field]
func TranslateValidationErrors(err validator.ValidationErrors, ut UtTranslator, useNamespace bool) map[string]string {
	keyFn := func(e validator.FieldError) string {
		if useNamespace {
			return e.Namespace()
		}
		return e.Field()
	}

	result := make(map[string]string, len(err))
	for _, fe := range err {
		result[keyFn(fe)] = fe.Translate(ut)
	}
	return result
}

// =================
// locale translator
// =================

// EnLocaleTranslator is a LocaleTranslator generated by en.New() from github.com/go-playground/locales/en.
func EnLocaleTranslator() LocaleTranslator {
	return loc_en.New()
}

// EsLocaleTranslator is a LocaleTranslator generated by es.New() from github.com/go-playground/locales/es.
func EsLocaleTranslator() LocaleTranslator {
	return loc_es.New()
}

// FrLocaleTranslator is a LocaleTranslator generated by fr.New() from github.com/go-playground/locales/fr.
func FrLocaleTranslator() LocaleTranslator {
	return loc_fr.New()
}

// IdLocaleTranslator is a LocaleTranslator generated by id.New() from github.com/go-playground/locales/id.
func IdLocaleTranslator() LocaleTranslator {
	return loc_id.New()
}

// JaLocaleTranslator is a LocaleTranslator generated by ja.New() from github.com/go-playground/locales/ja.
func JaLocaleTranslator() LocaleTranslator {
	return loc_ja.New()
}

// NlLocaleTranslator is a LocaleTranslator generated by nl.New() from github.com/go-playground/locales/nl.
func NlLocaleTranslator() LocaleTranslator {
	return loc_nl.New()
}

// PtLocaleTranslator is a LocaleTranslator generated by pt.New() from github.com/go-playground/locales/pt.
func PtLocaleTranslator() LocaleTranslator {
	return loc_pt.New()
}

// PtBrLocaleTranslator is a LocaleTranslator generated by pt_BR.New() from github.com/go-playground/locales/pt_BR.
func PtBrLocaleTranslator() LocaleTranslator {
	return loc_pt_BR.New()
}

// RuLocaleTranslator is a LocaleTranslator generated by ru.New() from github.com/go-playground/locales/ru.
func RuLocaleTranslator() LocaleTranslator {
	return loc_ru.New()
}

// TrLocaleTranslator is a LocaleTranslator generated by tr.New() from github.com/go-playground/locales/tr.
func TrLocaleTranslator() LocaleTranslator {
	return loc_tr.New()
}

// ZhLocaleTranslator is a LocaleTranslator generated by zh.New() from github.com/go-playground/locales/zh.
func ZhLocaleTranslator() LocaleTranslator {
	return loc_zh.New()
}

// ZhHantLocaleTranslator is a LocaleTranslator generated by zh_Hant.New() from github.com/go-playground/locales/zh_Hant.
func ZhHantLocaleTranslator() LocaleTranslator {
	return loc_zh_Hant.New()
}

// ============================
// translation register handler
// ============================

// EnTranslationRegisterFunc is a TranslationRegisterHandler generated by en.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/en.
func EnTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_en.RegisterDefaultTranslations
}

// EsTranslationRegisterFunc is a TranslationRegisterHandler generated by es.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/es.
func EsTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_es.RegisterDefaultTranslations
}

// FrTranslationRegisterFunc is a TranslationRegisterHandler generated by fr.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/fr.
func FrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_fr.RegisterDefaultTranslations
}

// IdTranslationRegisterFunc is a TranslationRegisterHandler generated by id.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/id.
func IdTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_id.RegisterDefaultTranslations
}

// JaTranslationRegisterFunc is a TranslationRegisterHandler generated by ja.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/ja.
func JaTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ja.RegisterDefaultTranslations
}

// NlTranslationRegisterFunc is a TranslationRegisterHandler generated by nl.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/nl.
func NlTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_nl.RegisterDefaultTranslations
}

// PtTranslationRegisterFunc is a TranslationRegisterHandler generated by pt.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/pt.
func PtTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_pt.RegisterDefaultTranslations
}

// PtBrTranslationRegisterFunc is a TranslationRegisterHandler generated by pt_BR.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/pt_BR.
func PtBrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_pt_BR.RegisterDefaultTranslations
}

// RuTranslationRegisterFunc is a TranslationRegisterHandler generated by ru.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translation/ru.
func RuTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ru.RegisterDefaultTranslations
}

// TrTranslationRegisterFunc is a TranslationRegisterHandler generated by tr.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/tr.
func TrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_tr.RegisterDefaultTranslations
}

// ZhTranslationRegisterFunc is a TranslationRegisterHandler generated by zh.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/zh.
func ZhTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_zh.RegisterDefaultTranslations
}

// ZhHantTranslationRegisterFunc is a TranslationRegisterHandler generated by zh_tw.RegisterDefaultTranslations from github.com/go-playground/validator/v10/translations/zh_tw.
func ZhHantTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_zh_tw.RegisterDefaultTranslations
}
