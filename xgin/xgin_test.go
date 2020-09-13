package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logrus2 "github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
	"time"
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

func TestBuildErrorDto(t *testing.T) {
	app := gin.New()
	app.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e := BuildErrorDto(err, c, 2, true)
				e.Others = map[string]interface{}{"a": "b"}

				log.Println(BuildFullErrorDto(err, c, e.Others, 2, false))
				c.JSON(200, e)
			}
		}()
		c.Next()
	})
	app.GET("panic", func(c *gin.Context) {
		panic("test panic")
	})
	app.GET("error", func(c *gin.Context) {
		c.JSON(200, BuildBasicErrorDto(fmt.Errorf("test error"), c))
	})
	_ = app.Run(":1234")
}

func TestLogger(t *testing.T) {
	app := gin.New()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	logrus := logrus2.New()

	PprofWrap(app)
	app.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		WithLogger(logger, start, c, "12345")
		WithLogrus(logrus, start, c, &LoggerExtra{
			OtherString: "12345",
			OtherFields: nil,
		})
	})

	_ = app.Run(":1234")
}

func TestBinding(t *testing.T) {
	app := gin.New()
	_ = EnableRegexpBinding()
	_ = EnableRFC3339DateBinding()
	_ = EnableRFC3339DateTimeBinding()

	type st struct {
		A string `binding:"regexp=^[abc]+$"`
		B string `binding:"date"`
		C string `binding:"datetime"`
	}

	app.GET("", func(ctx *gin.Context) {
		st := &st{}
		err := ctx.ShouldBindQuery(st)
		if err != nil {
			ctx.String(200, err.Error())
		}
	})

	_ = app.Run(":1234")
}

func TestRoute(t *testing.T) {
	app := gin.New()
	app.NoMethod(func(c *gin.Context) { c.String(200, "%s not found", c.Request.URL.String()) })
	GET(app, app.Group("v1"),
		Route("a", func(c *gin.Context) { log.Println(1, c.FullPath()) }),
		Route("b", func(c *gin.Context) { log.Println(2, c.FullPath()) }),
		Route(":a", func(c *gin.Context) { log.Println(3, c.FullPath(), "|", c.Param("a")) }),
		Route("a/b", func(c *gin.Context) { log.Println(4, c.FullPath()) }),
		Route("c/d", func(c *gin.Context) { log.Println(5, c.FullPath()) }),
		Route(":a/:b", func(c *gin.Context) { log.Println(6, c.FullPath(), "|", c.Param("a"), "|", c.Param("b")) }),
		Route("a/b/c", func(c *gin.Context) { log.Println(7, c.FullPath()) }),
		Route("d/e/f", func(c *gin.Context) { log.Println(8, c.FullPath()) }),
		Route(":a/:b/:c", func(c *gin.Context) { log.Println(9, c.FullPath(), "|", c.Param("a"), "|", c.Param("b"), "|", c.Param("c")) }),
		Route("a/b/c/d", func(c *gin.Context) { log.Println(10, c.FullPath()) }),
	)
	_ = app.Run(":1234")
}
