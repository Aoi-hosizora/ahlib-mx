package xgin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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
		{"GET", "debug/pprof/profile"},
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
}

func TestRequiredAndOmitempty(t *testing.T) {
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

// func TestValidator(t *testing.T) {
// 	app := gin.New()
// 	enTrans, err := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	log.Println(EnableRegexpBinding())
// 	log.Println(EnableRFC3339DateBinding())
// 	log.Println(EnableRFC3339DateTimeBinding())
// 	log.Println(EnableRegexpBindingTranslator(enTrans))
// 	log.Println(EnableRFC3339DateBindingTranslator(enTrans))
// 	log.Println(EnableRFC3339DateTimeBindingTranslator(enTrans))
//
// 	type st struct {
// 		A string `binding:"regexp=^[abc]+$"`
// 		B string `binding:"date"`
// 		C string `binding:"datetime"`
// 		D string `binding:"gte=2"`
// 	}
//
// 	app.GET("", func(ctx *gin.Context) {
// 		log.Println(ctx.Request.RequestURI)
// 		st := &st{}
// 		err := ctx.ShouldBindQuery(st)
// 		if err != nil {
// 			translations := err.(validator.ValidationErrors).Translate(enTrans)
// 			ctx.JSON(200, &gin.H{
// 				"msg":    err.Error(),
// 				"detail": translations,
// 			})
// 		}
// 	})
//
// 	// http://localhost:1234/?A=a&B=2020-11-16&C=2020-11-16T21:44:03Z&D=555
// 	/*
// 		"detail": {
// 			"st.A": "A must matches regexp /^[abc]+$/",
// 			"st.B": "B must be an RFC3339 Date",
// 			"st.C": "C must be an RFC3339 DateTime",
// 			"st.D": "D must be at least 2 characters in length"
// 		},
// 	*/
// 	server := &http.Server{Addr: ":12345", Handler: app}
// 	go server.ListenAndServe()
// 	defer server.Shutdown(context.Background())
// }

func TestGinRouter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	xtesting.NotPanic(t, func() {
		app := gin.New()
		app.GET(":a", func(*gin.Context) {})
		app.GET(":a/b", func(*gin.Context) {})
	})

	xtesting.PanicWithValue(t, "':b' in new path '/:b' conflicts with existing wildcard ':a' in existing prefix '/:a'", func() {
		app := gin.New()
		app.GET(":a", func(*gin.Context) {})
		app.GET(":b", func(*gin.Context) {})
	})

	xtesting.PanicWithValue(t, "'a' in new path '/a' conflicts with existing wildcard ':a' in existing prefix '/:a'", func() {
		app := gin.New()
		app.GET(":a", func(*gin.Context) {})
		app.GET("a", func(*gin.Context) {})
	})

	xtesting.PanicWithValue(t, "'a' in new path '/a' conflicts with existing wildcard ':a' in existing prefix '/:a'", func() {
		app := gin.New()
		app.GET(":a/b", func(*gin.Context) {})
		app.GET("a", func(*gin.Context) {})
	})

	xtesting.PanicWithValue(t, "'b' in new path '/:a/b' conflicts with existing wildcard ':b' in existing prefix '/:a/:b'", func() {
		app := gin.New()
		app.GET(":a/:b", func(*gin.Context) {})
		app.GET(":a/b", func(*gin.Context) {})
	})
}

