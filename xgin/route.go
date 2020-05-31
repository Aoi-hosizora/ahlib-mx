package xgin

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type CompositeOption interface {
	Check(param string) bool
	Do(c *gin.Context)
}

// Option for main handler
type MainHandler struct {
	Handlers []gin.HandlerFunc
}

// Option for specific prefix route
type PrefixOption struct {
	Prefix   string
	Handlers []gin.HandlerFunc
}

// Option for numeric route
type NumericOption struct {
	Handlers []gin.HandlerFunc
}

func M(handlers ...gin.HandlerFunc) *MainHandler {
	return &MainHandler{Handlers: handlers}
}

func P(prefix string, handlers ...gin.HandlerFunc) *PrefixOption {
	return &PrefixOption{Prefix: prefix, Handlers: handlers}
}

func N(handlers ...gin.HandlerFunc) *NumericOption {
	return &NumericOption{Handlers: handlers}
}

func (p *PrefixOption) Check(param string) bool {
	return p.Prefix == param
}

func (n *NumericOption) Check(param string) bool {
	_, err := strconv.Atoi(param)
	return err == nil
}

func (m *MainHandler) Do(c *gin.Context) {
	for _, handle := range m.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

func (p *PrefixOption) Do(c *gin.Context) {
	for _, handle := range p.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

func (n *NumericOption) Do(c *gin.Context) {
	for _, handle := range n.Handlers {
		handle(c)
		if c.IsAborted() {
			return
		}
	}
}

// panic: 'xxx' in new path '/user/xxx' conflicts with existing wildcard ':id' in existing Prefix '/user/:id' [recovered]
func Composite(key string, main *MainHandler, options ...CompositeOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		subPath := c.Param(key)
		do := main.Do
		for _, option := range options {
			if option.Check(subPath) {
				do = option.Do
				break
			}
		}
		do(c)
	}
}
