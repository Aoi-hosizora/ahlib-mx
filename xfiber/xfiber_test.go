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
				_ = c.JSON(BuildErrorDto(err, c, 2, true))
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

func TestBinding(t *testing.T) {
	app := fiber.New()
	_ = EnableRegexpBinding()
	_ = EnableRFC3339DateBinding()
	_ = EnableRFC3339DateTimeBinding()

	type st struct {
		A string `validate:"regexp=^[abc]+$"`
		B string `validate:"date"`
		C string `validate:"datetime"`
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	logrus := logrus2.New()

	app.Use(PprofHandler())
	app.Use(func(c *fiber.Ctx) {
		start := time.Now()
		c.Next()
		WithLogger(logger, start, c)
		WithLogrus(logrus, start, c)
	})

	app.Get("", func(ctx *fiber.Ctx) {
		a := ctx.Query("a")
		b := ctx.Query("b")
		c := ctx.Query("c")
		st := &st{A: a, B: b, C: c}
		if err := Struct(st); err != nil {
			ctx.SendString(err.Error())
		}
		if err := Var(st.A, "regexp=^[abc]+$"); err != nil {
			ctx.SendString(err.Error())
		}
	})

	_ = app.Listen("1234")
}
