package xgin

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

// PprofWrap adds several routes from package `net/http/pprof` to *gin.Engine object.
// Reference see https://github.com/DeanThompson/ginpprof.
func PprofWrap(router *gin.Engine) {
	for _, r := range []struct {
		Method  string
		Path    string
		Handler gin.HandlerFunc
	}{
		{"GET", "/debug/pprof/", indexHandler()},
		{"GET", "/debug/pprof/heap", heapHandler()},
		{"GET", "/debug/pprof/goroutine", goroutineHandler()},
		{"GET", "/debug/pprof/allocs", allocsHandler()},
		{"GET", "/debug/pprof/block", blockHandler()},
		{"GET", "/debug/pprof/threadcreate", threadCreateHandler()},
		{"GET", "/debug/pprof/cmdline", cmdlineHandler()},
		{"GET", "/debug/pprof/profile", profileHandler()},
		{"GET", "/debug/pprof/symbol", symbolHandler()},
		{"POST", "/debug/pprof/symbol", symbolHandler()},
		{"GET", "/debug/pprof/trace", traceHandler()},
		{"GET", "/debug/pprof/mutex", mutexHandler()},
	} {
		router.Handle(r.Method, r.Path, r.Handler)
	}
}

// /debug/pprof
func indexHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Index(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/heap
func heapHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("heap").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/goroutine
func goroutineHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("goroutine").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/allocs
func allocsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("allocs").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/block
func blockHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("block").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/threadcreate
func threadCreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/cmdline
func cmdlineHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Cmdline(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/profile
func profileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Profile(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/symbol
func symbolHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Symbol(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/trace
func traceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Trace(ctx.Writer, ctx.Request)
	}
}

// /debug/pprof/mutex
func mutexHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler("mutex").ServeHTTP(ctx.Writer, ctx.Request)
	}
}
