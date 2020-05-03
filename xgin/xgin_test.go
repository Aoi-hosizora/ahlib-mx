package xgin

import (
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
		testGroup.GET("/:id/:id2", MultiplePaths(
			"id", handle, // /?/?
			NewPrefixOption("test", handle2), // /test/?
			NewPrefixOption("test2", MultiplePaths(
				"id2", handle2, // /test2/?
				NewPrefixOption("test", handle3),  // /test2/test
				NewPrefixOption("test2", handle3), // /test2/test2
				NewNumericOption(handle4),         // /test2/0
			)),
			NewNumericOption(
				handle4, // /0/?
				MultiplePaths(
					"id2", handle5, // /0/?
					NewPrefixOption("test", handle6), // /0/test
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
		ctxGroup.GET("/:id", MultiplePaths("id", func(c *gin.Context) {
			log.Println(11)
			c.Abort()
		}, NewPrefixOption("test", func(c *gin.Context) {
			log.Println(12)
		}, func(c *gin.Context) {
			log.Println(13)
		}), NewPrefixOption("test2", func(c *gin.Context) {
			log.Println(21)
		}, func(c *gin.Context) {
			log.Println(22)
			c.Abort()
		}, func(c *gin.Context) {
			log.Println(23)
			c.Abort()
		}), NewNumericOption(func(c *gin.Context) {
			log.Println(31)
			c.Abort()
		}, func(c *gin.Context) {
			log.Println(32)
		})))
	}
	_ = engine.Run(":1234")
}
