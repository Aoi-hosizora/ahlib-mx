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
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDumpRequest(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("nil", func(c *gin.Context) { c.JSON(200, DumpRequest(nil)) })
	app.GET("nil_http", func(c *gin.Context) { c.JSON(200, DumpHttpRequest(nil)) })
	app.GET("all", func(c *gin.Context) { c.JSON(200, DumpRequest(c)) })
	app.GET("all_http", func(c *gin.Context) { c.JSON(200, DumpHttpRequest(c.Request)) })
	app.GET("reqline", func(c *gin.Context) { c.JSON(200, DumpRequest(c, WithIgnoreRequestLine(true))) })
	app.GET("reqline_http", func(c *gin.Context) { c.JSON(200, DumpRequest(c, WithIgnoreRequestLine(true))) })
	app.GET("retain1", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithRetainHeaders("X-Test")))
	})
	app.GET("retain2", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithRetainHeaders("X-TEST", "User-Agent")))
	})
	app.GET("retain3", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithRetainHeaders("X-Multi", "X-XXX"), WithIgnoreHeaders("Host")))
	})
	app.GET("ignore1", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithIgnoreHeaders("X-Test")))
	})
	app.GET("ignore2", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithIgnoreHeaders("X-TEST", "Host", "X-XXX")))
	})
	app.GET("ignore3", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithIgnoreHeaders("X-Multi")))
	})
	app.GET("secret1", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithSecretHeaders("X-Test"), WithIgnoreHeaders("X-Multi")))
	})
	app.GET("secret2", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithIgnoreHeaders("X-Multi"), WithSecretHeaders("X-TEST"), WithSecretReplace("***")))
	})
	app.GET("secret3", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithRetainHeaders("X-Multi"), WithSecretHeaders("X-Multi"), WithSecretReplace("***")))
	})

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	req := func(method, url string) []string {
		req, _ := http.NewRequest(method, url, nil)
		req.Header = http.Header{
			"Host":            []string{"localhost:12345"},
			"Accept-Encoding": []string{"gzip"},
			"User-Agent":      []string{"Go-http-client/1.1"},
			"X-Test":          []string{"xxx"},
			"X-Multi":         []string{"yyy", "zzz"},
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return []string{}
		}
		bs, _ := ioutil.ReadAll(resp.Body)
		arr := make([]string, 0)
		_ = json.Unmarshal(bs, &arr)
		return arr
	}

	for _, tc := range []struct {
		giveEp  string
		wantArr []string
	}{
		{"nil", nil},
		{"nil_http", nil},
		{"all", []string{"GET /all HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"all_http", []string{"GET /all_http HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"reqline", []string{"Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"reqline_http", []string{"Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"retain1", []string{"GET /retain1 HTTP/1.1", "X-Test: xxx"}},
		{"retain2", []string{"GET /retain2 HTTP/1.1", "User-Agent: Go-http-client/1.1", "X-Test: xxx"}},
		{"retain3", []string{"GET /retain3 HTTP/1.1", "X-Multi: yyy", "X-Multi: zzz"}},
		{"ignore1", []string{"GET /ignore1 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz"}},
		{"ignore2", []string{"GET /ignore2 HTTP/1.1", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz"}},
		{"ignore3", []string{"GET /ignore3 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Test: xxx"}},
		{"secret1", []string{"GET /secret1 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Test: *"}},
		{"secret2", []string{"GET /secret2 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Test: ***"}},
		{"secret3", []string{"GET /secret3 HTTP/1.1", "X-Multi: ***", "X-Multi: ***"}},
	} {
		t.Run(tc.giveEp, func(t *testing.T) {
			xtesting.Equal(t, req("GET", "http://127.0.0.1:12345/"+tc.giveEp), tc.wantArr)
		})
	}
}

func TestPprofWrap(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	PprofWrap(app)
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		giveMethod string
		giveUrl    string
	}{
		{"GET", "debug/pprof/"},
		{"GET", "debug/pprof/heap"},
		{"GET", "debug/pprof/goroutine"},
		{"GET", "debug/pprof/allocs"},
		{"GET", "debug/pprof/block"},
		{"GET", "debug/pprof/threadcreate"},
		{"GET", "debug/pprof/cmdline"},
		{"GET", "debug/pprof/profile"}, // <<< too slow
		{"GET", "debug/pprof/symbol"},
		{"POST", "debug/pprof/symbol"},
		{"GET", "debug/pprof/trace"},
		{"GET", "debug/pprof/mutex"},
	} {
		t.Run(tc.giveUrl, func(t *testing.T) {
			u := "http://localhost:12345/" + tc.giveUrl
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			req, _ := http.NewRequestWithContext(ctx, tc.giveMethod, u, nil) // after go113
			if tc.giveMethod == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if !errors.Is(err, context.DeadlineExceeded) {
				xtesting.Nil(t, err)
				xtesting.Equal(t, resp.StatusCode, 200)
			}
			cancel()
		})
	}
}

type mockValidator struct{}

func (f mockValidator) ValidateStruct(interface{}) error {
	return nil
}

func (f mockValidator) Engine() interface{} {
	return nil // fake
}

func TestValidatorAndTranslator(t *testing.T) {
	// validator
	val, err := GetValidatorEngine()
	xtesting.Nil(t, err)
	type testStruct1 struct {
		String string `binding:"required"`
	}
	xtesting.NotNil(t, val.Struct(&testStruct1{}))
	xtesting.Nil(t, val.Struct(&testStruct1{String: "xxx"}))
	val.SetTagName("validate") // change to use validate
	type testStruct2 struct {
		String string `validate:"required,gt=2"`
	}
	xtesting.NotNil(t, val.Struct(&testStruct2{}))
	xtesting.NotNil(t, val.Struct(&testStruct2{String: "xx"}))
	xtesting.Nil(t, val.Struct(&testStruct2{String: "xxx"}))
	val.SetTagName("binding") // default tag is binding

	// translator
	type testStruct3 struct {
		String string `binding:"required,ne=hhh" json:"str"`
	}
	xvalidator.UseTagAsFieldName(val, "json")
	xtesting.Equal(t, val.Struct(&testStruct3{}).Error(), "Key: 'testStruct3.str' Error:Field validation for 'str' failed on the 'required' tag")
	xtesting.Equal(t, val.Struct(&testStruct3{String: "hhh"}).Error(), "Key: 'testStruct3.str' Error:Field validation for 'str' failed on the 'ne' tag")
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

		reqText := val.Struct(&testStruct3{}).(validator.ValidationErrors).Translate(ut)["testStruct3.str"]
		xtesting.Equal(t, reqText, tc.wantRequiredText)
		neText := val.Struct(&testStruct3{String: "hhh"}).(validator.ValidationErrors).Translate(ut)["testStruct3.str"]
		xtesting.Equal(t, neText, tc.wantNotEqualText)

		reqText = xvalidator.TranslateValidationErrors(val.Struct(&testStruct3{}).(validator.ValidationErrors), ut, false)["str"]
		xtesting.Equal(t, reqText, tc.wantRequiredText)
		neText = xvalidator.TranslateValidationErrors(val.Struct(&testStruct3{String: "hhh"}).(validator.ValidationErrors), ut, false)["str"]
		xtesting.Equal(t, neText, tc.wantNotEqualText)
	}

	// mismatched validator engine
	motoVal := binding.Validator
	motoTrans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	binding.Validator = &mockValidator{}
	_, err = GetValidatorEngine()
	xtesting.NotNil(t, err)
	_, err = GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.NotNil(t, err)
	xtesting.NotNil(t, EnableParamRegexpBinding())
	xtesting.NotNil(t, EnableRFC3339DateBinding())
	xtesting.NotNil(t, EnableParamRegexpBindingTranslator(motoTrans))
	xtesting.NotNil(t, EnableRFC3339DateBindingTranslator(motoTrans))

	binding.Validator = motoVal
	_, err = GetValidatorEngine()
	xtesting.Nil(t, err)
	trans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)
	xtesting.Nil(t, EnableRFC3339DateTimeBinding())
	xtesting.Nil(t, EnableRFC3339DateTimeBindingTranslator(trans))
}

