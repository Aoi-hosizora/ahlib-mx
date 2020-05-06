package xroute

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Option interface {
	Check(param string) bool
	Handle(c *gin.Context)
}

// Option for main handler
type MainHandler struct {
	Handlers []gin.HandlerFunc
}

func M(handlers ...gin.HandlerFunc) *MainHandler {
	return &MainHandler{Handlers: handlers}
}

func (m *MainHandler) Handle(c *gin.Context) {
	for _, handle := range m.Handlers {
		if c.IsAborted() {
			return
		}
		handle(c)
	}
}

// Option for specific prefix route
type PrefixOption struct {
	Prefix   string
	Handlers []gin.HandlerFunc
}

func P(prefix string, handlers ...gin.HandlerFunc) *PrefixOption {
	return &PrefixOption{Prefix: prefix, Handlers: handlers}
}

func (p *PrefixOption) Check(param string) bool {
	return p.Prefix == param
}

func (p *PrefixOption) Handle(c *gin.Context) {
	for _, handle := range p.Handlers {
		if c.IsAborted() {
			return
		}
		handle(c)
	}
}

// Option for numeric route
type NumericOption struct {
	Handlers []gin.HandlerFunc
}

func N(handlers ...gin.HandlerFunc) *NumericOption {
	return &NumericOption{Handlers: handlers}
}

func (n *NumericOption) Check(param string) bool {
	_, err := strconv.Atoi(param)
	return err == nil
}

func (n *NumericOption) Handle(c *gin.Context) {
	for _, handle := range n.Handlers {
		if c.IsAborted() {
			return
		}
		handle(c)
	}
}

// panic: 'xxx' in new path '/user/xxx' conflicts with existing wildcard ':id' in existing Prefix '/user/:id' [recovered]
//
// Example: M() P() N()
//     testGroup.GET("", handle)
//     testGroup.GET("/:id/:id2", Composite("id",
//         M(handle),                   // /?/?
//         P("test", handle2),          // /test/?
//         P("test2", Composite("id2",
//             M(handle2),              // /test2/?
//             P("test", handle3),      // /test2/test
//             P("test2", handle3),     // /test2/test2
//             N(handle4),              // /test2/0
//         )),
//         N(handle4,                   // /0/?
//             Composite("id2",
//                 M(handle5),          // /0/?
//                 P("test", handle6),  // /0/test
//             ),
//         ),
//     ))
func Composite(key string, main *MainHandler, options ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		subPath := c.Param(key)
		for _, option := range options {
			if option.Check(subPath) {
				option.Handle(c)
				return
			}
		}
		main.Handle(c)
	}
}
