package xgin

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGinRouter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	// normal
	xtesting.NotPanic(t, func() {
		app := gin.New()
		app.GET(":a", func(*gin.Context) {})
		app.GET(":a/b", func(*gin.Context) {})
	})

	// conflict 0
	xtesting.PanicWithValue(t, "':b' in new path '/:b' conflicts with existing wildcard ':a' in existing prefix '/:a'", func() {
		app := gin.New()
		app.GET(":a", func(*gin.Context) {})
		app.GET(":b", func(*gin.Context) {})
	})

	// conflict 1: ":a" & "a"
	xtesting.PanicWithValue(t, "'a' in new path '/a' conflicts with existing wildcard ':a' in existing prefix '/:a'", func() {
		app := gin.New()
		app.GET(":a", func(*gin.Context) {})
		app.GET("a", func(*gin.Context) {})
	})

	// conflict 2: ":a/b" & "a"
	xtesting.PanicWithValue(t, "'a' in new path '/a' conflicts with existing wildcard ':a' in existing prefix '/:a'", func() {
		app := gin.New()
		app.GET(":a/b", func(*gin.Context) {})
		app.GET("a", func(*gin.Context) {})
	})

	// conflict 3: ":a/:b" & ":a/b"
	xtesting.PanicWithValue(t, "'b' in new path '/:a/b' conflicts with existing wildcard ':b' in existing prefix '/:a/:b'", func() {
		app := gin.New()
		app.GET(":a/:b", func(*gin.Context) {})
		app.GET(":a/b", func(*gin.Context) {})
	})
}

func TestAppRouterBasic(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	app := gin.New()
	app.HandleMethodNotAllowed = true

	// conflict panic
	xtesting.Panic(t, func() {
		ar := NewAppRouter(app, app)
		ar.GET(":a", func(*gin.Context) {})
		ar.GET(":b", func(*gin.Context) {})
	})
	xtesting.Panic(t, func() {
		ar := NewAppRouter(app, app)
		ar.GET(":a/:b", func(*gin.Context) {})
		ar.GET(":c/:d", func(*gin.Context) {})
	})

	// empty handler panic
	for _, tc := range []struct {
		giveFn func(string, ...gin.HandlerFunc)
	}{
		{NewAppRouter(app, app).GET},
		{NewAppRouter(app, app).POST},
		{NewAppRouter(app, app).DELETE},
		{NewAppRouter(app, app).PATCH},
		{NewAppRouter(app, app).PUT},
		{NewAppRouter(app, app).OPTIONS},
		{NewAppRouter(app, app).HEAD},
		{NewAppRouter(app, app).Any},
	} {
		xtesting.Panic(t, func() {
			tc.giveFn("")
		})
	}

	// 404 405 & logger
	PrintAppRouterRegisterFunc = func(index, count int, method, relativePath, handlerFuncname string, handlersCount int, layerFakePath string) {
		fmt.Printf("[XGIN]      %-6s ~/%-23s --> %s (%d handlers) ==> ~/%s\n", method, relativePath, handlerFuncname, handlersCount, layerFakePath)
	}
	defer func() {
		PrintAppRouterRegisterFunc = nil
	}()
	ar := NewAppRouter(app, app)
	ar.GET("a", func(c *gin.Context) {})
	ar.POST("b", func(c *gin.Context) {})
	ar.Register()
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	for _, tc := range []struct {
		giveMethod string
		giveUrl    string
		wantResp   string
	}{
		{"GET", "c", "404 page not found"},
		{"POST", "a", "405 method not allowed"},
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
			bs, _ := ioutil.ReadAll(resp.Body)
			xtesting.Equal(t, string(bs), tc.wantResp)
		}
	}
}

func TestAppRouterFunction(t *testing.T) {
	gin.SetMode(gin.DebugMode) // use debug mode
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

	// test app router
	g := app.Group("v1")
	ar := NewAppRouter(app, g)
	{
		// 0
		ar.GET("", fn)

		// 1
		ar.GET("a", fn)
		ar.GET(":x", fn)

		// 2
		ar.GET("a/b", fn)
		ar.GET("a/:y", fn)
		ar.GET(":x/b", fn)
		ar.GET(":x/:y", fn)

		// 3
		ar.GET("a/b/c", fn)
		ar.GET("a/b/:z", fn)
		ar.GET("a/:y/c", fn)
		ar.GET("a/:y/:z", fn)
		ar.GET(":x/b/c", fn)
		ar.GET(":x/b/:z", fn)
		ar.GET(":x/:y/c", fn)
		ar.GET(":x/:y/:z", fn)

		// 4
		ar.GET("a/b/c/d", fn)

		// POST
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
