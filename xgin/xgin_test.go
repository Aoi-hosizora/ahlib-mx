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
