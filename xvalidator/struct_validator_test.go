package xvalidator

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"testing"
)

func TestValidateFieldsErrorAndTranslate(t *testing.T) {
	// WrappedValidateFieldError
	v := validator.New()
	type s struct {
		Str string `validate:"required,gt=2,lt=10" json:"str"`
		Int int32  `validate:"required,gte=4,ne=5" json:"int"`
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string { return field.Tag.Get("json") })
	err1 := v.Struct(&s{}).(validator.ValidationErrors)[0]
	err2 := v.Struct(&s{Str: "abc"}).(validator.ValidationErrors)[0]
	err3 := v.Struct(&s{Str: "ab", Int: 4}).(validator.ValidationErrors)[0]
	err4 := v.Struct(&s{Str: "1234567890", Int: 4}).(validator.ValidationErrors)[0]
	err5 := v.Struct(&s{Str: "abc", Int: 3}).(validator.ValidationErrors)[0]
	err6 := v.Struct(&s{Str: "abc", Int: 5}).(validator.ValidationErrors)[0]

	err2_ := &WrappedValidateFieldError{err2.(validator.FieldError), "Int field must be set and can not be zero"}
	err4_ := &WrappedValidateFieldError{err4.(validator.FieldError), "The length of String must less then 10"}
	err6_ := &WrappedValidateFieldError{err6.(validator.FieldError), "The value of Int must less then 5"}
	xtesting.Equal(t, err2_.Origin().Error(), err2.Error())
	xtesting.Equal(t, err4_.Origin().Error(), "Key: 's.str' Error:Field validation for 'str' failed on the 'lt' tag")
	xtesting.Equal(t, err6_.Message(), "The value of Int must less then 5")
	xtesting.Equal(t, err6_.Error(), "Key: 's.int' Error:The value of Int must less then 5")

	// ValidateFieldsError
	ve := &ValidateFieldsError{fields: []error{err1, err2_}}
	xtesting.Equal(t, len(ve.Fields()), 2)
	xtesting.Equal(t, ve.Fields()[0].Error(), err1.Error())
	xtesting.Equal(t, ve.Fields()[1].Error(), "Key: 's.int' Error:Int field must be set and can not be zero")
	xtesting.Equal(t, ve.Error(), "Key: 's.str' Error:Field validation for 'str' failed on the 'required' tag\nKey: 's.int' Error:Int field must be set and can not be zero")
	ve = &ValidateFieldsError{fields: []error{err3, err5}}
	xtesting.Equal(t, ve.Error(), "Key: 's.str' Error:Field validation for 'str' failed on the 'gt' tag\nKey: 's.int' Error:Field validation for 'int' failed on the 'gte' tag")
	ve = &ValidateFieldsError{fields: []error{err4_, err6_}}
	xtesting.Equal(t, ve.Error(), "Key: 's.str' Error:The length of String must less then 10\nKey: 's.int' Error:The value of Int must less then 5")

	// TranslateValidationErrors
	tr, _ := ApplyTranslator(v, EnLocaleTranslator(), EnTranslationRegisterFunc())
	fe := validator.ValidationErrors{err1, err2}
	fe1 := TranslateValidationErrors(fe, tr, true) // => fe.Translate(tr)
	fe2 := TranslateValidationErrors(fe, tr, false)
	xtesting.Equal(t, fe1["s.int"], "int is a required field")
	xtesting.Equal(t, fe1["s.str"], "str is a required field")
	xtesting.Equal(t, fe2["int"], "int is a required field")
	xtesting.Equal(t, fe2["str"], "str is a required field")
	xtesting.Panic(t, func() { TranslateValidationErrors(fe, nil, false) })

	// ValidateFieldsError.Translate
	ve = &ValidateFieldsError{fields: []error{err1, err2_}}
	ve1 := ve.Translate(tr, true)
	ve2 := ve.Translate(tr, false)
	xtesting.Equal(t, ve1["s.int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve1["s.str"], "str is a required field")
	xtesting.Equal(t, ve2["int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve2["str"], "str is a required field")
	xtesting.Panic(t, func() { ve.Translate(nil, false) })
	err := &ValidateFieldsError{fields: []error{nil, errors.New("xxx")}}
	xtesting.Equal(t, len(err.Translate(tr, true)), 0)

	// SplitValidationErrors
	fe1 = SplitValidationErrors(fe, true)
	fe2 = SplitValidationErrors(fe, false)
	xtesting.Equal(t, fe1["s.int"], "Field validation for 'int' failed on the 'required' tag")
	xtesting.Equal(t, fe1["s.str"], "Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, fe2["int"], "Field validation for 'int' failed on the 'required' tag")
	xtesting.Equal(t, fe2["str"], "Field validation for 'str' failed on the 'required' tag")

	// ValidateFieldsError.SplitToMap
	ve1 = ve.SplitToMap(true)
	ve2 = ve.SplitToMap(false)
	xtesting.Equal(t, ve1["s.int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve1["s.str"], "Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, ve2["int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve2["str"], "Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, len(err.SplitToMap(true)), 0)
}

func TestCustomStructValidator(t *testing.T) {
	v := NewCustomStructValidator()
	v.SetValidatorTagName("binding")
	v.SetMessageTagName("message")
	UseTagAsFieldName(v.ValidateEngine(), "json")
	tr, _ := ApplyTranslator(v.ValidateEngine(), EnLocaleTranslator(), EnTranslationRegisterFunc())
	xtesting.Equal(t, v.Engine(), v.ValidateEngine())

	type s struct {
		Str   string  `binding:"required,gt=2,lt=10" json:"_str" message:"required|str must be set and can not be empty|gt|str length must be larger then 2|lt|str length must be less then 10"`
		Int   int32   `binding:"required,gte=4,ne=5" json:"_int" message:"required|int must be set and can not be zero|gte|int value must be larger then or equal to 4"`
		Float float32 `binding:"required,lte=0"`
	}
	for _, tc := range []struct {
		giveObj          interface{}
		wantInvalidErr   bool
		wantValidateErr  bool
		wantFieldErrors  []string
		wantTranslations map[string]string
	}{
		{nil, true, false, nil, nil},
		{new(*int), true, false, nil, nil},
		{"0", true, false, nil, nil},
		{0, true, false, nil, nil},
		{struct{}{}, false, false, nil, nil},
		{struct{ int }{1}, false, false, nil, nil},

		{&s{"333", 4, -0.5}, false, false, nil, nil},
		{&s{"", 4, -0.5}, false, true,
			[]string{"Key: 's._str' Error:str must be set and can not be empty"},
			map[string]string{"_str": "str must be set and can not be empty"}},
		{&s{"22", 4, -0.5}, false, true,
			[]string{"Key: 's._str' Error:str length must be larger then 2"},
			map[string]string{"_str": "str length must be larger then 2"}},
		{&s{"1010101010", 4, -0.5}, false, true,
			[]string{"Key: 's._str' Error:str length must be less then 10"},
			map[string]string{"_str": "str length must be less then 10"}},
		{&s{"333", 0, -0.5}, false, true,
			[]string{"Key: 's._int' Error:int must be set and can not be zero"},
			map[string]string{"_int": "int must be set and can not be zero"}},
		{&s{"333", 3, -0.5}, false, true,
			[]string{"Key: 's._int' Error:int value must be larger then or equal to 4"},
			map[string]string{"_int": "int value must be larger then or equal to 4"}},
		{&s{"333", 5, -0.5}, false, true,
			[]string{"Key: 's._int' Error:Field validation for '_int' failed on the 'ne' tag"},
			map[string]string{"_int": "_int should not be equal to 5"}},
		{&s{"333", 4, 0}, false, true,
			[]string{"Key: 's.Float' Error:Field validation for 'Float' failed on the 'required' tag"},
			map[string]string{"Float": "Float is a required field"}},
		{&s{"333", 4, 1.0}, false, true,
			[]string{"Key: 's.Float' Error:Field validation for 'Float' failed on the 'lte' tag"},
			map[string]string{"Float": "Float must be 0 or less"}},
		{&s{"11111111111", -1, -0.5}, false, true,
			[]string{"Key: 's._str' Error:str length must be less then 10", "Key: 's._int' Error:int value must be larger then or equal to 4"},
			map[string]string{"_str": "str length must be less then 10", "_int": "int value must be larger then or equal to 4"}},
		{&s{"", 5, 11.5}, false, true,
			[]string{"Key: 's._str' Error:str must be set and can not be empty", "Key: 's._int' Error:Field validation for '_int' failed on the 'ne' tag", "Key: 's.Float' Error:Field validation for 'Float' failed on the 'lte' tag"},
			map[string]string{"_str": "str must be set and can not be empty", "_int": "_int should not be equal to 5", "Float": "Float must be 0 or less"}},
	} {
		t.Run(fmt.Sprintf("%v", tc.giveObj), func(t *testing.T) {
			err := v.ValidateStruct(tc.giveObj)
			_, ok := err.(*validator.InvalidValidationError)
			xtesting.Equal(t, ok, tc.wantInvalidErr)
			ve, ok := err.(*ValidateFieldsError)
			xtesting.Equal(t, ok, tc.wantValidateErr)
			if ok && len(tc.wantFieldErrors) > 0 {
				xtesting.Equal(t, ve.Error(), strings.Join(tc.wantFieldErrors, "\n"))
			}
			if ok && len(tc.wantTranslations) > 0 {
				xtesting.Equal(t, ve.Translate(tr, false), tc.wantTranslations)
			}
		})
	}
}

func TestApplyCustomMessage(t *testing.T) {
	v := NewCustomStructValidator()
	type s struct {
		F1 string `validator_message:"required|required_1|gt|gt\\|1|lt|\\|lt\\|\\|_\\|\\|1\\|\\|\\|"`
		//                                                  ;gt  |1;lt;  |lt  |  |_  |  |1  |  |  |
		F2 string `validator_message:"x| |y| _|z|_ |w|_ _|lte|\\|_|gte|_\\|"`
		//                                               ;lte;  |_;gte;_  |
		F3 string `validator_message:"\\|eq|\\|_\\||_|\\||\\||_|\\|\\||\\|_\\|\\|_\\|"`
		//                              |eq;  |_  |;_;  |;  |;_;  |  |;  |_  |  |_  |
		F4 string `validator_message:""`
		F5 string
		F6 string `validator_message:"ne|ne_6|x"`
	}

	for _, tc := range []struct {
		giveField string
		giveTag   string
		wantOk    bool
		wantMsg   string
	}{
		{"F1", "required", true, "required_1"},
		{"F1", "required_", false, ""},
		{"F1", "gt", true, "gt|1"},
		{"F1", "lt", true, "|lt||_||1|||"},
		{"F1", "", false, ""},
		{"F1", "gte", false, ""},
		{"F2", "x", false, ""},
		{"F2", "y", true, "_"},
		{"F2", "z", true, "_"},
		{"F2", "w", true, "_ _"},
		{"F2", "lte", true, "|_"},
		{"F2", "gte", true, "_|"},
		{"F3", "eq", false, ""},
		{"F3", "|eq", true, "|_|"},
		{"F3", "_", true, "|"},
		{"F3", "|", true, "_"},
		{"F3", "||", true, "|_||_|"},
		{"F4", "x", false, ""},
		{"F5", "x", false, ""},
		{"F6", "ne", true, "ne_6"},
		{"F6", "x", false, ""},
		{"F6", "y", false, ""},
		{"F?", "x", false, ""},
	} {
		t.Run(tc.giveField+"_"+tc.giveTag, func(t *testing.T) {
			msg, ok := v.applyCustomMessage(reflect.TypeOf(s{}), tc.giveField, tc.giveTag)
			xtesting.Equal(t, ok, tc.wantOk)
			xtesting.Equal(t, msg, tc.wantMsg)
		})
	}

}