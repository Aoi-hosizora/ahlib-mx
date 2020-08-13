package xfiber

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gofiber/fiber"
	"log"
	"testing"
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
				j, _ := json.Marshal(BuildErrorDto(err, c, 2, true))
				fmt.Println(xstring.PrettifyJson(string(j), 4, " "))
			}
		}()
		c.Next()
	})
	app.Get("", func(c *fiber.Ctx) {
		panic("test panic")
	})
	_ = app.Listen("1234")
}
