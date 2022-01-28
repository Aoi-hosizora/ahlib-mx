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
	"os"
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
	app := NewEngineWithoutLogging()
	restore := HideDebugPrintRoute()
	WrapPprof(app)
	// [GIN-debug] GET    /debug/pprof/             --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func12 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/heap         --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func13 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/goroutine    --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func14 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/allocs       --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func15 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/block        --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func16 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/threadcreate --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func17 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/cmdline      --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func18 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/profile      --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func19 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/symbol       --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func20 (1 handlers)
	// [GIN-debug] POST   /debug/pprof/symbol       --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func20 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/trace        --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func21 (1 handlers)
	// [GIN-debug] GET    /debug/pprof/mutex        --> github.com/Aoi-hosizora/ahlib-web/xgin.glob..func22 (1 handlers)
	restore()
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

func TestGetProxyEnv(t *testing.T) {
	for _, key := range []string{
		"http_proxy", "https_proxy", "socks_proxy",
	} {
		e, ok := os.LookupEnv(key)
		if ok {
			//goland:noinspection GoDeferInLoop
			defer os.Setenv(key, e)
		}
	}

	os.Setenv("http_proxy", "")
	os.Setenv("https_proxy", "")
	os.Setenv("socks_proxy", "")
	hp, hsp, ssp := GetProxyEnv()
	xtesting.Equal(t, hp, "")
	xtesting.Equal(t, hsp, "")
	xtesting.Equal(t, ssp, "")

	os.Setenv("http_proxy", "http://localhost:9000")
	os.Setenv("https_proxy", "https://localhost:9000")
	os.Setenv("socks_proxy", "socks://localhost:9000")
	hp, hsp, ssp = GetProxyEnv()
	xtesting.Equal(t, hp, "http://localhost:9000")
	xtesting.Equal(t, hsp, "https://localhost:9000")
	xtesting.Equal(t, ssp, "socks://localhost:9000")

	os.Setenv("http_proxy", "")
	os.Setenv("https_proxy", "")
	os.Setenv("socks_proxy", "")
	hp, hsp, ssp = GetProxyEnv()
	xtesting.Equal(t, hp, "")
	xtesting.Equal(t, hsp, "")
	xtesting.Equal(t, ssp, "")
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
				msg := fmt.Sprintf("[Gin] %8d - %12s - %15s - %10s - %-7s %s", p.Status, p.Latency.String(), p.ClientIP, xnumber.FormatByteSize(float64(p.Length)), p.Method, p.Path)
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
