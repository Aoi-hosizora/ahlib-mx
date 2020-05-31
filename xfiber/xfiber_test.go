package xfiber

import (
	"github.com/gofiber/fiber"
	"log"
	"testing"
)

func TestComposite(t *testing.T) {
	app := fiber.New()
	app.Get("/:id", Composite("id",
		M(func(c *fiber.Ctx) {
			log.Println(11, c.Params("id"))
		}, func(c *fiber.Ctx) {
			log.Println(12, c.Params("id"))
		}),
		P("a", func(c *fiber.Ctx) {
			log.Println(21, c.Params("id"))
		}, func(c *fiber.Ctx) {
			log.Println(22, c.Params("id"))
			c.Next()
		}, func(c *fiber.Ctx) {
			log.Println(23, c.Params("id"))
		}),
		P("b", func(c *fiber.Ctx) {
			log.Println(31, c.Params("id"))
		}, func(c *fiber.Ctx) {
			log.Println(32, c.Params("id"))
		}),
		N(func(c *fiber.Ctx) {
			log.Println(41, c.Params("id"))
		}),
	), func(c *fiber.Ctx) {
		log.Println(51, c.Params("id"))
	})
	_ = app.Listen("1234")
}
