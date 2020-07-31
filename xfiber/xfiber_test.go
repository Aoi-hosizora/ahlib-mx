package xfiber

import (
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
