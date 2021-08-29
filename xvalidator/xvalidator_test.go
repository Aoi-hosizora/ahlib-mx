package xvalidator

import (
	"encoding/json"
	"errors"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
	"unsafe"
)

func TestRequiredAndOmitempty(t *testing.T) {
	// 1. `required` + non-pointer (common)
	// 	A uint64 `binding:"required"` // cannot be nil and 0
	// 	B string `binding:"required"` // cannot be nil and ""
	//
	// 2. `required` + pointer (common)
	// 	A *uint64 `binding:"required"` // cannot be nil, can be 0
	// 	B *string `binding:"required"` // cannot be nil, can be ""
	//
	// 3. `omitempty` + non-pointer (common)
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
	//
	// Also see https://godoc.org/github.com/go-playground/validator.

	v := validator.New()
	v.SetTagName("binding")

	type S1 struct {
		A uint64 `binding:"required"`
		B string `binding:"required"`
	}
	type S2 struct {
		A *uint64 `binding:"required"`
		B *string `binding:"required"`
	}
	type S3 struct {
		A uint64 `binding:"omitempty"`
		B string `binding:"omitempty"`
	}
	type S4 struct {
		A *uint64 `binding:"omitempty"`
		B *string `binding:"omitempty"`
	}
	type S5 struct {
		A uint64 `binding:"required,omitempty"`
		B string `binding:"required,omitempty"`
	}
	type S6 struct {
		A *uint64 `binding:"required,omitempty"`
		B *string `binding:"required,omitempty"`
	}

	for _, tc := range []struct {
		giveObj interface{}
		giveStr string
		wantOk  bool
	}{
		// typ required
		{&S1{}, `{}`, false},
		{&S1{}, `{"A": null, "B": null}`, false},
		{&S1{}, `{"A": 0, "B": ""}`, false},
		{&S1{}, `{"A": 1, "B": " "}`, true},
		// *typ required
		{&S2{}, `{}`, false},
		{&S2{}, `{"A": null, "B": null}`, false},
		{&S2{}, `{"A": 0, "B": ""}`, true},
		{&S2{}, `{"A": 1, "B": " "}`, true},
		// typ omitempty
		{&S3{}, `{}`, true},
		{&S3{}, `{"A": null, "B": null}`, true},
		{&S3{}, `{"A": 0, "B": ""}`, true},
		{&S3{}, `{"A": 1, "B": " "}`, true},
		// *typ omitempty => typ omitempty
		{&S4{}, `{}`, true},
		{&S4{}, `{"A": null, "B": null}`, true},
		{&S4{}, `{"A": 0, "B": ""}`, true},
		{&S4{}, `{"A": 1, "B": " "}`, true},
		// typ required,omitempty => typ required
		{&S5{}, `{}`, false},
		{&S5{}, `{"A": null, "B": null}`, false},
		{&S5{}, `{"A": 0, "B": ""}`, false},
		{&S5{}, `{"A": 1, "B": " "}`, true},
		// *typ required,omitempty => *typ required
		{&S6{}, `{}`, false},
		{&S6{}, `{"A": null, "B": null}`, false},
		{&S6{}, `{"A": 0, "B": ""}`, true},
		{&S6{}, `{"A": 1, "B": " "}`, true},
	} {
		_ = json.Unmarshal([]byte(tc.giveStr), tc.giveObj)
		xtesting.Equal(t, v.Struct(tc.giveObj) == nil, tc.wantOk)
	}
}

func TestIsXXXError(t *testing.T) {
	val := validator.New()
	type testStruct struct {
		Int int `validate:"required,lt=2"`
	}

	for _, tc := range []struct {
		giveErr   error
		wantReqOk bool
		wantValOk bool
	}{
		{nil, false, false},
		{errors.New("test"), false, false},
		{validator.ValidationErrors{}, false, true},
		{val.Struct(&testStruct{}), true, true},
		{val.Struct(&testStruct{Int: 0}), true, true},
		{val.Struct(&testStruct{Int: 1}), false, false},
		{val.Struct(&testStruct{Int: 3}), false, true},
	} {
		xtesting.Equal(t, IsValidationError(tc.giveErr), tc.wantValOk)
		xtesting.Equal(t, IsRequiredError(tc.giveErr), tc.wantReqOk)
	}
}

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
