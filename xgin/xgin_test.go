package xgin

import (
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

func TestDumpRequest(t *testing.T) {
	app := gin.New()
	app.GET("a", func(c *gin.Context) {
		for _, s := range DumpRequest(c) {
			log.Println(s)
		}
	})
	_ = app.Run(":1234")
}

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

func TestComposite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	testGroup := engine.Group("/test")
	{
		testGroup.GET("", handle)
		testGroup.GET("/:id/:id2", Composite("id",
			M(handle),          // /?/?
			P("test", handle2), // /test/?
			P("test2", Composite("id2",
				M(handle2),          // /test2/?
				P("test", handle3),  // /test2/test
				P("test2", handle3), // /test2/test2
				N(handle4),          // /test2/0
			)),
			N(handle4, // /0/?
				Composite("id2",
					M(handle5),         // /0/?
					P("test", handle6), // /0/test
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
		ctxGroup.GET("/:id", Composite("id",
			M(func(c *gin.Context) {
				log.Println(11)
			}, func(c *gin.Context) {
				log.Println(12)
				c.Abort()
			}),
			P("test", func(c *gin.Context) {
				log.Println(21)
			}, func(c *gin.Context) {
				log.Println(22)
			}),
			P("test2", func(c *gin.Context) {
				log.Println(31)
			}, func(c *gin.Context) {
				log.Println(32)
				c.Abort()
			}, func(c *gin.Context) {
				log.Println(33)
				c.Abort()
			}),
			N(func(c *gin.Context) {
				log.Println(41)
				c.Abort()
			}, func(c *gin.Context) {
				log.Println(42)
			}),
		))
	}
	_ = engine.Run(":1234")
}
