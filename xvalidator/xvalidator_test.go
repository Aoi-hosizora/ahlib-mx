package xvalidator

import (
	"encoding/json"
	"errors"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"testing"
	"unsafe"
)

// =============
// for demo only
// =============

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

func TestFieldErrorMethods(t *testing.T) {
	// 1. and ,
	// 	A uint64 `binding:"required,lt=10|eq=12,ne=8"` // Tag       => "ne"
	// 	B string `binding:"required,range_1,ne=4444"`  // Tag       => "range_1" or "ne"
	// 	v.RegisterAlias("range_1", "gt=2,lt=10")       // ActualTag => "gt" or "ne"
	//
	// 2. or |
	// 	A uint64 `binding:"required,lt=10|eq=12,ne=8"`  // Tag       => "lt=10|eq=12"
	// 	C string `binding:"required,range_2,ne=666666"` // Tag       => "range_2" or "ne"
	// 	v.RegisterAlias("range_2", "lt=8|len=12")       // ActualTag => "lt=8|len=12"
	//
	// 3. Error()
	// 	Key: 'S.a' Error:Field validation for 'a' failed on the 'lt=10|eq=12' tag
	// 	Key: 'S.b' Error:Field validation for 'b' failed on the 'range_1' tag
	// 	Key: 'S.c' Error:Field validation for 'c' failed on the 'range_2' tag
	//       -----                            ---               ---------
	//     Namespace                         Field                 Tag
	//
	// 4. Translate()
	// 	...
	//
	// Also see https://godoc.org/github.com/go-playground/validator.

	v := validator.New()
	v.SetTagName("binding")
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})
	v.RegisterAlias("range_1", "gt=2,lt=10")
	v.RegisterAlias("range_2", "lt=8|len=12")
	tr, _ := ApplyTranslator(v, EnLocaleTranslator(), EnTranslationRegisterFunc())

	type S struct {
		A uint64 `binding:"required,lt=10|eq=12,ne=8"  json:"a"`
		B string `binding:"required,range_1,ne=4444"   json:"b"`
		C string `binding:"required,range_2,ne=666666" json:"c"`
	}
	for _, tc := range []struct {
		giveObj *S
		wantOk  bool
	}{
		{&S{3, "333", "333"}, true},
		// A
		{&S{11, "333", "333"}, false}, // Key: 'S.a' Error:Field validation for 'a' failed on the 'lt=10|eq=12' tag		| X
		{&S{12, "333", "333"}, true},  // 																					|
		{&S{8, "333", "333"}, false},  // Key: 'S.a' Error:Field validation for 'a' failed on the 'ne' tag					| a should not be equal to 8
		// B
		{&S{3, "1", "333"}, false},           // Key: 'S.b' Error:Field validation for 'b' failed on the 'range_1' tag		| X
		{&S{3, "11111111111", "333"}, false}, // Key: 'S.b' Error:Field validation for 'b' failed on the 'range_1' tag		| X
		{&S{3, "4444", "333"}, false},        // Key: 'S.b' Error:Field validation for 'b' failed on the 'ne' tag			| b should not be equal to 4444
		// C
		{&S{3, "333", "999999999"}, false},   // Key: 'S.c' Error:Field validation for 'c' failed on the 'range_2' tag		| X
		{&S{3, "333", "121212121212"}, true}, // 																			|
		{&S{3, "333", "666666"}, false},      // Key: 'S.c' Error:Field validation for 'c' failed on the 'ne' tag			| c should not be equal to 666666
	} {
		j, _ := json.Marshal(tc.giveObj)
		log.Println(string(j))
		err := v.Struct(tc.giveObj)
		xtesting.Equal(t, err == nil, tc.wantOk)
		if err != nil {
			e := err.(validator.ValidationErrors)[0]
			// log.Println(e.Field(), e.StructField())
			// log.Println(e.Namespace(), e.StructNamespace())
			log.Println(e.Tag(), e.ActualTag())
			log.Println(e.Param())
			// log.Println(e.Type(), e.Kind(), e.Value())
			log.Println(e.Error())
			log.Println(e.Translate(tr))
		}
		log.Println("=======================")
	}
}

// ========================
// tests start from here ->
// ========================

