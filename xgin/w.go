package xgin

import (
	"github.com/gin-gonic/gin"
)

type HandlerFuncW func(c *gin.Context) (int, interface{})

func JsonW(fn HandlerFuncW) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data := fn(c)
		c.JSON(code, data)
	}
}

func XmlW(fn HandlerFuncW) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data := fn(c)
		c.XML(code, data)
	}
}
