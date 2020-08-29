package xfiber

import (
	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http/pprof"
	"strings"
)

var (
	ppIndex        = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index)
	ppCmdline      = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline)
	ppProfile      = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile)
	ppSymbol       = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol)
	ppTrace        = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace)
	ppAllocs       = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("allocs").ServeHTTP)
	ppBlock        = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("block").ServeHTTP)
	ppGoroutine    = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("goroutine").ServeHTTP)
	ppHeap         = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("heap").ServeHTTP)
	ppMutex        = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("mutex").ServeHTTP)
	ppThreadcreate = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler("threadcreate").ServeHTTP)
)

// New pprofhandler.
// Reference to https://github.com/gofiber/pprof/blob/master/main.go.
func PprofHandler() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		if !strings.HasPrefix(c.Path(), "/debug/pprof") {
			c.Next()
			return
		}

		switch string(c.Fasthttp.URI().Path()) {
		case "/debug/pprof/":
			// Set content-type to HTML to display index page
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
			ppIndex(c.Fasthttp)
		case "/debug/pprof/cmdline":
			ppCmdline(c.Fasthttp)
		case "/debug/pprof/profile":
			ppProfile(c.Fasthttp)
		case "/debug/pprof/symbol":
			ppSymbol(c.Fasthttp)
		case "/debug/pprof/trace":
			ppTrace(c.Fasthttp)
		case "/debug/pprof/allocs":
			ppAllocs(c.Fasthttp)
		case "/debug/pprof/block":
			ppBlock(c.Fasthttp)
		case "/debug/pprof/goroutine":
			ppGoroutine(c.Fasthttp)
		case "/debug/pprof/heap":
			ppHeap(c.Fasthttp)
		case "/debug/pprof/mutex":
			ppMutex(c.Fasthttp)
		case "/debug/pprof/threadcreate":
			ppThreadcreate(c.Fasthttp)
		default: // pprof index only works with trailing slash
			c.Redirect("/debug/pprof/", 302)
		}
	}
}
