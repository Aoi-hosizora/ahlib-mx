package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// AppRoute stores a group of methods (group of routes).
type AppRoute struct {
	engine *gin.Engine
	router gin.IRouter
	groups [][]*route
}

// NewAppRoute create an instance of AppRoute.
func NewAppRoute(engine *gin.Engine, router gin.IRouter) *AppRoute {
	return &AppRoute{engine: engine, router: router, groups: [][]*route{}}
}

// route represents a route in AppRoute, include handlers and relativePath.
type route struct {
	method       string
	relativePath string
	parameters   []string // used later
	handlers     []gin.HandlerFunc
}

// newRoute create an instance of route, panic if relativePath is empty.
func newRoute(method string, relativePath string, handlers ...gin.HandlerFunc) *route {
	return &route{method: method, relativePath: relativePath, handlers: handlers}
}

// addToGroups is used by http methods, used to insert handlers to AppRoute.groups.
func (a *AppRoute) addToGroups(method string, relativePath string, handlers []gin.HandlerFunc) {
	if len(handlers) == 0 {
		panic("a route must have at least one handler.")
	}

	r := newRoute(method, relativePath, handlers...)
	for idx, routes := range a.groups {
		if routes[0].method == method {
			a.groups[idx] = append(routes, r)
			return
		}
	}
	a.groups = append(a.groups, []*route{r})
}

// GET registers a new request handle and middleware with the given path and using Get method.
func (a *AppRoute) GET(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodGet, relativePath, handlers)
}

// POST registers a new request handle and middleware with the given path and using Post method.
func (a *AppRoute) POST(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodPost, relativePath, handlers)
}

// DELETE registers a new request handle and middleware with the given path and using Delete method.
func (a *AppRoute) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodDelete, relativePath, handlers)
}

// PATCH registers a new request handle and middleware with the given path and using Patch method.
func (a *AppRoute) PATCH(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodPatch, relativePath, handlers)
}

// PUT registers a new request handle and middleware with the given path and using Put method.
func (a *AppRoute) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodPut, relativePath, handlers)
}

// OPTIONS registers a new request handle and middleware with the given path and using Options method.
func (a *AppRoute) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodOptions, relativePath, handlers)
}

