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

type TranslationRegisterHandler func(v *validator.Validate, trans ut.Translator) error

const (
	panicNilValidator             = "xvalidator: nil validator"
	panicNilLocaleTranslator      = "xvalidator: nil locale translator"
	panicNilTranslationRegisterFn = "xvalidator: nil translation register function"
)

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

	uniTranslator := ut.New(locTranslator, locTranslator)
	translator, _ := uniTranslator.GetTranslator(locTranslator.Locale()) // must found
	err := registerFn(validator, translator)                             // register translator to validator
	if err != nil {
		return nil, err
	}

	return translator, nil
}

func AddToTranslatorFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

func DefaultTranslateFunc() validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param()) // field to {0}, param to {1}
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
}

func EnLocaleTranslator() locales.Translator {
	return loc_en.New()
}

func FrLocaleTranslator() locales.Translator {
	return loc_fr.New()
}

func JaLocaleTranslator() locales.Translator {
	return loc_ja.New()
}

func RuLocaleTranslator() locales.Translator {
	return loc_ru.New()
}

func ZhLocaleTranslator() locales.Translator {
	return loc_zh.New()
}

func ZhHantLocaleTranslator() locales.Translator {
	return loc_zh_Hant.New()
}

func EnTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_en.RegisterDefaultTranslations
}

func FrTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_fr.RegisterDefaultTranslations
}

func JaTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ja.RegisterDefaultTranslations
}

func RuTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_ru.RegisterDefaultTranslations
}

func ZhTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_zh.RegisterDefaultTranslations
}

func ZhTwTranslationRegisterFunc() TranslationRegisterHandler {
	return trans_zh_tw.RegisterDefaultTranslations
}
