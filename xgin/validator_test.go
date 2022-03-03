package xgin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type mockValidator struct{}

func (f mockValidator) ValidateStruct(interface{}) error { return nil }
func (f mockValidator) Engine() interface{}              { return nil } // fake





func TestValidator(t *testing.T) {
	val, err := GetValidatorEngine() // gin's default validator.Validator
	xtesting.Nil(t, err)

	// default
	type testStruct1 struct {
		String string `binding:"required"`
	}
	xtesting.NotNil(t, val.Struct(&testStruct1{}))
	xtesting.Nil(t, val.Struct(&testStruct1{String: "xxx"}))

	// custom tag
	val.SetTagName("validate") // let the global validator use `validate` tag
	type testStruct2 struct {
		String string `validate:"required,gt=2"`
	}
	xtesting.NotNil(t, val.Struct(&testStruct2{}))             // err: required
	xtesting.NotNil(t, val.Struct(&testStruct2{String: "xx"})) // err: gt
	xtesting.Nil(t, val.Struct(&testStruct2{String: "xxx"}))
	val.SetTagName("binding") // default to `binding`

	// custom field name
	type testStruct3 struct {
		String string `binding:"required,ne=hhh" json:"str"`
	}
	xvalidator.UseTagAsFieldName(val, "json") // let the global validator use `json` as field name
	xtesting.Equal(t, val.Struct(&testStruct3{}).Error(), "Key: 'testStruct3.str' Error:Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, val.Struct(&testStruct3{String: "hhh"}).Error(), "Key: 'testStruct3.str' Error:Field validation for 'str' failed on the 'ne' tag")
	xvalidator.UseTagAsFieldName(val, "") // defaults to no use
	xtesting.Equal(t, val.Struct(&testStruct3{}).Error(), "Key: 'testStruct3.str' Error:Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, val.Struct(&testStruct3{String: "hhh"}).Error(), "Key: 'testStruct3.str' Error:Field validation for 'str' failed on the 'ne' tag") // <<< with structCache
	type testStruct4 struct {
		String string `binding:"required,lt=3" json:"str"`
	}
	xtesting.Equal(t, val.Struct(&testStruct4{}).Error(), "Key: 'testStruct4.String' Error:Field validation for 'String' failed on the 'required' tag")
	xtesting.Equal(t, val.Struct(&testStruct4{String: "hhh"}).Error(), "Key: 'testStruct4.String' Error:Field validation for 'String' failed on the 'lt' tag")
}

func TestTranslator(t *testing.T) {
	val, err := GetValidatorEngine()
	xtesting.Nil(t, err)
	ut, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)

	// global
	xtesting.Nil(t, GetGlobalTranslator())
	SetGlobalTranslator(ut)
	xtesting.Equal(t, GetGlobalTranslator(), ut)
	SetGlobalTranslator(nil)
	xtesting.Nil(t, GetGlobalTranslator())

	// translate
	type testStruct struct {
		String string `binding:"required,ne=hhh" json:"str"`
	}
	xvalidator.UseTagAsFieldName(val, "json")
	defer func() { xvalidator.UseTagAsFieldName(val, "") }()
	for _, tc := range []struct {
		giveLoc          xvalidator.LocaleTranslator
		giveReg          xvalidator.TranslationRegisterHandler
		wantRequiredText string
		wantNotEqualText string
	}{
		{xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc(), "str is a required field", "str should not be equal to hhh"},
		{xvalidator.FrLocaleTranslator(), xvalidator.FrTranslationRegisterFunc(), "str est un champ obligatoire", "str ne doit pas être égal à hhh"},
		{xvalidator.JaLocaleTranslator(), xvalidator.JaTranslationRegisterFunc(), "strは必須フィールドです", "strはhhhと異ならなければなりません"},
		{xvalidator.RuLocaleTranslator(), xvalidator.RuTranslationRegisterFunc(), "str обязательное поле", "str должен быть не равен hhh"},
		{xvalidator.ZhLocaleTranslator(), xvalidator.ZhTranslationRegisterFunc(), "str为必填字段", "str不能等于hhh"},
		{xvalidator.ZhHantLocaleTranslator(), xvalidator.ZhHantTranslationRegisterFunc(), "str為必填欄位", "str不能等於hhh"},
	} {
		ut, err := GetValidatorTranslator(tc.giveLoc, tc.giveReg)
		xtesting.Nil(t, err)
		reqErr := val.Struct(&testStruct{}).(validator.ValidationErrors)
		neErr := val.Struct(&testStruct{String: "hhh"}).(validator.ValidationErrors)
		xtesting.Equal(t, reqErr.Translate(ut)["testStruct.str"], tc.wantRequiredText)
		xtesting.Equal(t, neErr.Translate(ut)["testStruct.str"], tc.wantNotEqualText)
		xtesting.Equal(t, xvalidator.TranslateValidationErrors(reqErr, ut, false)["str"], tc.wantRequiredText)
		xtesting.Equal(t, xvalidator.TranslateValidationErrors(neErr, ut, false)["str"], tc.wantNotEqualText)
	}

	t.Run("mismatch validator type", func(t *testing.T) {
		// mismatched validator engine
		originVal := binding.Validator
		originTrans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
		binding.Validator = &mockValidator{}
		_, err = GetValidatorEngine()
		xtesting.NotNil(t, err)
		_, err = GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
		xtesting.NotNil(t, err)
		xtesting.NotNil(t, EnableParamRegexpBinding())
		xtesting.NotNil(t, EnableRFC3339DateBinding())
		xtesting.NotNil(t, EnableParamRegexpBindingTranslator(originTrans))
		xtesting.NotNil(t, EnableRFC3339DateBindingTranslator(originTrans))

		binding.Validator = originVal
		_, err = GetValidatorEngine()
		xtesting.Nil(t, err)
		trans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
		xtesting.Nil(t, err)
		xtesting.Nil(t, EnableRFC3339DateTimeBinding())
		xtesting.Nil(t, EnableRFC3339DateTimeBindingTranslator(trans))
	})
}

func TestAddBindingAndAddTranslator(t *testing.T) {
	val, _ := GetValidatorEngine()
	xvalidator.UseTagAsFieldName(val, "json")
	trans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())

	xtesting.Nil(t, EnableParamRegexpBinding())
	xtesting.Nil(t, EnableParamRegexpBindingTranslator(trans))
	xtesting.Nil(t, EnableRFC3339DateBinding())
	xtesting.Nil(t, EnableRFC3339DateBindingTranslator(trans))
	xtesting.Nil(t, EnableRFC3339DateTimeBinding())
	xtesting.Nil(t, EnableRFC3339DateTimeBindingTranslator(trans))
	xtesting.Nil(t, AddBinding("re_number", xvalidator.RegexpValidator(regexp.MustCompile(`^[0-9]+$`))))
	xtesting.Nil(t, AddTranslation(trans, "re_number", "{0} should be a number string", true))
	xtesting.Nil(t, AddBinding("range_name", xvalidator.LengthInRangeValidator(3, 10)))
	xtesting.Nil(t, AddTranslation(trans, "range_name", "{0} should be in range of [3, 10]", true))

	type testStruct struct {
		Number   string `json:"number"   form:"number"   binding:"required,re_number"`
		Abc      string `json:"abc"      form:"abc"      binding:"required,regexp=^[abc]+$"`
		Date     string `json:"date"     form:"date"     binding:"required,date"`
		Datetime string `json:"datetime" form:"datetime" binding:"required,datetime"`
		Username string `json:"username" form:"username" binding:"required,range_name"`
		Num      *int32 `json:"num"      form:"num"      binding:"required,lt=2"`
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("", func(ctx *gin.Context) {
		test := &testStruct{}
		err := ctx.ShouldBindQuery(test)
		if err != nil {
			translations := xvalidator.TranslateValidationErrors(err.(validator.ValidationErrors), trans, false)
			ctx.JSON(400, gin.H{"success": false, "details": translations})
		} else {
			ctx.JSON(200, gin.H{"success": true})
		}
	})
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		giveMap     map[string]interface{}
		wantSuccess bool
		wantMap     map[string]interface{}
	}{
		{nil, false,
			map[string]interface{}{
				"number":   "number is a required field",
				"abc":      "abc is a required field",
				"date":     "date is a required field",
				"datetime": "datetime is a required field",
				"username": "username is a required field",
				"num":      "num is a required field",
			},
		},
		{map[string]interface{}{"number": "", "abc": "", "date": "", "datetime": "", "username": "   ", "num": 0}, false,
			map[string]interface{}{
				"number":   "number is a required field",
				"abc":      "abc is a required field",
				"date":     "date is a required field",
				"datetime": "datetime is a required field",
			},
		},
		{map[string]interface{}{"number": "a", "abc": "d", "date": "2021/02/03", "datetime": "2021/02/03 02:10:13", "username": "u", "num": 5}, false,
			map[string]interface{}{
				"number":   "number should be a number string",
				"abc":      "abc should match regexp /^[abc]+$/",
				"date":     "date should be an RFC3339 date",
				"datetime": "datetime should be an RFC3339 datetime",
				"username": "username should be in range of [3, 10]",
				"num":      "num must be less than 2",
			},
		},
		{map[string]interface{}{"number": "", "abc": "a", "date": "2021-2-3", "datetime": "2021-02-03T02:10:13", "username": "11111111111"}, false,
			map[string]interface{}{
				"number":   "number is a required field",
				"date":     "date should be an RFC3339 date",
				"datetime": "datetime should be an RFC3339 datetime",
				"username": "username should be in range of [3, 10]",
				"num":      "num is a required field",
			},
		},
		{map[string]interface{}{"number": "0", "abc": "abcd", "date": "2021/02/03", "datetime": "", "username": "55555", "num": 0}, false,
			map[string]interface{}{
				"abc":      "abc should match regexp /^[abc]+$/",
				"date":     "date should be an RFC3339 date",
				"datetime": "datetime is a required field",
			},
		},
		{map[string]interface{}{"number": "123", "abc": "aabbcc", "date": "2021-02-03", "datetime": "2021-02-03T02:10:13+08:00", "username": "333", "num": -1}, true,
			nil,
		},
	} {
		value := url.Values{}
		for k, v := range tc.giveMap {
			value[k] = []string{fmt.Sprintf("%v", v)}
		}
		query := value.Encode()
		t.Run(query, func(t *testing.T) {
			resp, err := http.Get("http://localhost:12345?" + query)
			xtesting.Nil(t, err)
			bs, _ := ioutil.ReadAll(resp.Body)
			r := make(map[string]interface{})
			err = json.Unmarshal(bs, &r)
			xtesting.Nil(t, err)
			xtesting.Equal(t, r["success"].(bool), tc.wantSuccess)
			if !tc.wantSuccess {
				xtesting.Equal(t, r["details"].(map[string]interface{}), tc.wantMap)
			}
		})
	}
}

