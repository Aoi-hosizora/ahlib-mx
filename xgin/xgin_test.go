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

func TestParam(t *testing.T) {
	app := gin.New()
	app.GET(":a/:b", Param(func(c *gin.Context) {
		log.Println(c.Param("a"))
		log.Println(c.Param("b"))
		log.Println(c.Param("c"))
		log.Println(c.Param("d"))
		log.Println(c.Param("e"))
	}, Padd("c", "a"), Padd("d", "b"), Pdel("a"), Pdel("b"), Pdel("b")))
	_ = app.Run(":1234")
}

func TestRoute(t *testing.T) {
	app := gin.New()
	app.GET(":a/:b", Composite("a",
		P("p", func(c *gin.Context) {
			log.Println("P1", c.Param("a"), c.Param("b"))
		}, func(c *gin.Context) {
			log.Println("P2", c.Param("a"), c.Param("b"))
			c.Abort()
		}, func(c *gin.Context) {
			log.Println("P3", c.Param("a"), c.Param("b"))
		}),

		P("1", func(c *gin.Context) {
			log.Println("P4", c.Param("a"), c.Param("b"))
		}, func(c *gin.Context) {
			log.Println("P5", c.Param("a"), c.Param("b"))
			c.Abort()
		}, func(c *gin.Context) {
			log.Println("P6", c.Param("a"), c.Param("b"))
		}),

		I(func(c *gin.Context) {
			log.Println("I1", c.Param("a"), c.Param("b"))
		}, func(c *gin.Context) {
			log.Println("I2", c.Param("a"), c.Param("b"))
			c.Abort()
		}, func(c *gin.Context) {
			log.Println("I3", c.Param("a"), c.Param("b"))
		}),

		F(func(c *gin.Context) {
			log.Println("F1", c.Param("a"), c.Param("b"))
		}, func(c *gin.Context) {
			log.Println("F2", c.Param("a"), c.Param("b"))
			c.Abort()
		}, func(c *gin.Context) {
			log.Println("F3", c.Param("a"), c.Param("b"))
		}),

		M(func(c *gin.Context) {
			log.Println("M1", c.Param("a"), c.Param("b"))
		}, func(c *gin.Context) {
			log.Println("M2", c.Param("a"), c.Param("b"))
			c.Abort()
		}, func(c *gin.Context) {
			log.Println("M3", c.Param("a"), c.Param("b"))
		}),
	))
	_ = app.Run(":1234")
}
