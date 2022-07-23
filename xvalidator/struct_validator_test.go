package xvalidator

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTranslateAndFlat(t *testing.T) {
	// WrappedFieldError
	v := validator.New()
	type s struct {
		Str string `validate:"required,gt=2,lt=10" json:"str"`
		Int int32  `validate:"required,gte=4,ne=5" json:"int"`
	}
	UseTagAsFieldName(v, "json")
	err1 := v.Struct(&s{}).(validator.ValidationErrors)[0]
	err2 := v.Struct(&s{Str: "abc"}).(validator.ValidationErrors)[0]
	err3 := v.Struct(&s{Str: "ab", Int: 4}).(validator.ValidationErrors)[0]
	err4 := v.Struct(&s{Str: "1234567890", Int: 4}).(validator.ValidationErrors)[0]
	err5 := v.Struct(&s{Str: "abc", Int: 3}).(validator.ValidationErrors)[0]
	err6 := v.Struct(&s{Str: "abc", Int: 5}).(validator.ValidationErrors)[0]

	err2_ := &WrappedFieldError{err2.(validator.FieldError), "Int field must be set and can not be zero"}
	err4_ := &WrappedFieldError{err4.(validator.FieldError), "The length of String must less than 10"}
	err6_ := &WrappedFieldError{err6.(validator.FieldError), "The value of Int must less than 5"}
	xtesting.Equal(t, err2_.Origin().Error(), err2.Error())
	xtesting.Equal(t, err2_.Unwrap().Error(), err2.Error())
	xtesting.Equal(t, errors.Is(err2_, err2), true)
	xtesting.Equal(t, err4_.Origin().Error(), "Key: 's.str' Error:Field validation for 'str' failed on the 'lt' tag")
	xtesting.Equal(t, err6_.Message(), "The value of Int must less than 5")
	xtesting.Equal(t, err6_.Error(), "Key: 's.int' Error:The value of Int must less than 5")

	// MultiFieldsError
	ve := &MultiFieldsError{fields: []error{err1, err2_}}
	xtesting.Equal(t, len(ve.Errors()), 2)
	xtesting.Equal(t, ve.Errors()[0].Error(), err1.Error())
	xtesting.Equal(t, ve.Errors()[1].Error(), "Key: 's.int' Error:Int field must be set and can not be zero")
	xtesting.Equal(t, ve.Error(), "Key: 's.str' Error:Field validation for 'str' failed on the 'required' tag; Key: 's.int' Error:Int field must be set and can not be zero")
	ve = &MultiFieldsError{fields: []error{err3, err5}}
	xtesting.Equal(t, ve.Error(), "Key: 's.str' Error:Field validation for 'str' failed on the 'gt' tag; Key: 's.int' Error:Field validation for 'int' failed on the 'gte' tag")
	ve = &MultiFieldsError{fields: []error{err4_, err6_}}
	xtesting.Equal(t, ve.Error(), "Key: 's.str' Error:The length of String must less than 10; Key: 's.int' Error:The value of Int must less than 5")

	// TranslateValidationErrors
	tr, _ := ApplyEnglishTranslator(v)
	fe := validator.ValidationErrors{err1, err2}
	fe1 := TranslateValidationErrors(fe, tr, true) // => fe.Translate(tr)
	fe2 := TranslateValidationErrors(fe, tr, false)
	xtesting.Equal(t, fe1["s.int"], "int is a required field")
	xtesting.Equal(t, fe1["s.str"], "str is a required field")
	xtesting.Equal(t, fe2["int"], "int is a required field")
	xtesting.Equal(t, fe2["str"], "str is a required field")
	xtesting.Panic(t, func() { TranslateValidationErrors(fe, nil, false) })

	// MultiFieldsError.Translate
	ve = &MultiFieldsError{fields: []error{err1, err2_}}
	ve1 := ve.Translate(tr, true)
	ve2 := ve.Translate(tr, false)
	xtesting.Equal(t, ve1["s.int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve1["s.str"], "str is a required field")
	xtesting.Equal(t, ve2["int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve2["str"], "str is a required field")
	xtesting.Panic(t, func() { ve.Translate(nil, false) })
	err := &MultiFieldsError{fields: []error{nil, errors.New("xxx")}}
	xtesting.Equal(t, len(err.Translate(tr, true)), 0)

	// FlatValidationErrors
	fe3 := FlatValidationErrors(fe, true)
	fe4 := FlatValidationErrors(fe, false)
	xtesting.Equal(t, fe3["s.int"], "Field validation for 'int' failed on the 'required' tag")
	xtesting.Equal(t, fe3["s.str"], "Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, fe4["int"], "Field validation for 'int' failed on the 'required' tag")
	xtesting.Equal(t, fe4["str"], "Field validation for 'str' failed on the 'required' tag")

	// MultiFieldsError.FlatToMap
	ve3 := ve.FlatToMap(true)
	ve4 := ve.FlatToMap(false)
	xtesting.Equal(t, ve3["s.int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve3["s.str"], "Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, ve4["int"], "Int field must be set and can not be zero")
	xtesting.Equal(t, ve4["str"], "Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, len(err.FlatToMap(true)), 0)

	// MapToError
	ferr := MapToError(nil)
	xtesting.Nil(t, ferr)
	ferr = MapToError(map[string]string{})
	xtesting.Nil(t, ferr)
	ferr = MapToError(fe2)
	xtesting.Equal(t, ferr.Error(), "int is a required field; str is a required field")
	ferr = MapToError(ve2)
	xtesting.Equal(t, ferr.Error(), "Int field must be set and can not be zero; str is a required field")
	ferr = MapToError(fe4)
	xtesting.Equal(t, ferr.Error(), "Field validation for 'int' failed on the 'required' tag; Field validation for 'str' failed on the 'required' tag")
	ferr = MapToError(ve4)
	xtesting.Equal(t, ferr.Error(), "Int field must be set and can not be zero; Field validation for 'str' failed on the 'required' tag")
}

func TestMessagedValidator(t *testing.T) {
	mv := NewMessagedValidator()
	mv.SetValidateTagName("binding")
	mv.SetMessageTagName("message")
	mv.UseTagAsFieldName("json", "yaml", "form")
	tr, _ := ApplyEnglishTranslator(mv.ValidateEngine())
	xtesting.Equal(t, mv.Engine(), mv.ValidateEngine())

	err := mv.ValidateStruct(time.Time{}) // validator.InvalidValidationError
	xtesting.NotNil(t, err)

	type s struct {
		Str   string  `binding:"required,gt=2,lt=10" json:"_str" message:"required|str must be set and can not be empty|gt|str length must be larger than 2|lt|str length must be less than 10"`
		Int   int32   `binding:"required,gte=4,ne=5" json:"_int" message:"required|int must be set and can not be zero|gte|int value must be larger than or equal to 4"`
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
			[]string{"Key: 's._str' Error:str length must be larger than 2"},
			map[string]string{"_str": "str length must be larger than 2"}},
		{&s{"1010101010", 4, -0.5}, false, true,
			[]string{"Key: 's._str' Error:str length must be less than 10"},
			map[string]string{"_str": "str length must be less than 10"}},
		{&s{"333", 0, -0.5}, false, true,
			[]string{"Key: 's._int' Error:int must be set and can not be zero"},
			map[string]string{"_int": "int must be set and can not be zero"}},
		{&s{"333", 3, -0.5}, false, true,
			[]string{"Key: 's._int' Error:int value must be larger than or equal to 4"},
			map[string]string{"_int": "int value must be larger than or equal to 4"}},
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
			[]string{"Key: 's._str' Error:str length must be less than 10", "Key: 's._int' Error:int value must be larger than or equal to 4"},
			map[string]string{"_str": "str length must be less than 10", "_int": "int value must be larger than or equal to 4"}},
		{&s{"", 5, 11.5}, false, true,
			[]string{"Key: 's._str' Error:str must be set and can not be empty", "Key: 's._int' Error:Field validation for '_int' failed on the 'ne' tag", "Key: 's.Float' Error:Field validation for 'Float' failed on the 'lte' tag"},
			map[string]string{"_str": "str must be set and can not be empty", "_int": "_int should not be equal to 5", "Float": "Float must be 0 or less"}},
	} {
		t.Run(fmt.Sprintf("%v", tc.giveObj), func(t *testing.T) {
			err := mv.ValidateStruct(tc.giveObj)
			_, ok := err.(*validator.InvalidValidationError)
			xtesting.Equal(t, ok, tc.wantInvalidErr)
			ve, ok := err.(*MultiFieldsError)
			xtesting.Equal(t, ok, tc.wantValidateErr)
			if ok && len(tc.wantFieldErrors) > 0 {
				xtesting.Equal(t, ve.Error(), strings.Join(tc.wantFieldErrors, "; "))
			}
			if ok && len(tc.wantTranslations) > 0 {
				xtesting.Equal(t, ve.Translate(tr, false), tc.wantTranslations)
			}
		})
	}
}

func TestApplyCustomMessage(t *testing.T) {
	v := NewMessagedValidator()
	type s struct {
		F1 string `validate_message:"required|required_1|gt|gt\\|1|lt|\\|lt\\|\\|_\\|\\|1\\|\\|\\|"`
		//                                                  ;gt  |1;lt;  |lt  |  |_  |  |1  |  |  |
		F2 string `validate_message:"x| |y| _|z|_ |w|_ _|lte|\\|_|gte|_\\|"`
		//                                               ;lte;  |_;gte;_  |
		F3 string `validate_message:"\\|eq|\\|_\\||_|\\||\\||_|\\|\\||\\|_\\|\\|_\\|"`
		//                              |eq;  |_  |;_;  |;  |;_;  |  |;  |_  |  |_  |
		F4 string `validate_message:""`
		F5 string
		F6 string `validate_message:"ne|ne_6|x"`
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
