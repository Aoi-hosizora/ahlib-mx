package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
	"strings"
)

/*
	v1.GET(":a", func)
	v1.GET(":a/:b", func) // pass

	v1.GET("a", func)
	v1.GET(":a/:b", func) // panic

	v1.GET(":_a", func)
	v1.GET(":_a/:_b", func)
	v1.GET(":_a/:_b/:_c", func)

	xgin.XXX(group,
		xgin.Route(":a", func, func),
		xgin.Route("a/:b", func, func),
		xgin.Route(":a/:b", func, func),
	)
	->
	group.GET(":_a", func)
	group.GET(":_a/:_b", func)
*/

type route struct {
	relativePath string
	parameters   []string
	handlers     []gin.HandlerFunc
}

func Route(relativePath string, handlers ...gin.HandlerFunc) *route {
	return &route{relativePath: relativePath, handlers: handlers}
}

var _ = METHOD
var _ = ANY
var _ = GET
var _ = POST
var _ = PUT
var _ = DELETE
var _ = PATCH
var _ = HEAD
var _ = OPTIONS

func METHOD(app *gin.Engine, method func(string, ...gin.HandlerFunc) gin.IRoutes, routes ...*route) {
	// arrange by layer
	routeGroups := make(map[int][]*route)
	for _, r := range routes {
		r.relativePath = strings.TrimPrefix(r.relativePath, "/")
		r.parameters = strings.Split(r.relativePath, "/")
		layerCount := len(r.parameters)
		if _, ok := routeGroups[layerCount]; !ok {
			routeGroups[layerCount] = []*route{r}
		} else {
			routeGroups[layerCount] = append(routeGroups[layerCount], r)
		}
	}

	// get unexported field
	noMethodHandler := xreflect.GetUnexportedField(reflect.ValueOf(app).Elem().FieldByName("noMethod")).(gin.HandlersChain)

	// build handler
	for layerCount, routes := range routeGroups {
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

		// build handler
		path := pathSb.String()
		handler := func(c *gin.Context) {
			handlers := findRoute(c, routes, path, true)
			if handlers == nil {
				handlers = noMethodHandler
			}

			for _, handler := range handlers {
				if !c.IsAborted() {
					handler(c)
				}
			}
		}

		// handle
		method(path, handler)
	}
}

func findRoute(c *gin.Context, routes []*route, fakePath string, do bool) []gin.HandlerFunc {
	if routes == nil {
		return nil
	}

	// find route
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

func ANY(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.Any, routes...)
}

func GET(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.GET, routes...)
}

func POST(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.POST, routes...)
}

func DELETE(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.DELETE, routes...)
}

func PATCH(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.PATCH, routes...)
}

func PUT(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.PUT, routes...)
}

func OPTIONS(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.OPTIONS, routes...)
}

func HEAD(app *gin.Engine, group gin.IRouter, routes ...*route) {
	METHOD(app, group.HEAD, routes...)
}
