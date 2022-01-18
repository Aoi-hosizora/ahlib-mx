package xgin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
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
		c.JSON(200, DumpRequest(c, WithIgnoreHeaders("X-Multi"), WithSecretHeaders("X-TEST"), WithSecretPlaceholder("***")))
	})
	app.GET("secret3", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithRetainHeaders("X-Multi"), WithSecretHeaders("X-Multi"), WithSecretPlaceholder("***")))
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

func TestWrapPprof(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	app := NewEngineWithoutDebugWarning()
	rfn := HideDebugPrintRoute()
	WrapPprof(app)
	rfn()
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

func (f mockValidator) ValidateStruct(interface{}) error { return nil }
func (f mockValidator) Engine() interface{}              { return nil /* fake */ }

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
	xtesting.Nil(t, GetGlobalTranslator())
	ut, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)
	SetGlobalTranslator(ut)
	xtesting.NotNil(t, GetGlobalTranslator())
	SetGlobalTranslator(nil)
	xtesting.Nil(t, GetGlobalTranslator())
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
			ctx.JSON(400, gin.H{"success": false, "details": err.(*xvalidator.ValidateFieldsError).Translate(trans, false)})
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

func TestResponseLogger(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	rand.Seed(time.Now().UnixNano())

	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.Default()
	std := false
	custom := false

	app := gin.New()
	loggerMiddleware := func(c *gin.Context) {
		start := time.Now()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(20))) // fake duration
		c.Next()
		end := time.Now()

		// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
		if custom {
			FormatResponseFunc = func(p *ResponseLoggerParam) string {
				msg := fmt.Sprintf("[Gin] %8d - %12s - %15s - %10s - %-7s %s", p.Status, p.Latency.String(), p.ClientIP, xnumber.RenderByte(float64(p.Length)), p.Method, p.Path)
				if p.ErrorMsg != "" {
					msg += fmt.Sprintf(" - err: %s", p.ErrorMsg)
				}
				return msg
			}
			FieldifyResponseFunc = func(p *ResponseLoggerParam) logrus.Fields {
				return logrus.Fields{"module": "gin", "method": p.Method, "path": p.Path, "status": p.Status}
			}
		}
		if !std {
			LogResponseToLogrus(l1, c, start, end)
			LogResponseToLogrus(l1, c, start, end, WithExtraText(" | extra"))
			LogResponseToLogrus(l1, c, start, end, WithExtraFields(map[string]interface{}{"k": "v"}))
			LogResponseToLogrus(l1, c, start, end, WithExtraFieldsV("k", "v"))
		} else {
			LogResponseToLogger(l2, c, start, end)
			LogResponseToLogger(l2, c, start, end, WithExtraText(" | extra"))
			LogResponseToLogger(l2, c, start, end, WithExtraFields(map[string]interface{}{"k": "v"}))
			LogResponseToLogger(l2, c, start, end, WithExtraFieldsV("k", "v"))
		}
		if custom {
			FormatResponseFunc = nil
			FieldifyResponseFunc = nil
		}
	}
	app.Use(loggerMiddleware)
	app.GET("/200", func(c *gin.Context) { c.JSON(200, gin.H{"status": "200 success"}) })
	app.GET("/304", func(c *gin.Context) { c.Status(304) })
	app.GET("/403", func(c *gin.Context) { c.JSON(403, gin.H{"status": "403 forbidden"}) })
	app.GET("/500", func(c *gin.Context) { c.JSON(500, gin.H{"status": "500 internal server error"}) })
	app.POST("/XX", func(c *gin.Context) { _ = c.Error(errors.New("test error")) })

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, s := range []bool{false, true} {
		std = s
		for _, c := range []bool{false, true} {
			custom = c
			log.Printf("std: %t, custom: %t", std, custom)
			_, _ = http.Get("http://127.0.0.1:12345/200")
			_, _ = http.Get("http://127.0.0.1:12345/403?query=string")
			_, _ = http.Get("http://127.0.0.1:12345/304")
			_, _ = http.Get("http://127.0.0.1:12345/500#anchor")
			_, _ = http.Post("http://127.0.0.1:12345/XX", "application/json", nil)
		}
	}

	xtesting.NotPanic(t, func() {
		LogResponseToLogrus(l1, nil, time.Now(), time.Now())
		LogResponseToLogger(l2, nil, time.Now(), time.Now())
		LogResponseToLogrus(nil, nil, time.Now(), time.Now())
		LogResponseToLogger(nil, nil, time.Now(), time.Now())
	})
}

func TestRecoveryLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.Default()

	for _, std := range []bool{false, true} {
		for _, custom := range []bool{false, true} {
			for _, tc := range []struct {
				giveErr     interface{}
				giveStack   xruntime.TraceStack
				giveOptions []LoggerOption
			}{
				{nil, nil, nil},
				{"test string", nil, nil},
				{errors.New("test error"), nil, nil},
				{nil, xruntime.RuntimeTraceStack(0), nil},
				{errors.New("test error"), xruntime.RuntimeTraceStack(0), nil},

				{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraText(" | extra")}},
				{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
				{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraFieldsV("k", "v")}},
				{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraText(" | extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
				{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraText(" | extra"), WithExtraFieldsV("k", "v")}},
			} {
				// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
				if custom {
					FormatRecoveryFunc = func(p *RecoveryLoggerParam) string {
						return fmt.Sprintf("[Recovery] %s, %s:%d", p.PanicMsg, p.Filename, p.LineIndex)
					}
					FieldifyRecoveryFunc = func(p *RecoveryLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "recovery", "panic_msg": p.PanicMsg, "#trace_stack": len(p.Stack)}
					}
				}
				if !std {
					LogRecoveryToLogrus(l1, tc.giveErr, tc.giveStack, tc.giveOptions...)
				} else {
					LogRecoveryToLogger(l2, tc.giveErr, tc.giveStack, tc.giveOptions...)
				}
				if custom {
					FormatRecoveryFunc = nil
					FieldifyRecoveryFunc = nil
				}
			}
		}
	}

	xtesting.NotPanic(t, func() {
		LogRecoveryToLogrus(l1, nil, nil)
		LogRecoveryToLogger(l2, nil, nil)
		LogRecoveryToLogrus(nil, nil, nil)
		LogRecoveryToLogger(nil, nil, nil)
	})
}
