package xgin

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
)

// An interface used for Composite parameter.
type CompositeHandler interface {
	Check(param string) bool
	Do(c *gin.Context)
}

// CompositeHandler for main handler.
type MainHandler struct {
	Handlers []gin.HandlerFunc
}

// Create an instance of MainHandler.
func M(handlers ...gin.HandlerFunc) *MainHandler {
	return &MainHandler{Handlers: handlers}
}

func (m *MainHandler) Check(string) bool {
	return true
}

func (m *MainHandler) Do(c *gin.Context) {
	for _, handle := range m.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

// CompositeHandler for specific prefix.
type PrefixHandler struct {
	Prefix   string
	Handlers []gin.HandlerFunc
}

// Create an instance of PrefixHandler.
func P(prefix string, handlers ...gin.HandlerFunc) *PrefixHandler {
	return &PrefixHandler{Prefix: prefix, Handlers: handlers}
}

func (p *PrefixHandler) Check(param string) bool {
	return p.Prefix == param
}

func (p *PrefixHandler) Do(c *gin.Context) {
	for _, handle := range p.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

// CompositeHandler for int64 parameter.
type IntegerHandler struct {
	Handlers []gin.HandlerFunc
}

// Create an instance of IntegerHandler.
func I(handlers ...gin.HandlerFunc) *IntegerHandler {
	return &IntegerHandler{Handlers: handlers}
}

func (n *IntegerHandler) Check(param string) bool {
	_, err := xnumber.Atoi64(param)
	return err == nil
}

func (n *IntegerHandler) Do(c *gin.Context) {
	for _, handle := range n.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

// CompositeHandler for float64 parameter.
type FloatHandler struct {
	Handlers []gin.HandlerFunc
}

// Create an instance of FloatHandler.
func F(handlers ...gin.HandlerFunc) *FloatHandler {
	return &FloatHandler{Handlers: handlers}
}

func (f *FloatHandler) Check(param string) bool {
	_, err := xnumber.Atof64(param)
	return err == nil
}

func (f *FloatHandler) Do(c *gin.Context) {
	for _, handle := range f.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

// Composite some CompositeHandler for `wildcard route`. This route will check handlers in order.
// 	panic: 'xxx' in new path '/user/xxx' conflicts with existing wildcard ':id' in existing Prefix '/user/:id' [recovered]
func Composite(key string, handlers ...CompositeHandler) gin.HandlerFunc {
	if len(handlers) == 0 {
		panic("a composite route must have at least one CompositeHandler.")
	}

	return func(c *gin.Context) {
		subPath := c.Param(key)
		do := func(c *gin.Context) {}
		for _, option := range handlers {
			if option.Check(subPath) {
				do = option.Do
				break
			}
		}
		do(c)
	}
}