func TestAppRouter(t *testing.T) {
	app := gin.New()
	app.HandleMethodNotAllowed = true
	app.NoRoute(func(c *gin.Context) {
		c.String(404, "%s %s %s 404", c.Request.Method, c.FullPath(), c.Request.URL.Path)
	})
	app.NoMethod(func(c *gin.Context) {
		c.String(405, "%s %s %s 405", c.Request.Method, c.FullPath(), c.Request.URL.Path)
	})
	fn := func(c *gin.Context) {
		c.String(200, "%s %s %s %s %s %s", c.Request.Method, c.FullPath(), c.Request.URL.Path, c.Param("x"), c.Param("y"), c.Param("z"))
	}

	xtesting.Panic(t, func() {
		ap := NewAppRouter(app, app)
		ap.GET(":a", fn)
		ap.GET(":b", fn)
	})
	xtesting.Panic(t, func() {
		ap := NewAppRouter(app, app)
		ap.GET(":a/:b", fn)
		ap.GET(":c/:d", fn)
	})

	g := app.Group("v1")
	ar := NewAppRouter(app, g)
	{
		ar.GET("", fn)

		ar.GET("a", fn)
		ar.GET(":x", fn)

		ar.GET("a/b", fn)
		ar.GET("a/:y", fn)
		ar.GET(":x/b", fn)
		ar.GET(":x/:y", fn)

		ar.GET("a/b/c", fn)
		ar.GET("a/b/:z", fn)
		ar.GET("a/:y/c", fn)
		ar.GET("a/:y/:z", fn)
		ar.GET(":x/b/c", fn)
		ar.GET(":x/b/:z", fn)
		ar.GET(":x/:y/c", fn)
		ar.GET(":x/:y/:z", fn)

		ar.GET("a/b/c/d", fn)

		ar.POST("a", fn)
		ar.POST(":x", fn)
		ar.POST("a/b", fn)
		ar.POST("a/b/c", fn)
	}
	ar.Register()

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		giveMethod   string
		giveUrl      string
		want404      bool
		want405      bool
		wantFullPath string
		wantPath     string
		wantX        string
		wantY        string
		wantZ        string
	}{
		{http.MethodGet, "v1", false, false, "/v1", "/v1", "", "", ""},

		{http.MethodGet, "v1/a", false, false, "/v1/a", "/v1/a", "", "", ""},
		{http.MethodGet, "v1/m", false, false, "/v1/:x", "/v1/m", "m", "", ""},

		{http.MethodGet, "v1/a/b", false, false, "/v1/a/b", "/v1/a/b", "", "", ""},
		{http.MethodGet, "v1/a/m", false, false, "/v1/a/:y", "/v1/a/m", "", "m", ""},
		{http.MethodGet, "v1/m/b", false, false, "/v1/:x/b", "/v1/m/b", "m", "", ""},
		{http.MethodGet, "v1/m/n", false, false, "/v1/:x/:y", "/v1/m/n", "m", "n", ""},

		{http.MethodGet, "v1/a/b/c", false, false, "/v1/a/b/c", "/v1/a/b/c", "", "", ""},
		{http.MethodGet, "v1/a/b/m", false, false, "/v1/a/b/:z", "/v1/a/b/m", "", "", "m"},
		{http.MethodGet, "v1/a/m/c", false, false, "/v1/a/:y/c", "/v1/a/m/c", "", "m", ""},
		{http.MethodGet, "v1/a/m/n", false, false, "/v1/a/:y/:z", "/v1/a/m/n", "", "m", "n"},
		{http.MethodGet, "v1/m/b/c", false, false, "/v1/:x/b/c", "/v1/m/b/c", "m", "", ""},
		{http.MethodGet, "v1/m/b/n", false, false, "/v1/:x/b/:z", "/v1/m/b/n", "m", "", "n"},
		{http.MethodGet, "v1/m/n/c", false, false, "/v1/:x/:y/c", "/v1/m/n/c", "m", "n", ""},
		{http.MethodGet, "v1/m/n/o", false, false, "/v1/:x/:y/:z", "/v1/m/n/o", "m", "n", "o"},

		{http.MethodGet, "v1/a/b/c/d", false, false, "/v1/a/b/c/d", "/v1/a/b/c/d", "", "", ""},
		{http.MethodGet, "v1/m/n/o/p", true, false, "", "/v1/m/n/o/p", "", "", ""},

		{http.MethodPost, "v1/a", false, false, "/v1/a", "/v1/a", "", "", ""},
		{http.MethodPost, "v1/m", false, false, "/v1/:x", "/v1/m", "m", "", ""},
		{http.MethodPost, "v1/a/b", false, false, "/v1/a/b", "/v1/a/b", "", "", ""},
		{http.MethodPost, "v1/m/n", false, true, "", "/v1/m/n", "", "", ""},
		{http.MethodPost, "v1/a/b/c", false, false, "/v1/a/b/c", "/v1/a/b/c", "", "", ""},
		{http.MethodPost, "v1/m/n/o", false, true, "", "/v1/m/n/o", "", "", ""},
		{http.MethodPost, "v1/a/b/c/d", false, true, "", "/v1/a/b/c/d", "", "", ""}, // <<<
		{http.MethodPost, "v1/a/b/c/e", false, true, "", "/v1/a/b/c/e", "", "", ""}, // <<<
		{http.MethodPost, "v1/m/n/o/p", false, true, "", "/v1/m/n/o/p", "", "", ""}, // <<<
		{http.MethodPost, "v1/a/b/c/d/e", true, false, "", "/v1/a/b/c/d/e", "", "", ""},
		{http.MethodPost, "v1/m/n/o/p/q", true, false, "", "/v1/m/n/o/p/q", "", "", ""},
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
		if err != nil {
			continue
		}
		bs, _ := ioutil.ReadAll(resp.Body)
		text := string(bs)

		if tc.want404 {
			xtesting.Equal(t, resp.StatusCode, 404)
			xtesting.Equal(t, text, fmt.Sprintf("%s %s %s 404", tc.giveMethod, tc.wantFullPath, tc.wantPath))
		} else if tc.want405 {
			xtesting.Equal(t, resp.StatusCode, 405)
			xtesting.Equal(t, text, fmt.Sprintf("%s %s %s 405", tc.giveMethod, tc.wantFullPath, tc.wantPath))
		} else {
			xtesting.Equal(t, resp.StatusCode, 200)
			xtesting.Equal(t, text, fmt.Sprintf("%s %s %s %s %s %s", tc.giveMethod, tc.wantFullPath, tc.wantPath, tc.wantX, tc.wantY, tc.wantZ))
		}
	}
}

func TestLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)
	std := false

	rand.Seed(time.Now().UnixNano())
	gin.SetMode(gin.ReleaseMode)
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
