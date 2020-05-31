package xfiber

import (
	"github.com/gofiber/fiber"
	"strconv"
)

type CompositeOption interface {
	Check(param string) bool
	Do(c *fiber.Ctx)
}

// Option for main handler
type MainHandler struct {
	Handlers []fiber.Handler
}

// Option for specific prefix route
type PrefixOption struct {
	Prefix   string
	Handlers []fiber.Handler
}

// Option for numeric route
type NumericOption struct {
	Handlers []fiber.Handler
}

func M(handlers ...fiber.Handler) *MainHandler {
	return &MainHandler{Handlers: handlers}
}

func P(prefix string, handlers ...fiber.Handler) *PrefixOption {
	return &PrefixOption{Prefix: prefix, Handlers: handlers}
}

func N(handlers ...fiber.Handler) *NumericOption {
	return &NumericOption{Handlers: handlers}
}

func (p *PrefixOption) Check(param string) bool {
	return p.Prefix == param
}

func (n *NumericOption) Check(param string) bool {
	_, err := strconv.Atoi(param)
	return err == nil
}

func (m *MainHandler) Do(c *fiber.Ctx) {
	for _, handle := range m.Handlers {
		handle(c)
		if c.Error() != nil {
			break
		}
	}
	// c.Next(c.Error())
}

func (p *PrefixOption) Do(c *fiber.Ctx) {
	for _, handle := range p.Handlers {
		handle(c)
		if c.Error() != nil {
			break
		}
	}
	// c.Next(c.Error())
}

func (n *NumericOption) Do(c *fiber.Ctx) {
	for _, handle := range n.Handlers {
		handle(c)
		if c.Error() != nil {
			break
		}
	}
	// c.Next(c.Error())
}

func Composite(key string, main *MainHandler, options ...CompositeOption) fiber.Handler {
	return func(c *fiber.Ctx) {
		subPath := c.Params(key)
		do := main.Do
		for _, option := range options {
			if option.Check(subPath) {
				do = option.Do
				break
			}
		}
		do(c)
		// c.Next(c.Error())
	}
}
