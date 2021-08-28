package xvalidator

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
	"unsafe"
)

func TestApplyTranslator(t *testing.T) {
	v := validator.New()
	type testStruct struct {
		String string `validate:"required"`
	}

	for _, tc := range []struct {
		giveTranslator   locales.Translator
		giveRegisterFn   TranslationRegisterHandler
		wantRequiredText string
	}{
		{EnLocaleTranslator(), EnTranslationRegisterFunc(), "String is a required field"},
		{FrLocaleTranslator(), FrTranslationRegisterFunc(), "String est un champ obligatoire"},
		{JaLocaleTranslator(), JaTranslationRegisterFunc(), "Stringは必須フィールドです"},
		{ZhLocaleTranslator(), ZhTranslationRegisterFunc(), "String为必填字段"},
		{ZhHantLocaleTranslator(), ZhTwTranslationRegisterFunc(), "String為必填欄位"},
	} {
		translator, err := ApplyTranslator(v, tc.giveTranslator, tc.giveRegisterFn)
		xtesting.Nil(t, err)
		err = v.Struct(&testStruct{})
		xtesting.NotNil(t, err)
		xtesting.Equal(t, err.(validator.ValidationErrors).Translate(translator)["testStruct.String"], tc.wantRequiredText)
	}

	xtesting.Panic(t, func() { _, _ = ApplyTranslator(nil, EnLocaleTranslator(), EnTranslationRegisterFunc()) })
	xtesting.Panic(t, func() { _, _ = ApplyTranslator(v, nil, EnTranslationRegisterFunc()) })
	xtesting.Panic(t, func() { _, _ = ApplyTranslator(v, EnLocaleTranslator(), nil) })
	_, err := ApplyTranslator(v, EnLocaleTranslator(), func(v *validator.Validate, trans ut.Translator) error {
		return errors.New("test error")
	})
	xtesting.NotNil(t, err)
	xtesting.Equal(t, err.Error(), "test error")
}

func TestTranslationRegister(t *testing.T) {
	// 1. normal
	val := validator.New()
	type testStruct1 struct {
		String string `validate:"required"`
	}
	trans, _ := ApplyTranslator(val, EnLocaleTranslator(), EnTranslationRegisterFunc())
	fn := AddToTranslatorFunc("required", "required {0}!!!", true)
	_ = val.RegisterTranslation("required", trans, fn, DefaultTranslateFunc())

	err := val.Struct(&testStruct1{}).(validator.ValidationErrors)
	xtesting.NotNil(t, err)
	transResults := err.Translate(trans)
	xtesting.Equal(t, transResults["testStruct1.String"], "required String!!!")

	// 2. error with ut.T
	val = validator.New()
	_ = val.RegisterValidation("test", EqualValidator("test"))
	type testStruct2 struct {
		String string `validate:"test"`
	}
	trans, _ = ApplyTranslator(val, EnLocaleTranslator(), EnTranslationRegisterFunc())
	fn = AddToTranslatorFunc("no_test", "translator for test tag", true)
	_ = val.RegisterTranslation("test", trans, fn, DefaultTranslateFunc())

	err = val.Struct(&testStruct2{}).(validator.ValidationErrors)
	xtesting.NotNil(t, err)
	transResults = err.Translate(trans)
	xtesting.Equal(t, transResults["testStruct2.String"], "Key: 'testStruct2.String' Error:Field validation for 'String' failed on the 'test' tag")

	// 3. param with no param
	val = validator.New()
	_ = val.RegisterValidation("test", EqualValidator("test"))
	type testStruct3 struct {
		String1 string `validate:"test"`
		String2 string `validate:"test=hhh"`
	}
	trans, _ = ApplyTranslator(val, EnLocaleTranslator(), EnTranslationRegisterFunc())
	fn = AddToTranslatorFunc("test", "{0} <- {1}", true)
	_ = val.RegisterTranslation("test", trans, fn, DefaultTranslateFunc())

	err = val.Struct(&testStruct3{}).(validator.ValidationErrors)
	xtesting.NotNil(t, err)
	transResults = err.Translate(trans)
	xtesting.Equal(t, transResults["testStruct3.String1"], "String1 <- ")
	xtesting.Equal(t, transResults["testStruct3.String2"], "String2 <- hhh")
}

func TestLocaleTranslators(t *testing.T) {
	for _, tc := range []struct {
		give     locales.Translator
		wantName string
	}{
		{EnLocaleTranslator(), "en"},
		{FrLocaleTranslator(), "fr"},
		{IdLocaleTranslator(), "id"},
		{JaLocaleTranslator(), "ja"},
		{NlLocaleTranslator(), "nl"},
		{PtBrLocaleTranslator(), "pt_BR"},
		{RuLocaleTranslator(), "ru"},
		{TrLocaleTranslator(), "tr"},
		{ZhLocaleTranslator(), "zh"},
		{ZhHantLocaleTranslator(), "zh_Hant"},
	} {
		xtesting.Equal(t, tc.give.Locale(), tc.wantName)
	}
}

func TestTranslationRegisterFuncs(t *testing.T) {
	type transText struct {
		text    string
		indexes []int
	}

	for _, tc := range []struct {
		giveFn           TranslationRegisterHandler
		wantRequiredText string
	}{
		{EnTranslationRegisterFunc(), "{0} is a required field"},
		{FrTranslationRegisterFunc(), "{0} est un champ obligatoire"},
		{IdTranslationRegisterFunc(), "{0} wajib diisi"},
		{JaTranslationRegisterFunc(), "{0}は必須フィールドです"},
		{NlTranslationRegisterFunc(), "{0} is een verplicht veld"},
		{PtBrTranslationRegisterFunc(), "{0} é um campo requerido"},
		{RuTranslationRegisterFunc(), "{0} обязательное поле"},
		{TrTranslationRegisterFunc(), "{0} zorunlu bir alandır"},
		{ZhTranslationRegisterFunc(), "{0}为必填字段"},
		{ZhTwTranslationRegisterFunc(), "{0}為必填欄位"},
	} {
		val := validator.New()
		uniTrans := ut.New(EnLocaleTranslator(), EnLocaleTranslator())
		trans, _ := uniTrans.GetTranslator(EnLocaleTranslator().Locale())
		err := tc.giveFn(val, trans)
		xtesting.Nil(t, err)

		field := reflect.ValueOf(trans).Elem().FieldByName("translations")
		fieldValue := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		ptr := fieldValue.MapIndex(reflect.ValueOf("required"))
		xtesting.Equal(t, (*transText)(unsafe.Pointer(ptr.Elem().UnsafeAddr())).text, tc.wantRequiredText)
	}
}