func TestIsXXXError(t *testing.T) {
	val := validator.New()
	type testStruct struct {
		Int int `validate:"required,lt=2"`
	}

	for _, tc := range []struct {
		giveErr   error
		wantValOk bool
		wantReqOk bool
	}{
		{nil, false, false},
		{errors.New("test"), false, false},
		{validator.ValidationErrors{}, true, false},
		{val.Struct(&testStruct{}), true, true},
		{val.Struct(&testStruct{Int: 0}), true, true},
		{val.Struct(&testStruct{Int: 1}), false, false},
		{val.Struct(&testStruct{Int: 3}), true, false},
	} {
		xtesting.Equal(t, IsValidationError(tc.giveErr), tc.wantValOk)
		xtesting.Equal(t, IsRequiredError(tc.giveErr), tc.wantReqOk)
	}
}

func TestUseTagAsFieldName(t *testing.T) {
	v1 := validator.New()
	v2 := validator.New()
	xtesting.Panic(t, func() { UseTagAsFieldName(nil, "a") })
	UseTagAsFieldName(v1, "json")
	UseTagAsFieldName(v2, "json")
	UseTagAsFieldName(v2, "") // unregister

	type S struct {
		F1 string  `validate:"required,gt=3" json:"f_1"`
		F2 int32   `validate:"required,ne=5" json:"f_2"`
		F3 float32 `validate:"required"      json:"-"`
		F4 float32 `validate:"required"`
	}
	for _, tc := range []struct {
		giveErr   validator.FieldError
		giveField string
		wantField string
		wantError string
	}{
		{v1.Struct(&S{}).(validator.ValidationErrors)[0], "F1", "f_1", "Key: 'S.f_1' Error:Field validation for 'f_1' failed on the 'required' tag"},
		{v2.Struct(&S{}).(validator.ValidationErrors)[0], "F1", "F1", "Key: 'S.F1' Error:Field validation for 'F1' failed on the 'required' tag"},
		{v1.Struct(&S{F1: "4444"}).(validator.ValidationErrors)[0], "F2", "f_2", "Key: 'S.f_2' Error:Field validation for 'f_2' failed on the 'required' tag"},
		{v2.Struct(&S{F1: "4444"}).(validator.ValidationErrors)[0], "F2", "F2", "Key: 'S.F2' Error:Field validation for 'F2' failed on the 'required' tag"},
		{v1.Struct(&S{F1: "333", F2: 4}).(validator.ValidationErrors)[0], "F1", "f_1", "Key: 'S.f_1' Error:Field validation for 'f_1' failed on the 'gt' tag"},
		{v2.Struct(&S{F1: "333", F2: 4}).(validator.ValidationErrors)[0], "F1", "F1", "Key: 'S.F1' Error:Field validation for 'F1' failed on the 'gt' tag"},
		{v1.Struct(&S{F1: "4444", F2: 5}).(validator.ValidationErrors)[0], "F2", "f_2", "Key: 'S.f_2' Error:Field validation for 'f_2' failed on the 'ne' tag"},
		{v2.Struct(&S{F1: "4444", F2: 5}).(validator.ValidationErrors)[0], "F2", "F2", "Key: 'S.F2' Error:Field validation for 'F2' failed on the 'ne' tag"},
		{v2.Struct(&S{F1: "4444", F2: 4, F3: 0, F4: 1.0}).(validator.ValidationErrors)[0], "F3", "F3", "Key: 'S.F3' Error:Field validation for 'F3' failed on the 'required' tag"},
		{v2.Struct(&S{F1: "4444", F2: 4, F3: 1.0, F4: 0}).(validator.ValidationErrors)[0], "F4", "F4", "Key: 'S.F4' Error:Field validation for 'F4' failed on the 'required' tag"},
	} {
		xtesting.Equal(t, tc.giveErr.Field(), tc.wantField)
		xtesting.Equal(t, tc.giveErr.Namespace(), "S."+tc.wantField)
		xtesting.Equal(t, tc.giveErr.StructField(), tc.giveField)
		xtesting.Equal(t, tc.giveErr.StructNamespace(), "S."+tc.giveField)
		xtesting.Equal(t, tc.giveErr.Error(), tc.wantError)
	}
}