func TestCustomStructValidator(t *testing.T) {
	originVal := binding.Validator
	mv := xvalidator.NewMessagedValidator()
	mv.SetValidateTagName("binding")
	mv.SetMessageTagName("message")
	mv.UseTagAsFieldName("json")
	binding.Validator = mv
	defer func() { binding.Validator = originVal }()

	val, err := GetValidatorEngine()
	xtesting.Nil(t, err)
	trans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)
	xvalidator.UseTagAsFieldName(val, "json")
	xtesting.Nil(t, EnableParamRegexpBinding())
	xtesting.Nil(t, EnableParamRegexpBindingTranslator(trans))
	xtesting.Nil(t, AddBinding("re_number", xvalidator.RegexpValidator(regexp.MustCompile(`^[0-9]+$`))))
	xtesting.Nil(t, AddTranslation(trans, "re_number", "{0} should be a number string", true))
	xtesting.Nil(t, AddBinding("rg_int", xvalidator.Or(xvalidator.EqualValidator(0), xvalidator.And(xvalidator.GreaterThenOrEqualValidator(4), xvalidator.NotEqualValidator(5)))))

	type testStruct struct {
		String string  `binding:"required,re_number" json:"_string" form:"_string" message:"required|_string should be set and can not be empty"`
		Int    *int32  `binding:"required,rg_int"    json:"_int"    form:"_int"    message:"rg_int|_int should be in specific range"`
		Float  float64 `binding:"required,gt=0.3"    json:"_float"  form:"_float"`
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("", func(ctx *gin.Context) {
		s := &testStruct{}
		if err := ctx.ShouldBindQuery(s); err != nil {
			// if s.Int != nil {
			// 	log.Print(*s.Int)
			// }
			if verr, ok := err.(*xvalidator.MultiFieldsError); ok {
				ctx.JSON(400, gin.H{"success": false, "details": verr.Translate(trans, false)})
			} else if nerr, ok := err.(*strconv.NumError); ok {
				ctx.JSON(400, gin.H{"success": false, "details": gin.H{"__json": nerr.Error()}})
			} else {
				ctx.JSON(400, gin.H{"success": false, "details": "?"})
			}
		} else {
			ctx.JSON(200, gin.H{"success": true})
		}
	})
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		giveMap     map[string]interface{}
		wantSuccess bool
		wantMap     map[string]interface{}
	}{
		{nil, false,
			map[string]interface{}{"_string": "_string should be set and can not be empty", "_int": "_int is a required field", "_float": "_float is a required field"},
		},
		{map[string]interface{}{"_int": "0"}, false,
			map[string]interface{}{"_string": "_string should be set and can not be empty", "_float": "_float is a required field"},
		},
		{map[string]interface{}{"_string": "", "_int": "", "_float": ""}, false,
			map[string]interface{}{"_string": "_string should be set and can not be empty", "_float": "_float is a required field"}, // _int == 0 in query
		},
		{map[string]interface{}{"_string": "", "_int": "_", "_float": ""}, false,
			map[string]interface{}{"__json": "strconv.ParseInt: parsing \"_\": invalid syntax"},
		},
		{map[string]interface{}{"_string": "a0", "_int": " ", "_float": "0.0"}, false,
			map[string]interface{}{"__json": "strconv.ParseInt: parsing \" \": invalid syntax"},
		},
		{map[string]interface{}{"_string": " 0", "_int": "3", "_float": "0.3"}, false,
			map[string]interface{}{"_string": "_string should be a number string", "_int": "_int should be in specific range", "_float": "_float must be greater than 0.3"},
		},
		{map[string]interface{}{"_string": "1", "_int": "5.1", "_float": "1"}, false,
			map[string]interface{}{"__json": "strconv.ParseInt: parsing \"5.1\": invalid syntax"},
		},
		{map[string]interface{}{"_string": "1.2", "_int": "5", "_float": "0_1"}, false,
			map[string]interface{}{"_string": "_string should be a number string", "_int": "_int should be in specific range"},
		},
		{map[string]interface{}{"_string": "abc", "_int": "0", "_float": "0"}, false,
			map[string]interface{}{"_string": "_string should be a number string", "_float": "_float is a required field"},
		},
		{map[string]interface{}{"_string": "123", "_int": "4", "_float": "1.1"}, true, nil},
	} {
		value := url.Values{}
		for k, v := range tc.giveMap {
			value[k] = []string{fmt.Sprintf("%v", v)}
		}
		query := value.Encode()
		t.Run(query, func(t *testing.T) {
			resp, err := http.Get("http://localhost:12345?" + query)
			xtesting.Nil(t, err)
			bs, _ := ioutil.ReadAll(resp.Body)

			r := make(map[string]interface{})
			err = json.Unmarshal(bs, &r)
			xtesting.Nil(t, err)
			xtesting.Equal(t, r["success"].(bool), tc.wantSuccess)
			if !tc.wantSuccess {
				xtesting.Equal(t, r["details"].(map[string]interface{}), tc.wantMap)
			}
		})
	}
}

