package xgin

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
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

func METHOD(method func(string, ...gin.HandlerFunc) gin.IRoutes, routes ...*route) {
	routeGroups := make(map[int][]*route)
	for _, r := range routes {
		r.relativePath = strings.TrimLeft(r.relativePath, "/")
		r.parameters = strings.Split(r.relativePath, "/")
		layerCount := len(r.parameters)
		if _, ok := routeGroups[layerCount]; !ok {
			routeGroups[layerCount] = []*route{r}
		} else {
			routeGroups[layerCount] = append(routeGroups[layerCount], r)
		}
	}

	for layerCount, routes := range routeGroups {
		routes := routes
		pathSb := strings.Builder{}
		for i := 1; i <= layerCount; i++ {
			if i > 1 {
				pathSb.WriteString("/")
			}
			pathSb.WriteString(":_")
			pathSb.WriteString(xnumber.Itoa(i))
		}

		path := pathSb.String()
		handler := func(c *gin.Context) {
			handlers := findRoute(c, routes, true)
			if handlers != nil {
				for _, handler := range handlers {
					if !c.IsAborted() {
						handler(c)
					}
				}
			} else {
				// not found
			}
		}

		method(path, handler)
	}
}

func findRoute(c *gin.Context, routes []*route, do bool) []gin.HandlerFunc {
	if routes == nil {
		return nil
	}
	for _, route := range routes {
		accept := true
		for idx, parameter := range route.parameters {
			if strings.HasPrefix(parameter, ":") {
				continue
			} else {
				if parameter != c.Param("_"+xnumber.Itoa(idx+1)) {
					accept = false
				}
			}
		}
		if accept {
			if do {
				for idx, parameter := range route.parameters {
					if strings.HasPrefix(parameter, ":") {
						from := "_" + xnumber.Itoa(idx+1)
						c.Params = append(c.Params, gin.Param{Key: parameter[1:], Value: c.Param(from)})
					}
				}
			}
			return route.handlers
		}
	}
	return nil
}

func ANY(group gin.IRouter, routes ...*route) {
	METHOD(group.Any, routes...)
}

func GET(group gin.IRouter, routes ...*route) {
	METHOD(group.GET, routes...)
}

func POST(group gin.IRouter, routes ...*route) {
	METHOD(group.POST, routes...)
}

func DELETE(group gin.IRouter, routes ...*route) {
	METHOD(group.DELETE, routes...)
}

func PATCH(group gin.IRouter, routes ...*route) {
	METHOD(group.PATCH, routes...)
}

func PUT(group gin.IRouter, routes ...*route) {
	METHOD(group.PUT, routes...)
}

func OPTIONS(group gin.IRouter, routes ...*route) {
	METHOD(group.OPTIONS, routes...)
}

func HEAD(group gin.IRouter, routes ...*route) {
	METHOD(group.HEAD, routes...)
}