func TestApplyTranslator(t *testing.T) {
	v := validator.New()
	type testStruct struct {
		String string `validate:"required,ne=hhh,lte=8,gt=1"`
	}

	for _, tc := range []struct {
		giveLoc     LocaleTranslator
		giveReg     TranslationRegisterHandler
		wantReqText string
		wantNeText  string
		wantLteText string
		wantGtText  string
	}{
		{EnLocaleTranslator(), EnTranslationRegisterFunc(),
			"String is a required field", "String should not be equal to hhh",
			"String must be at maximum 8 characters in length", "String must be greater than 1 character in length"},
		{FrLocaleTranslator(), FrTranslationRegisterFunc(),
			"String est un champ obligatoire", "String ne doit pas être égal à hhh",
			"String doit faire une taille maximum de 8 caractères", "String doit avoir une taille supérieur à 1 caractère"},
		{JaLocaleTranslator(), JaTranslationRegisterFunc(),
			"Stringは必須フィールドです", "Stringはhhhと異ならなければなりません",
			"Stringの長さは最大でも8文字でなければなりません", "Stringの長さは1文字よりも多くなければなりません"},
		// {RuLocaleTranslator(), RuTranslationRegisterFunc(),
		// 	"String обязательное поле", "String должен быть не равен hhh",
		// 	"", "String должен быть длиннее 1 символ"},
		{ZhLocaleTranslator(), ZhTranslationRegisterFunc(),
			"String为必填字段", "String不能等于hhh",
			"String长度不能超过8个字符", "String长度必须大于1个字符"},
		{ZhHantLocaleTranslator(), ZhHantTranslationRegisterFunc(),
			"String為必填欄位", "String不能等於hhh",
			"String長度不能超過8個字元", "String長度必須大於1個字元"},
	} {
		translator, err := ApplyTranslator(v, tc.giveLoc, tc.giveReg)
		xtesting.Nil(t, err)

		err = v.Struct(&testStruct{}) // required
		xtesting.NotNil(t, err)
		result := err.(validator.ValidationErrors).Translate(translator)
		xtesting.Equal(t, result["testStruct.String"], tc.wantReqText)

		err = v.Struct(&testStruct{String: "hhh"}) // ne
		xtesting.NotNil(t, err)
		result = err.(validator.ValidationErrors).Translate(translator)
		xtesting.Equal(t, result["testStruct.String"], tc.wantNeText)

		err = v.Struct(&testStruct{String: "999999999"}) // lte
		xtesting.NotNil(t, err)
		result = err.(validator.ValidationErrors).Translate(translator)
		xtesting.Equal(t, result["testStruct.String"], tc.wantLteText)

		err = v.Struct(&testStruct{String: "1"}) // gt
		xtesting.NotNil(t, err)
		result = err.(validator.ValidationErrors).Translate(translator)
		xtesting.Equal(t, result["testStruct.String"], tc.wantGtText)
	}

	xtesting.Panic(t, func() { _, _ = ApplyTranslator(nil, EnLocaleTranslator(), EnTranslationRegisterFunc()) })
	xtesting.Panic(t, func() { _, _ = ApplyTranslator(v, nil, EnTranslationRegisterFunc()) })
	xtesting.Panic(t, func() { _, _ = ApplyTranslator(v, EnLocaleTranslator(), nil) })
	errRegisterFn := func(v *validator.Validate, trans UtTranslator) error { return errors.New("test error") }
	_, err := ApplyTranslator(v, EnLocaleTranslator(), errRegisterFn)
	xtesting.Equal(t, err.Error(), "test error")
}

func TestRegisterTranslation(t *testing.T) {
	// 1. normal
	val := validator.New()
	type testStruct1 struct {
		String string `validate:"required"`
	}
	trans, _ := ApplyTranslator(val, EnLocaleTranslator(), EnTranslationRegisterFunc())
	fn := DefaultRegistrationFunc("required", "required {0}!!!", true)
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
	fn = DefaultRegistrationFunc("no_test", "translator for test tag", true)
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
	fn = DefaultRegistrationFunc("test", "{0} <- {1}", true)
	_ = val.RegisterTranslation("test", trans, fn, DefaultTranslateFunc())

	err = val.Struct(&testStruct3{}).(validator.ValidationErrors)
	xtesting.NotNil(t, err)
	transResults = err.Translate(trans)
	xtesting.Equal(t, transResults["testStruct3.String1"], "String1 <- ")
	xtesting.Equal(t, transResults["testStruct3.String2"], "String2 <- hhh")
}