func TestAddBindingAndTranslator(t *testing.T) {
	val, _ := GetValidatorEngine()
	xvalidator.UseTagAsFieldName(val, "json")
	trans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, AddBinding("re_number", xvalidator.RegexpValidator(regexp.MustCompile(`^[0-9]+$`))))
	xtesting.Nil(t, AddTranslation(trans, "re_number", "{0} should be a number string", true))
	xtesting.Nil(t, AddBinding("range_name", xvalidator.LengthInRangeValidator(3, 10)))
	xtesting.Nil(t, AddTranslation(trans, "range_name", "{0} should be in range of [3, 10]", true))
	xtesting.Nil(t, EnableParamRegexpBinding())
	xtesting.Nil(t, EnableParamRegexpBindingTranslator(trans))
	xtesting.Nil(t, EnableRFC3339DateBinding())
	xtesting.Nil(t, EnableRFC3339DateBindingTranslator(trans))
	xtesting.Nil(t, EnableRFC3339DateTimeBinding())
	xtesting.Nil(t, EnableRFC3339DateTimeBindingTranslator(trans))

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
			ctx.JSON(400, &gin.H{"success": false, "details": translations})
		} else {
			ctx.JSON(200, &gin.H{"success": true})
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
		{map[string]interface{}{"number": "123", "abc": "aabbcc", "date": "2021-02-03", "datetime": "2021-02-03T02:10:13+08:00", "username": "333", "num": -1}, true, nil},
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
	motoVal := binding.Validator
	sv := xvalidator.NewCustomStructValidator()
	sv.SetMessageTagName("message")
	sv.SetValidatorTagName("binding")
	binding.Validator = sv
	defer func() { binding.Validator = motoVal }()
	val, err := GetValidatorEngine()
	xtesting.Nil(t, err)
	trans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)
	xvalidator.UseTagAsFieldName(val, "json")
	xtesting.Nil(t, EnableParamRegexpBinding())
	xtesting.Nil(t, EnableParamRegexpBindingTranslator(trans))
	xtesting.Nil(t, AddBinding("re_number", xvalidator.RegexpValidator(regexp.MustCompile(`^[0-9]+$`))))
	xtesting.Nil(t, AddTranslation(trans, "re_number", "{0} should be a number string", true))
	xtesting.Nil(t, AddBinding("range_int", xvalidator.Or(xvalidator.EqualValidator(0), xvalidator.And(xvalidator.GreaterThenOrEqualValidator(4), xvalidator.NotEqualValidator(5)))))

	type testStruct struct {
		Str string `binding:"required,re_number" json:"_str" form:"_str" message:"required|_str should be set and can not be empty"`
		Int *int32 `binding:"required,range_int" json:"_int" form:"_int" message:"range_int|_int should be in (x == 0) \\|\\| (x >= 4 && x != 5)"`
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("", func(ctx *gin.Context) {
		if err := ctx.ShouldBindQuery(&testStruct{}); err != nil {
			ctx.JSON(400, &gin.H{"success": false, "details": err.(*xvalidator.ValidateFieldsError).Translate(trans, false)})
		} else {
			ctx.JSON(200, &gin.H{"success": true})
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
			map[string]interface{}{"_str": "_str should be set and can not be empty", "_int": "_int is a required field"},
		},
		{map[string]interface{}{"_str": "", "_int": "0"}, false,
			map[string]interface{}{"_str": "_str should be set and can not be empty"},
		},
		{map[string]interface{}{"_str": "a0", "_int": "3"}, false,
			map[string]interface{}{"_str": "_str should be a number string", "_int": "_int should be in (x == 0) || (x >= 4 && x != 5)"},
		},
		{map[string]interface{}{"_str": "abc", "_int": "5"}, false,
			map[string]interface{}{"_str": "_str should be a number string", "_int": "_int should be in (x == 0) || (x >= 4 && x != 5)"},
		},
		{map[string]interface{}{"_str": "123", "_int": "4"}, true, nil},
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

func TestTranslateBindingError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	val, _ := GetValidatorEngine()
	xvalidator.UseTagAsFieldName(val, "json")
	trans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	newVal := xvalidator.NewCustomStructValidator()
	newVal.SetMessageTagName("message")
	newVal.SetValidatorTagName("binding")
	xvalidator.UseTagAsFieldName(newVal.ValidateEngine(), "json")

	app := gin.New()
	type testStruct struct {
		Str string `json:"str" form:"str" binding:"required" message:"required|str should be not null and not empty"`
		Int int32  `json:"int" form:"int" binding:"required" message:"required|int should be not null and not zero"`
	}
	respond := func(c *gin.Context, code int, details map[string]string) {
		if code == 200 {
			c.JSON(200, &gin.H{"success": true})
		} else {
			c.JSON(code, &gin.H{"success": false, "details": details})
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
			motoVal := binding.Validator
			binding.Validator = newVal
			defer func() { binding.Validator = motoVal }()
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
		{"body", `{}`, "", 500, nil},
		{"body", `{}`, "useTrans=true", 400,
			map[string]interface{}{"str": "str is a required field", "int": "int is a required field"},
		},
		{"body", `{"str": "", "int": 0}`, "", 500, nil}, // TODO
		{"body", `{"str": "", "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"str": "str is a required field", "int": "int is a required field"},
		},
		// xvalidator required
		{"body", `{}`, "useCustom=true", 500, nil},
		{"body", `{}`, "useCustom=true&useTrans=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{"str": "", "int": 0}`, "useCustom=true", 500, nil}, // TODO
		{"body", `{"str": "", "int": 0}`, "useCustom=true&useTrans=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		// ok
		{"body", `{"str": "abc", "int": 1}`, "", 200, nil},
		// strconv number error
		{"id/a", ``, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter is not a number"},
		},
		{"id/3.14", ``, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter is not a number"},
		},
		{"id/999999999999999999999999999999", `useNumError=true`, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter is out of range"},
		},
		// router decode error
		{"id/a", ``, "", 400,
			map[string]interface{}{"id": "router parameter id is not a number"},
		},
		{"id/3.14", ``, "", 400,
			map[string]interface{}{"id": "router parameter id is not a number"},
		},
		{"id/3.14", ``, "ignoreField=true", 400,
			map[string]interface{}{"router parameter": "router parameter is not a number"},
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

	// other
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
		{"ValidateFieldsError", &xvalidator.ValidateFieldsError{}, []TranslateOption{WithXvalidatorValidateFieldsError(
			func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
				return nil, false
			})}, nil, false},
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

func TestLogger(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	rand.Seed(time.Now().UnixNano())

	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)
	std := false

	app := gin.New()
	app.Use(func(c *gin.Context) {
		start := time.Now()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(20))) // fake duration
		c.Next()
		end := time.Now()

		if !std {
			LogToLogrus(l1, c, start, end)
			LogToLogrus(l1, c, start, end, WithExtraText("extra"))
			LogToLogrus(l1, c, start, end, WithExtraFields(map[string]interface{}{"k": "v"}))
			LogToLogrus(l1, c, start, end, WithExtraFieldsV("k", "v"))
		} else {
			LogToLogger(l2, c, start, end)
			LogToLogger(l2, c, start, end, WithExtraText("extra"))
			LogToLogger(l2, c, start, end, WithExtraFields(map[string]interface{}{"k": "v"}))
			LogToLogger(l2, c, start, end, WithExtraFieldsV("k", "v"))
		}
	})
	app.GET("/200", func(c *gin.Context) { c.JSON(200, &gin.H{"status": "200 success"}) })
	app.GET("/304", func(c *gin.Context) { c.Status(304) })
	app.GET("/403", func(c *gin.Context) { c.JSON(403, &gin.H{"status": "403 forbidden"}) })
	app.GET("/500", func(c *gin.Context) { c.JSON(500, &gin.H{"status": "500 internal server error"}) })
	app.POST("/XX", func(c *gin.Context) { _ = c.Error(errors.New("test error")) })

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, s := range []bool{false, true} {
		std = s
		_, _ = http.Get("http://127.0.0.1:12345/200")
		_, _ = http.Get("http://127.0.0.1:12345/403?query=string")
		_, _ = http.Get("http://127.0.0.1:12345/304")
		_, _ = http.Get("http://127.0.0.1:12345/500#anchor")
		_, _ = http.Post("http://127.0.0.1:12345/XX", "application/json", nil)
	}
}
