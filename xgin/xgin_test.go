package xgin

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin/xroute"
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

func handle(c *gin.Context) {
	log.Println(c.Request.URL.Path, "!", c.Param("id"))
}

func handle2(c *gin.Context) {
	log.Println(c.Request.URL.Path, "!!", c.Param("id"))
}

func handle3(c *gin.Context) {
	log.Println(c.Request.URL.Path, "!!!", c.Param("id"))
}

func handle4(c *gin.Context) {
	log.Println(c.Request.URL.Path, "!!!!", c.Param("id"))
}

func handle5(c *gin.Context) {
	log.Println(c.Request.URL.Path, "!!!!!", c.Param("id"))
}

func handle6(c *gin.Context) {
	log.Println(c.Request.URL.Path, "!!!!!!", c.Param("id"))
}

func TestGin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	testGroup := engine.Group("/test")
	{
		testGroup.GET("", handle)
		testGroup.GET("/:id/:id2", xroute.Composite("id",
			xroute.M(handle),          // /?/?
			xroute.P("test", handle2), // /test/?
			xroute.P("test2", xroute.Composite("id2",
				xroute.M(handle2),          // /test2/?
				xroute.P("test", handle3),  // /test2/test
				xroute.P("test2", handle3), // /test2/test2
				xroute.N(handle4),          // /test2/0
			)),
			xroute.N(handle4, // /0/?
				xroute.Composite("id2",
					xroute.M(handle5),         // /0/?
					xroute.P("test", handle6), // /0/test
				),
			),
		))
	}
	ctxGroup := engine.Group("/ctx")
	{
		ctxGroup.GET("", func(c *gin.Context) {
			log.Println(1)
		}, func(c *gin.Context) {
			log.Println(2)
		})
		ctxGroup.GET("/:id", xroute.Composite("id",
			xroute.M(func(c *gin.Context) {
				log.Println(11)
			}, func(c *gin.Context) {
				log.Println(12)
				c.Abort()
			}),
			xroute.P("test", func(c *gin.Context) {
				log.Println(21)
			}, func(c *gin.Context) {
				log.Println(22)
			}),
			xroute.P("test2", func(c *gin.Context) {
				log.Println(31)
			}, func(c *gin.Context) {
				log.Println(32)
				c.Abort()
			}, func(c *gin.Context) {
				log.Println(33)
				c.Abort()
			}),
			xroute.N(func(c *gin.Context) {
				log.Println(41)
				c.Abort()
			}, func(c *gin.Context) {
				log.Println(42)
			}),
		))
	}
	_ = engine.Run(":1234")
}
