package xgin

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
)

type ParamOption struct {
	From   string
	To     string
	Delete bool
}

// Padd create an instance of ParamOption (delete: false).
func Padd(from string, to string) *ParamOption {
	return &ParamOption{From: from, To: to, Delete: false}
}

// Pdel create an instance of ParamOption (delete: true)
func Pdel(from string) *ParamOption {
	return &ParamOption{From: from, Delete: true}
}

// Param copy some route param to new param in gin.Context.
func Param(handler func(c *gin.Context), params ...*ParamOption) func(c *gin.Context) {
	if len(params) == 0 {
		panic("a param mapper route must have at least two params string.")
	}

	added := make([]*ParamOption, 0)
	deleted := make([]string, 0)
	for _, param := range params {
		if !param.Delete {
			added = append(added, param)
		} else {
			deleted = append(deleted, param.From)
		}
	}

	indexOf := func(params []gin.Param, key string) int {
		idx := -1
		for i, param := range params {
			if param.Key == key {
				idx = i
			}
		}
		return idx
	}

	return func(c *gin.Context) {
		// add
		for _, param := range added {
			c.Params = append(c.Params, gin.Param{
				Key:   param.From,
				Value: c.Param(param.To),
			})
		}

		// del
		for _, del := range deleted {
			idx := indexOf(c.Params, del)
			for idx != -1 {
				if len(c.Params) == idx+1 {
					c.Params = c.Params[:idx]
				} else {
					c.Params = append(c.Params[:idx], c.Params[idx+1:]...)
				}
				idx = indexOf(c.Params, del)
			}
		}

		// handler
		if !c.IsAborted() {
			handler(c)
		}
	}
}

// An interface used for Composite parameter.
type CompositeHandler interface {
	Check(param string) bool
	Do(c *gin.Context)
}

// CompositeHandler for main handler.
type MainHandler struct {
	Handlers []gin.HandlerFunc
}

// M Create an instance of MainHandler.
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

// P Create an instance of PrefixHandler.
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

// I Create an instance of IntegerHandler.
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

// F Create an instance of FloatHandler.
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
