package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// AppRoute stores a group of methods (group of routes).
type AppRoute struct {
	engine *gin.Engine
	router gin.IRouter
	groups map[string][]*route
}

// NewAppRoute create an instance of AppRoute.
func NewAppRoute(engine *gin.Engine, router gin.IRouter) *AppRoute {
	return &AppRoute{
		engine: engine,
		router: router,
		groups: map[string][]*route{},
	}
}

// route represents a route in AppRoute, include handlers and relativePath.
type route struct {
	relativePath string
	parameters   []string // used later
	handlers     []gin.HandlerFunc
}

// newRoute create an instance of route, panic if relativePath is empty.
func newRoute(relativePath string, handlers ...gin.HandlerFunc) *route {
	if relativePath == "" {
		panic("AppRoute only allow to create not empty route")
	}
	return &route{relativePath: relativePath, handlers: handlers}
}

// addToGroups is used by http methods, used to insert handlers to AppRoute.groups.
func (a *AppRoute) addToGroups(method string, relativePath string, handlers []gin.HandlerFunc) {
	r := newRoute(relativePath, handlers...)
	if _, ok := a.groups[method]; !ok {
		a.groups[method] = []*route{r}
	} else {
		a.groups[method] = append(a.groups[method], r)
	}
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
	for method, routes := range a.groups {
		method := method
		routes := routes

		// arrange by layer
		layerRoutes := make(map[int][]*route)
		for _, r := range routes {
			r.relativePath = strings.TrimPrefix(strings.TrimSuffix(r.relativePath, "/"), "/")
			r.parameters = strings.Split(r.relativePath, "/")
			layerCount := len(r.parameters)
			if _, ok := layerRoutes[layerCount]; !ok {
				layerRoutes[layerCount] = []*route{r}
			} else {
				layerRoutes[layerCount] = append(layerRoutes[layerCount], r)
			}
		}

		// get unexported field
		noRouteHandler := xreflect.GetUnexportedField(reflect.ValueOf(a.engine).Elem().FieldByName("noRoute")).(gin.HandlersChain)
		noMethodHandler := xreflect.GetUnexportedField(reflect.ValueOf(a.engine).Elem().FieldByName("noMethod")).(gin.HandlersChain)

		// build handler
		for layerCount, routes := range layerRoutes {
			routes := routes
			// build parameter path sting
			pathSb := strings.Builder{}
			for i := 1; i <= layerCount; i++ {
				if i > 1 {
					pathSb.WriteString("/")
				}
				pathSb.WriteString(":_")
				pathSb.WriteString(xnumber.Itoa(i))
			}

			// build handler !!! core
			path := pathSb.String()

			handler := func(c *gin.Context) {
				handlers := findRoute(c, routes, path, true)

				// route not found
				if handlers == nil {
					// if handle 405
					if a.engine.HandleMethodNotAllowed {
						// finding is 405?
						for otherMethod, routes := range a.groups {
							if method == otherMethod {
								continue
							}
							// found, use noMethod handler
							if findRoute(c, routes, path, false) != nil {
								log.Println(otherMethod, path) // TODO wrong 405
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
				a.router.GET(path, handler)
			case http.MethodPost:
				a.router.POST(path, handler)
			case http.MethodDelete:
				a.router.DELETE(path, handler)
			case http.MethodPatch:
				a.router.PATCH(path, handler)
			case http.MethodPut:
				a.router.PUT(path, handler)
			case http.MethodOptions:
				a.router.OPTIONS(path, handler)
			case http.MethodHead:
				a.router.HEAD(path, handler)
			}
		}
	}
}

// findRoute can find a correspond []*gin.HandleChain through array of route.
// Using fakePath (:_1/:_2...) and do (need to change gin.Context).
func findRoute(c *gin.Context, routes []*route, fakePath string, do bool) []gin.HandlerFunc {
	if routes == nil {
		return nil
	}

	// find route !!!! core
	for _, route := range routes {
		accept := true
		for idx, parameter := range route.parameters {
			if strings.HasPrefix(parameter, ":") { // is `:` path
				continue
			} else {
				if parameter != c.Param("_"+xnumber.Itoa(idx+1)) { // is specific path
					accept = false
				}
			}
		}
		if accept { // accept, found
			if do { // if need to use this route to register handler
				for idx, parameter := range route.parameters {
					if strings.HasPrefix(parameter, ":") {
						from := "_" + xnumber.Itoa(idx+1)
						c.Params = append(c.Params, gin.Param{Key: parameter[1:], Value: c.Param(from)}) // set new c.Params
					}
				}
			}

			// set unexported field

			// fullPath
			fullPath := xreflect.GetUnexportedField(reflect.ValueOf(c).Elem().FieldByName("fullPath")).(string)
			fullPath = strings.TrimSuffix(strings.TrimSuffix(fullPath, fakePath), "/")
			fullPath = fmt.Sprintf("%s/%s", fullPath, route.relativePath)
			xreflect.SetUnexportedField(reflect.ValueOf(c).Elem().FieldByName("fullPath"), fullPath) // TODO use reflect

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
