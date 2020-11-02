package xfiber

import (
	"fmt"
	"github.com/gofiber/fiber"
	logrus2 "github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
	"time"
)

func TestDumpRequest(t *testing.T) {
	app := fiber.New()
	app.Get("a", func(c *fiber.Ctx) {
		for _, s := range DumpRequest(c) {
			log.Println(s)
		}
	})
	_ = app.Listen("1234")
}

func TestBuildErrorDto(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) {
		defer func() {
			if err := recover(); err != nil {
				e := BuildErrorDto(err, c, 2, true)
				e.Others = map[string]interface{}{"a": "b"}
				_ = c.JSON(e)
			}
		}()
		c.Next()
	})
	app.Get("panic", func(c *fiber.Ctx) {
		panic("test panic")
	})
	app.Get("error", func(c *fiber.Ctx) {
		_ = c.JSON(BuildBasicErrorDto(fmt.Errorf("test error"), c))
	})
	_ = app.Listen("1234")
}

func TestLogger(t *testing.T) {
	app := fiber.New()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	logrus := logrus2.New()
	logrus.SetFormatter(&logrus2.TextFormatter{})

	app.Use(PprofHandler())
	app.Use(func(c *fiber.Ctx) {
		start := time.Now()
		c.Next()
		WithLogger(logger, start, c, "12345")
		WithLogrus(logrus, start, c, &LoggerExtra{
			OtherString: "12345",
			OtherFields: nil,
		})
	})

	_ = app.Listen("1234")
}
