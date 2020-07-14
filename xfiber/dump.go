package xfiber

import (
	"github.com/gofiber/fiber"
	"strings"
)

func DumpRequest(c *fiber.Ctx) []string {
	request := make([]string, 0)
	if c == nil {
		return request
	}

	bytes := c.Fasthttp.Request.String()
	params := strings.Split(bytes, "\r\n")
	for _, param := range params {
		if strings.HasPrefix(param, "Authorization:") { // Authorization header
			request = append(request, "Authorization: *")
		} else if param != "" { // other param
			request = append(request, param)
		}
	}
	return request
}
