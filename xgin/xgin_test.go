package xgin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestDumpRequest(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	req := func(ep string) []string {
		req, _ := http.NewRequest("GET", "http://127.0.0.1:12345/"+ep, nil)
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
		giveEp string
		giveFn func(c *gin.Context) []string
		want   []string
	}{
		{"nil", func(c *gin.Context) []string { return DumpRequest(nil) }, nil},
		{"nil_http", func(c *gin.Context) []string { return DumpHttpRequest(nil) }, nil},
		{"all", func(c *gin.Context) []string {
			return DumpRequest(c)
		}, []string{"GET /all HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"all_http", func(c *gin.Context) []string {
			return DumpHttpRequest(c.Request)
		}, []string{"GET /all_http HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"reqline", func(c *gin.Context) []string {
			return DumpRequest(c, WithIgnoreRequestLine(true))
		}, []string{"Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"reqline_http", func(c *gin.Context) []string {
			return DumpRequest(c, WithIgnoreRequestLine(true))
		}, []string{"Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz", "X-Test: xxx"}},
		{"retain1", func(c *gin.Context) []string {
			return DumpRequest(c, WithRetainHeaders("X-Test"))
		}, []string{"GET /retain1 HTTP/1.1", "X-Test: xxx"}},
		{"retain2", func(c *gin.Context) []string {
			return DumpRequest(c, WithRetainHeaders("X-TEST", "User-Agent"))
		}, []string{"GET /retain2 HTTP/1.1", "User-Agent: Go-http-client/1.1", "X-Test: xxx"}},
		{"retain3", func(c *gin.Context) []string {
			return DumpRequest(c, WithRetainHeaders("X-Multi", "X-XXX"), WithIgnoreHeaders("Host"))
		}, []string{"GET /retain3 HTTP/1.1", "X-Multi: yyy", "X-Multi: zzz"}},
		{"ignore1", func(c *gin.Context) []string {
			return DumpRequest(c, WithIgnoreHeaders("X-Test"))
		}, []string{"GET /ignore1 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz"}},
		{"ignore2", func(c *gin.Context) []string {
			return DumpRequest(c, WithIgnoreHeaders("X-TEST", "Host", "X-XXX"))
		}, []string{"GET /ignore2 HTTP/1.1", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Multi: yyy", "X-Multi: zzz"}},
		{"ignore3", func(c *gin.Context) []string {
			return DumpRequest(c, WithIgnoreHeaders("X-Multi"))
		}, []string{"GET /ignore3 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Test: xxx"}},
		{"secret1", func(c *gin.Context) []string {
			return DumpRequest(c, WithSecretHeaders("X-Test"), WithIgnoreHeaders("X-Multi"))
		}, []string{"GET /secret1 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Test: *"}},
		{"secret2", func(c *gin.Context) []string {
			return DumpRequest(c, WithIgnoreHeaders("X-Multi"), WithSecretHeaders("X-TEST"), WithSecretPlaceholder("***"))
		}, []string{"GET /secret2 HTTP/1.1", "Host: 127.0.0.1:12345", "Accept-Encoding: gzip", "User-Agent: Go-http-client/1.1", "X-Test: ***"}},
		{"secret3", func(c *gin.Context) []string {
			return DumpRequest(c, WithRetainHeaders("X-Multi"), WithSecretHeaders("X-Multi"), WithSecretPlaceholder("***"))
		}, []string{"GET /secret3 HTTP/1.1", "X-Multi: ***", "X-Multi: ***"}},
	} {
		t.Run(tc.giveEp, func(t *testing.T) {
			app.GET(tc.giveEp, func(c *gin.Context) {
				c.JSON(200, tc.giveFn(c))
			})
			xtesting.Equal(t, req(tc.giveEp), tc.want)
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	// app.Use(gin.Logger())
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		name            string
		giveCode        int
		giveAndWantData interface{}
	}{
		{"301", http.StatusMovedPermanently, "code_301"},
		{"302", http.StatusFound, "code_302"},
		{"303", http.StatusSeeOther, "code_303"},
		// {"304", http.StatusNotModified, "code_304"}, => 304 is not used for redirect
		{"307", http.StatusTemporaryRedirect, "code_307"},
		{"308", http.StatusPermanentRedirect, "code_308"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := "/" + tc.name
			index := root + "/index"
			app.GET(root, RedirectHandler(tc.giveCode, index))
			app.GET(index, func(c *gin.Context) {
				c.JSON(200, gin.H{"data": tc.giveAndWantData})
			})

			req, err := http.NewRequest("GET", "http://localhost:12345"+root, nil) // root -> index
			xtesting.Nil(t, err)
			client := http.Client{}
			resp, err := client.Do(req)
			xtesting.Nil(t, err)
			xtesting.Equal(t, resp.StatusCode, 200)
			bs, err := ioutil.ReadAll(resp.Body)
			xtesting.Nil(t, err)
			resp.Body.Close()
			data := map[string]interface{}{}
			xtesting.Nil(t, json.Unmarshal(bs, &data))
			xtesting.Equal(t, data["data"], tc.giveAndWantData)
		})
	}
}

func TestWrapPprof(t *testing.T) {
	gin.SetMode(gin.DebugMode)

	// 1.
	log.Println("============ 1")
	app := gin.New() // <<< with warning
	WrapPprof(app)   // <<< with [Gin]
	// [GIN] GET    /debug/pprof/             --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func12 (1 handlers)
	// [GIN] GET    /debug/pprof/heap         --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func13 (1 handlers)
	// [GIN] GET    /debug/pprof/goroutine    --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func14 (1 handlers)
	// [GIN] GET    /debug/pprof/allocs       --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func15 (1 handlers)
	// [GIN] GET    /debug/pprof/block        --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func16 (1 handlers)
	// [GIN] GET    /debug/pprof/threadcreate --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func17 (1 handlers)
	// [GIN] GET    /debug/pprof/cmdline      --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func18 (1 handlers)
	// [GIN] GET    /debug/pprof/profile      --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func19 (1 handlers)
	// [GIN] GET    /debug/pprof/symbol       --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func20 (1 handlers)
	// [GIN] POST   /debug/pprof/symbol       --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func20 (1 handlers)
	// [GIN] GET    /debug/pprof/trace        --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func21 (1 handlers)
	// [GIN] GET    /debug/pprof/mutex        --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func22 (1 handlers)

	// 2.
	log.Println("============ 2")
	app = gin.New() // <<< with warning
	gin.DebugPrintRouteFunc = DefaultColorizedPrintRouteFunc
	WrapPprof(app) // <<< with colorized [Gin]

	// 3.
	log.Println("============ 3")
	app = NewWithoutLogging()    // <<< no warning
	WrapPprofWithoutLogging(app) // <<< no [Gin]

	// 4.
	log.Println("============ 4")
	restore := HideDebugLogging()
	app = gin.New() // <<< no warning
	restore()
	restore = HideDebugPrintRoute()
	WrapPprof(app) // <<< no [Gin]
	restore()

	// ...
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			req, _ := http.NewRequestWithContext(ctx, tc.giveMethod, "http://localhost:12345/"+tc.giveUrl, nil) // after go113
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

func TestRouterDecodeError(t *testing.T) {
	_, err := strconv.Atoi("1a")
	rerr := NewRouterDecodeError("", "1a", err, "")
	xtesting.Equal(t, rerr.Field, "")
	xtesting.Equal(t, rerr.Input, "1a")
	xtesting.Equal(t, rerr.Err, err)
	xtesting.Equal(t, rerr.Message, "")
	xtesting.Equal(t, rerr.Error(), "parsing \"1a\": strconv.Atoi: parsing \"1a\": invalid syntax")
	xtesting.Equal(t, rerr.Unwrap(), err)
	xtesting.True(t, errors.Is(rerr, err))

	err = errors.New("non-positive number")
	rerr = NewRouterDecodeError("id", "0", err, "must be a positive number")
	xtesting.Equal(t, rerr.Field, "id")
	xtesting.Equal(t, rerr.Input, "0")
	xtesting.Equal(t, rerr.Err, err)
	xtesting.Equal(t, rerr.Message, "must be a positive number")
	xtesting.Equal(t, rerr.Error(), "parsing id \"0\": non-positive number")
	xtesting.Equal(t, rerr.Unwrap(), err)
	xtesting.True(t, errors.Is(rerr, err))

	xtesting.Panic(t, func() { _ = NewRouterDecodeError("", "", nil, "") })
}

// ATTENTION: loggerOptions related code and unit tests in xgin package and xtelebot package should keep the same as each other.
func TestLoggerOptions(t *testing.T) {
	for _, tc := range []struct {
		give       []LoggerOption
		wantMsg    string
		wantFields logrus.Fields
	}{
		{[]LoggerOption{}, "", logrus.Fields{}},
		{[]LoggerOption{nil}, "", logrus.Fields{}},
		{[]LoggerOption{nil, nil, nil}, "", logrus.Fields{}},

		{[]LoggerOption{WithExtraText("")}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("   ")}, "   ", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("  x x  ")}, "  x x  ", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("test")}, "test", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("test1"), WithExtraText(" | test2")}, " | test2", logrus.Fields{}},

		{[]LoggerOption{WithExtraFields(map[string]interface{}{})}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"true": 2, "3": 4.4})}, "", logrus.Fields{"true": 2, "3": 4.4}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"true": 2, "3": 4.4}), WithExtraFields(map[string]interface{}{"k": "v"})}, "", logrus.Fields{"k": "v"}},

		{[]LoggerOption{WithExtraFieldsV()}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil)}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil, "a", nil)}, "", logrus.Fields{"a": nil}},
		{[]LoggerOption{WithExtraFieldsV(nil, nil, "a")}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil, nil, 1, nil)}, "", logrus.Fields{"1": nil}},
		{[]LoggerOption{WithExtraFieldsV(true, 2, 3.3, 4, 5)}, "", logrus.Fields{"true": 2, "3.3": 4}},
		{[]LoggerOption{WithExtraFieldsV(true, 2, 3.3, 4, 5), WithExtraFieldsV("k", "v")}, "", logrus.Fields{"k": "v"}},

		{[]LoggerOption{WithExtraText("test"), WithExtraFields(map[string]interface{}{"1": 2})}, "test", logrus.Fields{"1": 2}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"1": 2}), WithExtraText(" | test")}, " | test", logrus.Fields{"1": 2}},
		{[]LoggerOption{WithExtraText("test"), WithExtraFieldsV(3, 4)}, "test", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFieldsV(3, 4), WithExtraText(" | test")}, " | test", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"1": 2}), WithExtraFieldsV(3, 4)}, "", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFieldsV(3, 4), WithExtraFields(map[string]interface{}{"1": 2})}, "", logrus.Fields{"1": 2}},
	} {
		ops := buildLoggerOptions(tc.give)
		msg := ""
		fields := logrus.Fields{}
		ops.ApplyToMessage(&msg)
		ops.ApplyToFields(fields)
		xtesting.Equal(t, msg, tc.wantMsg)
		xtesting.Equal(t, fields, tc.wantFields)
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
				path := p.Path
				if p.Query != "" {
					path += "?" + p.Query
				}
				msg := fmt.Sprintf("[Gin] %8d - %12s - %15s - %10s - %-7s %s", p.Status, p.Latency.String(), p.ClientIP, xnumber.FormatByteSize(float64(p.Length)), p.Method, path)
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
						return fmt.Sprintf("[Recovery] %s, %s:%d, %s", p.PanicMsg, p.FullFilename, p.LineIndex, p.FullFuncname)
					}
					FieldifyRecoveryFunc = func(p *RecoveryLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "recovery", "panic_msg": p.PanicMsg, "#trace_stack": len(p.TraceStack)}
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
