package xgin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestDumpRequest(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("nil", func(c *gin.Context) {
		c.JSON(200, DumpRequest(nil))
	})
	app.GET("all", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c))
	})
	app.GET("retain", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithRetainHeaders("Host")))
	})
	app.GET("ignore", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithIgnoreHeaders("Host")))
	})
	app.GET("secret1", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithSecretHeaders("Host")))
	})
	app.GET("secret2", func(c *gin.Context) {
		c.JSON(200, DumpRequest(c, WithSecretHeaders("Host"), WithSecretReplace("***")))
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
		{"nil", []string{}},
		{"all", []string{"GET /all HTTP/1.1", "Host: localhost:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1"}},
		{"retain", []string{"GET /retain HTTP/1.1", "Host: localhost:12345"}},
		{"ignore", []string{"GET /ignore HTTP/1.1", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1"}},
		{"secret1", []string{"GET /secret1 HTTP/1.1", "Host: *", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1"}},
		{"secret2", []string{"GET /secret2 HTTP/1.1", "Host: ***", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1"}},
	} {
		xtesting.Equal(t, req("GET", "http://localhost:12345/"+tc.giveEp), tc.wantArr)
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
		// {"GET", "debug/pprof/profile"}, // <<< too slow
		{"GET", "debug/pprof/symbol"},
		{"POST", "debug/pprof/symbol"},
		{"GET", "debug/pprof/trace"},
		{"GET", "debug/pprof/mutex"},
	} {
		var resp *http.Response
		var err error
		url := "http://localhost:12345/" + tc.giveUrl
		if tc.giveMethod == http.MethodGet {
			resp, err = http.Get(url)
		} else if tc.giveMethod == http.MethodPost {
			resp, err = http.Post(url, "application/json", nil)
		}
		xtesting.Nil(t, err)
		if err == nil {
			xtesting.Equal(t, resp.StatusCode, 200)
		}
	}

	// for slow debug/pprof/profile
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:12345/debug/pprof/profile", nil) // after go113
	client := &http.Client{}
	_, _ = client.Do(req) // ignore result
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
	val.SetTagName("validate")      // change to use validate
	defer val.SetTagName("binding") // default tag is binding
	type testStruct2 struct {
		String string `validate:"required"`
	}
	xtesting.NotNil(t, val.Struct(&testStruct2{}))

	// translator
	for _, tc := range []struct {
		giveTranslator   locales.Translator
		giveRegisterFn   xvalidator.TranslationRegisterHandler
		wantRequiredText string
	}{
		{nil, nil, "Key: 'testStruct2.String' Error:Field validation for 'String' failed on the 'required' tag"},
		{xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc(), "String is a required field"},
		{xvalidator.FrLocaleTranslator(), xvalidator.FrTranslationRegisterFunc(), "String est un champ obligatoire"},
		{xvalidator.JaLocaleTranslator(), xvalidator.JaTranslationRegisterFunc(), "Stringは必須フィールドです"},
		{xvalidator.RuLocaleTranslator(), xvalidator.RuTranslationRegisterFunc(), "String обязательное поле"},
		{xvalidator.ZhLocaleTranslator(), xvalidator.ZhTranslationRegisterFunc(), "String为必填字段"},
		{xvalidator.ZhHantLocaleTranslator(), xvalidator.ZhTwTranslationRegisterFunc(), "String為必填欄位"},
	} {
		text := ""
		if tc.giveTranslator == nil || tc.giveRegisterFn == nil {
			text = val.Struct(&testStruct2{}).Error()
		} else {
			trans, err := GetValidatorTranslator(tc.giveTranslator, tc.giveRegisterFn)
			xtesting.Nil(t, err)
			var ok bool
			text, ok = val.Struct(&testStruct2{}).(validator.ValidationErrors).Translate(trans)["testStruct2.String"]
			xtesting.True(t, ok)
		}
		xtesting.Equal(t, text, tc.wantRequiredText)
	}

	// error
	trans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)
	motoVal := binding.Validator
	binding.Validator = &mockValidator{}
	defer func() { binding.Validator = motoVal }()

	_, err = GetValidatorEngine()
	xtesting.NotNil(t, err)
	_, err = GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.NotNil(t, err)
	xtesting.NotNil(t, EnableParamRegexpBinding())
	xtesting.NotNil(t, EnableParamRegexpBindingTranslator(trans))
}

func TestAddBindingAndAddTranslator(t *testing.T) {
	trans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	xtesting.Nil(t, err)

	xtesting.Nil(t, AddBinding("re_number", xvalidator.RegexpValidator(regexp.MustCompile(`[0-9]+`))))
	xtesting.Nil(t, AddTranslator(trans, "re_number", "{0} must be a number string", true))
	xtesting.Nil(t, EnableParamRegexpBinding())
	xtesting.Nil(t, EnableParamRegexpBindingTranslator(trans))
	xtesting.Nil(t, EnableRFC3339DateBinding())
	xtesting.Nil(t, EnableRFC3339DateBindingTranslator(trans))
	xtesting.Nil(t, EnableRFC3339DateTimeBinding())
	xtesting.Nil(t, EnableRFC3339DateTimeBindingTranslator(trans))

	type testStruct struct {
		Number   string `json:"number"   form:"number"   binding:"re_number"`
		Abc      string `json:"abc"      form:"abc"      binding:"regexp=^[abc]+$"`
		Date     string `json:"date"     form:"date"     binding:"date"`
		Datetime string `json:"datetime" form:"datetime" binding:"datetime"`
		Gte2     int    `json:"gte2"     form:"gte2"     binding:"gte=2"`
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("", func(ctx *gin.Context) {
		test := &testStruct{}
		err := ctx.ShouldBindQuery(test)
		if err != nil {
			translations := err.(validator.ValidationErrors).Translate(trans)
			ctx.JSON(400, &gin.H{
				"msg":    "failed",
				"detail": translations,
			})
		} else {
			ctx.JSON(200, &gin.H{
				"msg": "success",
			})
		}
	})
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		giveQuery string
		wantMap   map[string]interface{}
	}{
		{"", map[string]interface{}{
			"msg": "failed",
			"detail": map[string]interface{}{
				"testStruct.Number":   "Number must be a number string",
				"testStruct.Abc":      "Abc must matches regexp /^[abc]+$/",
				"testStruct.Date":     "Date must be an RFC3339 date",
				"testStruct.Datetime": "Datetime must be an RFC3339 datetime",
				"testStruct.Gte2":     "Gte2 must be 2 or greater",
			},
		}},
		{"?number=aaa&abc=def&date=2021-02-03&datetime=2021-02-03T02%3A10%3A13%2B08%3A00&gte2=1", map[string]interface{}{ // 2021-02-03T02:10:13+08:00
			"msg": "failed",
			"detail": map[string]interface{}{
				"testStruct.Number": "Number must be a number string",
				"testStruct.Abc":    "Abc must matches regexp /^[abc]+$/",
				"testStruct.Gte2":   "Gte2 must be 2 or greater",
			},
		}},
		{"?number=0&abc=aabbcc&date=2021/02/03&datetime=2021-02-03T02%3A10%3A13&gte2=3", map[string]interface{}{
			"msg": "failed",
			"detail": map[string]interface{}{
				"testStruct.Date":     "Date must be an RFC3339 date",
				"testStruct.Datetime": "Datetime must be an RFC3339 datetime",
			},
		}},
		{"?number=0&abc=aabbcc&date=2021-02-03&datetime=2021-02-03T02%3A10%3A13%2B08%3A00&gte2=2", map[string]interface{}{
			"msg": "success",
		}},
	} {
		resp, err := http.Get("http://localhost:12345" + tc.giveQuery)
		xtesting.Nil(t, err)
		bs, _ := ioutil.ReadAll(resp.Body)

		resultMap := make(map[string]interface{})
		err = json.Unmarshal(bs, &resultMap)
		xtesting.Nil(t, err)
		xtesting.Equal(t, resultMap, tc.wantMap)
	}
}

func TestCustomStructValidatorAndTranslate(t *testing.T) {
	// ...
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