type translatableError string

func (t translatableError) Error() string {
	return string(t)
}

func (t translatableError) Translate() (map[string]string, bool) {
	return map[string]string{"_": string(t)}, true
}

func TestTranslateBindingError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	val, _ := GetValidatorEngine()
	xvalidator.UseTagAsFieldName(val, "json")
	trans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	mv := xvalidator.NewMessagedValidator()
	mv.SetValidateTagName("binding")
	mv.SetMessageTagName("message")
	mv.UseTagAsFieldName("json")

	app := gin.New()
	type testStruct struct {
		Str string `json:"str" form:"str" binding:"required" message:"required|str should be not null and not empty"`
		Int int32  `json:"int" form:"int" binding:"required" message:"required|int should be not null and not zero"`
	}
	respond := func(c *gin.Context, code int, details map[string]string) {
		if code == 200 {
			c.JSON(200, gin.H{"success": true})
		} else {
			c.JSON(code, gin.H{"success": false, "details": details})
		}
	}

	app.POST("/body", func(c *gin.Context) {
		opts := make([]TranslateOption, 0)
		var ptr interface{} = &testStruct{}
		if c.Query("useInvalidType") == "true" {
			ptr = 0
		}
		if c.Query("useTrans") == "true" {
			opts = append(opts, WithUtTranslator(trans))
		}
		if c.Query("useCustom") == "true" {
			originVal := binding.Validator
			binding.Validator = mv
			defer func() { binding.Validator = originVal }()
		}
		if err := c.ShouldBind(ptr); err != nil {
			if result, need4xx := TranslateBindingError(err, opts...); need4xx {
				respond(c, 400, result)
			} else {
				respond(c, 500, result)
			}
		} else {
			respond(c, 200, nil)
		}
	})
	app.POST("/id/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			if c.Query("useNumError") != "true" {
				err = NewRouterDecodeError("id", idStr, err, "")
			}
		} else if id <= 0 {
			err = errors.New("id <= 0")
			err = NewRouterDecodeError("id", idStr, err, "should be larger then zero")
		}
		if rErr, ok := err.(*RouterDecodeError); ok {
			_ = rErr.Error()
			if c.Query("ignoreField") == "true" {
				rErr.Field = ""
			}
			if c.Query("ignoreMessage") == "true" {
				rErr.Message = ""
			}
		}
		if err != nil {
			if result, need4xx := TranslateBindingError(err); need4xx {
				respond(c, 400, result)
			} else {
				respond(c, 500, result)
			}
		} else {
			respond(c, 200, nil)
		}
	})

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	// normal
	for _, tc := range []struct {
		giveRoute string
		giveBody  string
		giveQuery string
		wantCode  int
		wantMap   map[string]interface{}
	}{
		// io eof
		{"body", ``, "", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position -1"},
		},
		{"body", `{`, "", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position -1"},
		},
		// json invalid unmarshal
		{"body", `{}`, "useInvalidType=true", 500, nil},
		// json syntax
		{"body", `{"str": a, "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position 9"},
		},
		{"body", `{"str": "a, "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position 14"},
		},
		// json type
		{"body", `{"str": 0, "int": 0}`, "", 400,
			map[string]interface{}{"__decode": "type of 'number' in 'str' mismatches with required 'string'"},
		},
		{"body", `{"str": "", "int": ""}`, "", 400,
			map[string]interface{}{"__decode": "type of 'string' in 'int' mismatches with required 'int32'"},
		},
		{"body", `{"str": "abc", "int": 999999999999999999999999999999}`, "", 400,
			map[string]interface{}{"__decode": "type of 'number 999999999999999999999999999999' in 'int' mismatches with required 'int32'"},
		},
		{"body", `{"str": "abc", "int": 3.14}`, "", 400,
			map[string]interface{}{"__decode": "type of 'number 3.14' in 'int' mismatches with required 'int32'"},
		},
		// validator
		{"body", `{}`, "", 400,
			map[string]interface{}{"str": "Field validation for 'str' failed on the 'required' tag", "int": "Field validation for 'int' failed on the 'required' tag"},
		},
		{"body", `{}`, "useTrans=true", 400,
			map[string]interface{}{"str": "str is a required field", "int": "int is a required field"},
		},
		{"body", `{"str": "", "int": 0}`, "", 400,
			map[string]interface{}{"str": "Field validation for 'str' failed on the 'required' tag", "int": "Field validation for 'int' failed on the 'required' tag"},
		},
		{"body", `{"str": "", "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"str": "str is a required field", "int": "int is a required field"},
		},
		// xvalidator required
		{"body", `{}`, "useCustom=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{}`, "useCustom=true&useTrans=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{"str": "", "int": 0}`, "useCustom=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{"str": "", "int": 0}`, "useCustom=true&useTrans=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		// ok
		{"body", `{"str": "abc", "int": 1}`, "", 200, nil},
		// strconv number error
		{"id/a", ``, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter must be a number"},
		},
		{"id/3.14", ``, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter must be a number"},
		},
		{"id/999999999999999999999999999999", `useNumError=true`, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter is out of range"},
		},
		// router decode error
		{"id/a", ``, "", 400,
			map[string]interface{}{"id": "router parameter id must be a number"},
		},
		{"id/3.14", ``, "", 400,
			map[string]interface{}{"id": "router parameter id must be a number"},
		},
		{"id/3.14", ``, "ignoreField=true", 400,
			map[string]interface{}{"router parameter": "router parameter must be a number"},
		},
		{"id/999999999999999999999999999999", ``, "", 400,
			map[string]interface{}{"id": "router parameter id is out of range"},
		},
		{"id/999999999999999999999999999999", ``, "ignoreMessage=true", 400,
			map[string]interface{}{"id": "router parameter id is out of range"},
		},
		{"id/0", ``, "", 400,
			map[string]interface{}{"id": "router parameter id should be larger then zero"},
		},
		{"id/0", ``, "ignoreField=true", 400,
			map[string]interface{}{"router parameter": "router parameter should be larger then zero"},
		},
		{"id/0", ``, "ignoreMessage=true", 500, nil},
		// ok
		{"id/1", ``, "", 200, nil},
	} {
		t.Run(tc.giveRoute+"_"+tc.giveBody, func(t *testing.T) {
			u := "http://localhost:12345/" + tc.giveRoute + "?" + tc.giveQuery
			req, _ := http.NewRequest("POST", u, strings.NewReader(tc.giveBody))
			req.Header.Set("Content-type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			xtesting.Nil(t, err)
			bs, _ := ioutil.ReadAll(resp.Body)

			r := make(map[string]interface{})
			err = json.Unmarshal(bs, &r)
			xtesting.Nil(t, err)
			xtesting.Equal(t, resp.StatusCode, tc.wantCode)
			if resp.StatusCode != 200 && r["details"] != nil {
				xtesting.Equal(t, r["details"].(map[string]interface{}), tc.wantMap)
			}
		})
	}

	// for coverage
	for _, tc := range []struct {
		name     string
		giveErr  error
		giveOpts []TranslateOption
		wantMap  map[string]string
		want4xx  bool
	}{
		{"nil", nil, nil, nil, false},
		{"NumError", &strconv.NumError{}, nil, nil, false},
		{"InvalidValidationError", &validator.InvalidValidationError{}, nil, nil, false},
		{"ExtraErrors", errors.New("TODO"), nil, nil, false},
		{"InvalidUnmarshalError", &json.InvalidUnmarshalError{}, []TranslateOption{WithJsonInvalidUnmarshalError(
			func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool) { return nil, true })}, nil, true},
		{"UnmarshalTypeError", &json.UnmarshalTypeError{}, []TranslateOption{WithJsonUnmarshalTypeError(
			func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"SyntaxError", &json.SyntaxError{}, []TranslateOption{WithJsonSyntaxError(
			func(*json.SyntaxError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"ErrUnexpectedEOF", io.ErrUnexpectedEOF, []TranslateOption{WithIoEOFError(
			func(error) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"NumError", &strconv.NumError{}, []TranslateOption{WithStrconvNumErrorError(
			func(*strconv.NumError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"RouterDecodeError", &RouterDecodeError{}, []TranslateOption{WithXginRouterDecodeError(
			func(*RouterDecodeError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"InvalidValidationError", &validator.InvalidValidationError{}, []TranslateOption{WithValidatorInvalidTypeError(
			func(*validator.InvalidValidationError) (result map[string]string, need4xx bool) { return nil, true })}, nil, true},
		{"ValidationErrors", validator.ValidationErrors{}, []TranslateOption{WithValidatorFieldsError(
			func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
				return nil, false
			})}, nil, false},
		{"MultiFieldsError", &xvalidator.MultiFieldsError{}, []TranslateOption{WithXvalidatorMultiFieldsError(
			func(*xvalidator.MultiFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
				return nil, false
			})}, nil, false},
		{"WithTranslatableError", translatableError("TODO"), []TranslateOption{}, map[string]string{"_": "TODO"}, true},
		{"WithTranslatableError", translatableError("TODO"), []TranslateOption{WithTranslatableError(
			func(e TranslatableError) (result map[string]string, need4xx bool) {
				return map[string]string{"_x_": e.Error()}, true
			})}, map[string]string{"_x_": "TODO"}, true},
		{"NilExtraErrors", errors.New("TODO"), []TranslateOption{WithExtraErrorsTranslate(nil)}, nil, false},
		{"ExtraErrors", errors.New("TODO"), []TranslateOption{WithExtraErrorsTranslate(
			func(e error) (result map[string]string, need4xx bool) {
				return map[string]string{"_": e.Error()}, true
			})}, map[string]string{"_": "TODO"}, true},
	} {
		t.Run("other_"+tc.name, func(t *testing.T) {
			result, need4xx := TranslateBindingError(tc.giveErr, tc.giveOpts...)
			xtesting.Equal(t, result, tc.wantMap)
			xtesting.Equal(t, need4xx, tc.want4xx)
		})
	}
}
