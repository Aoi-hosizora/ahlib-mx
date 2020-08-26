package xfiber

import (
	"github.com/Aoi-hosizora/ahlib-web/xdto"
	"github.com/gofiber/fiber"
	"strings"
)

func DumpRequest(c *fiber.Ctx) []string {
	request := make([]string, 0)
	if c == nil {
		return request
	}

	str := c.Fasthttp.Request.String()
	params := strings.Split(str, "\r\n")
	for _, param := range params {
		if strings.HasPrefix(param, "Authorization:") { // Authorization header
			request = append(request, "Authorization: *")
		} else if param != "" { // other param
			request = append(request, param)
		}
	}
	return request
}

func BuildBasicErrorDto(err interface{}, c *fiber.Ctx) *xdto.ErrorDto {
	return xdto.BuildBasicErrorDto(err, DumpRequest(c))
}

func BuildErrorDto(err interface{}, c *fiber.Ctx, skip int, print bool) *xdto.ErrorDto {
	skip++
	return xdto.BuildErrorDto(err, DumpRequest(c), skip, print)
}