func TestLocaleTranslatorsAndTranslationRegisterFunc(t *testing.T) {
	for _, tc := range []struct {
		give        LocaleTranslator
		wantName    string
		wantSunday  string
		wantJanuary string
	}{
		{EnLocaleTranslator(), "en", "Sunday", "January"},
		{EsLocaleTranslator(), "es", "domingo", "enero"},
		{FrLocaleTranslator(), "fr", "dimanche", "janvier"},
		{IdLocaleTranslator(), "id", "Minggu", "Januari"},
		{JaLocaleTranslator(), "ja", "日曜日", "1月"},
		{NlLocaleTranslator(), "nl", "zondag", "januari"},
		{PtLocaleTranslator(), "pt", "domingo", "janeiro"},
		{PtBrLocaleTranslator(), "pt_BR", "domingo", "janeiro"},
		{RuLocaleTranslator(), "ru", "воскресенье", "января"},
		{TrLocaleTranslator(), "tr", "Pazar", "Ocak"},
		{ZhLocaleTranslator(), "zh", "星期日", "一月"},
		{ZhHantLocaleTranslator(), "zh_Hant", "星期日", "1月"},
	} {
		xtesting.Equal(t, tc.give.Locale(), tc.wantName)
		xtesting.Equal(t, tc.give.WeekdaysWide()[0], tc.wantSunday)
		xtesting.Equal(t, tc.give.MonthsWide()[0], tc.wantJanuary)
	}

	type transText struct {
		text string
		// ...
	}
	for _, tc := range []struct {
		giveLoc          LocaleTranslator
		giveReg          TranslationRegisterHandler
		wantRequiredText string
		wantNotEqualText string
	}{
		{EnLocaleTranslator(), EnTranslationRegisterFunc(), "{0} is a required field", "{0} should not be equal to {1}"},
		{EsLocaleTranslator(), EsTranslationRegisterFunc(), "{0} es un campo requerido", "{0} no debería ser igual a {1}"},
		{FrLocaleTranslator(), FrTranslationRegisterFunc(), "{0} est un champ obligatoire", "{0} ne doit pas être égal à {1}"},
		{IdLocaleTranslator(), IdTranslationRegisterFunc(), "{0} wajib diisi", "{0} tidak sama dengan {1}"},
		{JaLocaleTranslator(), JaTranslationRegisterFunc(), "{0}は必須フィールドです", "{0}は{1}と異ならなければなりません"},
		{NlLocaleTranslator(), NlTranslationRegisterFunc(), "{0} is een verplicht veld", "{0} mag niet gelijk zijn aan {1}"},
		{PtLocaleTranslator(), PtTranslationRegisterFunc(), "{0} é obrigatório", "{0} não deve ser igual a {1}"},
		{PtLocaleTranslator(), PtBrTranslationRegisterFunc(), "{0} é um campo requerido", "{0} não deve ser igual a {1}"},
		{RuLocaleTranslator(), RuTranslationRegisterFunc(), "{0} обязательное поле", "{0} должен быть не равен {1}"},
		{TrLocaleTranslator(), TrTranslationRegisterFunc(), "{0} zorunlu bir alandır", "{0}, {1} değerine eşit olmamalıdır"},
		{ZhLocaleTranslator(), ZhTranslationRegisterFunc(), "{0}为必填字段", "{0}不能等于{1}"},
		{ZhHantLocaleTranslator(), ZhHantTranslationRegisterFunc(), "{0}為必填欄位", "{0}不能等於{1}"},
	} {
		val := validator.New()
		trans, err := ApplyTranslator(val, tc.giveLoc, tc.giveReg)
		xtesting.Nil(t, err)

		translations := xreflect.GetUnexportedField(reflect.ValueOf(trans).Elem().FieldByName("translations"))
		transTextPtr := translations.MapIndex(reflect.ValueOf("required")).Elem()
		xtesting.Equal(t, (*transText)(unsafe.Pointer(transTextPtr.UnsafeAddr())).text, tc.wantRequiredText)
		transTextPtr = translations.MapIndex(reflect.ValueOf("ne")).Elem()
		xtesting.Equal(t, (*transText)(unsafe.Pointer(transTextPtr.UnsafeAddr())).text, tc.wantNotEqualText)
	}
}