// HEAD registers a new request handle and middleware with the given path and using Head method.
func (a *AppRoute) HEAD(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodHead, relativePath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, DELETE, PATCH, PUT, OPTIONS, HEAD.
func (a *AppRoute) Any(relativePath string, handlers ...gin.HandlerFunc) {
	a.addToGroups(http.MethodGet, relativePath, handlers)
	a.addToGroups(http.MethodPost, relativePath, handlers)
	a.addToGroups(http.MethodDelete, relativePath, handlers)
	a.addToGroups(http.MethodPatch, relativePath, handlers)
	a.addToGroups(http.MethodPut, relativePath, handlers)
	a.addToGroups(http.MethodOptions, relativePath, handlers)
	a.addToGroups(http.MethodHead, relativePath, handlers)
}

// Do handle all registered routes to gin.IRouter using setting of gin.Engine.
func (a *AppRoute) Do() {
	// for all methods
	for _, groupRoutes := range a.groups {
		method := groupRoutes[0].method
		allRoutes := groupRoutes

		// pre handle route param, get max layer
		maxLayer := 0
		for _, route := range allRoutes {
			route.relativePath = strings.Trim(route.relativePath, "/")
			route.parameters = strings.Split(route.relativePath, "/")
			if route.relativePath == "" {
				route.parameters = []string{}
			}
			if len(route.parameters) > maxLayer {
				maxLayer = len(route.parameters)
			}
		}

		// arrange by layer
		layerRoutes := make([][]*route, maxLayer+1)
		for _, route := range allRoutes {
			layerCount := len(route.parameters)
			layerRoutes[layerCount] = append(layerRoutes[layerCount], route)
		}

		// get unexported field
		noRouteHandler := xreflect.GetUnexportedField(reflect.ValueOf(a.engine).Elem().FieldByName("noRoute")).(gin.HandlersChain)
		noMethodHandler := xreflect.GetUnexportedField(reflect.ValueOf(a.engine).Elem().FieldByName("noMethod")).(gin.HandlersChain)

		// build handler
		for layer, routes := range layerRoutes {
			layer := layer
			routes := routes

			// handle route that is empty relativePath first
			if layer == 0 {
				for _, route := range routes { // 0 or 1
					a.router.Handle(method, "", route.handlers...)
				}
				continue
			}

			// build fake path (:_1/:_2/...)
			fakePathSb := strings.Builder{}
			for i := 1; i <= layer; i++ {
				if i > 1 {
					fakePathSb.WriteString("/")
				}
				fakePathSb.WriteString(":_")
				fakePathSb.WriteString(xnumber.Itoa(i)) // :_1
			}
			fakePath := fakePathSb.String()

			// build handler !!! core <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
			handler := func(c *gin.Context) {
				handlers := findRoute(c, routes, fakePath, true)

				// route not found, 404 or 405
				if handlers == nil {
					// if handle 405
					if a.engine.HandleMethodNotAllowed {
						// finding is 405?
						for _, groupRoutes := range a.groups {
							if method == groupRoutes[0].method {
								continue
							}
							// found, use noMethod handler
							if findRoute(c, groupRoutes, fakePath, false) != nil {
								if noMethodHandler == nil {
									c.String(405, "405 method not allowed")
									return
								}
								handlers = noMethodHandler
								break
							}
						}
					}

					// no 405 handler, use 404
					if handlers == nil {
						if noRouteHandler == nil {
							c.String(404, "404 page not found")
							return
						}
						handlers = noRouteHandler
					}
				}

				// run handler (exist | 404 | 405)
				for _, handler := range handlers {
					if !c.IsAborted() {
						handler(c)
					}
				}
			}

			// handle to router
			switch method {
			case http.MethodGet:
				a.router.GET(fakePath, handler)
			case http.MethodPost:
				a.router.POST(fakePath, handler)
			case http.MethodDelete:
				a.router.DELETE(fakePath, handler)
			case http.MethodPatch:
				a.router.PATCH(fakePath, handler)
			case http.MethodPut:
				a.router.PUT(fakePath, handler)
			case http.MethodOptions:
				a.router.OPTIONS(fakePath, handler)
			case http.MethodHead:
				a.router.HEAD(fakePath, handler)
			}

			// print log
			if gin.Mode() == gin.DebugMode {
				for idx, route := range routes {
					pre := "├─"
					if idx == len(routes)-1 {
						pre = "└─"
					}

					funcname := runtime.FuncForPC(reflect.ValueOf(route.handlers[0]).Pointer()).Name()
					fmt.Printf("[XGIN]   %2s %-6s _/%-23s --> %s (--> /%s)\n", pre, method, route.relativePath, funcname, fakePath)
				}
			}
		}
	}
}

// findRoute can find a correspond []*gin.HandleChain through array of route.
// Using fakePath (:_1/:_2...) and do (need to change gin.Context).
func findRoute(c *gin.Context, routes []*route, fakePath string, do bool) []gin.HandlerFunc {
	for _, route := range routes {
		// different parameter size
		if len(c.Params) != len(routes[0].parameters) {
			continue
		}

		// check route accept
		accept := true
		for idx, parameter := range route.parameters {
			if strings.HasPrefix(parameter, ":") { // is `:` path
				continue
			} else {
				if parameter != c.Param("_"+xnumber.Itoa(idx+1)) { // is specific path
					accept = false
					break
				}
			}
		}
		if accept { // accept, found
			// if need to use this route to register handler
			if do {
				for idx, parameter := range route.parameters {
					if strings.HasPrefix(parameter, ":") {
						from := "_" + xnumber.Itoa(idx+1)
						c.Params = append(c.Params, gin.Param{Key: parameter[1:], Value: c.Param(from)}) // set new c.Params
					}
				}
			}

			// set fullPath
			fullPath := strings.TrimSuffix(strings.TrimSuffix(c.FullPath(), fakePath), "/")
			fullPath = fmt.Sprintf("%s/%s", fullPath, route.relativePath)
			xreflect.SetUnexportedField(reflect.ValueOf(c).Elem().FieldByName("fullPath"), fullPath)

			// return
			return route.handlers
		}
	}

	// not found
	return nil
}

/*
	v1.GET(":a", func)
	v1.GET(":a/:b", func) // pass

	v1.GET("a", func)
	v1.GET(":a/:b", func) // panic

	// use this
	v1.GET(":_a", func)
	v1.GET(":_a/:_b", func)
	v1.GET(":_a/:_b/:_c", func)
*/
