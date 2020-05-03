package xgin

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PathOption interface {
	Check(param string) bool
	Handle(c *gin.Context)
}

// Option for string route Prefix
type PrefixOption struct {
	Prefix   string
	Handlers []gin.HandlerFunc
}

func NewPrefixOption(prefix string, handlers ...gin.HandlerFunc) *PrefixOption {
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

// Option for numeric route Prefix
type NumericOption struct {
	Handlers []gin.HandlerFunc
}

func NewNumericOption(handlers ...gin.HandlerFunc) *NumericOption {
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
// Example:
//     testGroup.GET("", handle)
//     testGroup.GET("/:id/:id2", MultiplePaths(
//         "id", handle,                             // /?/?
//         NewPrefixOption("test", handle2),         // /test/?
//         NewPrefixOption("test2", MultiplePaths(
//             "id2", handle2,                       // /test2/?
//             NewPrefixOption("test", handle3),     // /test2/test
//             NewPrefixOption("test2", handle3),    // /test2/test2
//             NewNumericOption(handle4),            // /test2/0
//         )),
//         NewNumericOption(handle4,                 // /0/?
//             MultiplePaths(
//                 "id2", handle5,                   // /0/?
//                 NewPrefixOption("test", handle6), // /0/test
//             ),
//         ),
//     ))
func MultiplePaths(key string, mainHandler gin.HandlerFunc, options ...PathOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		subPath := c.Param(key)
		for _, option := range options {
			if option.Check(subPath) {
				option.Handle(c)
				return
			}
		}
		mainHandler(c)
	}
}
