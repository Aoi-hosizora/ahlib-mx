package xgin

import (
	"github.com/Aoi-hosizora/ahlib-web/xdto"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"strings"
)

func DumpRequest(c *gin.Context) []string {
	request := make([]string, 0)
	if c == nil {
		return request
	}

	bytes, _ := httputil.DumpRequest(c.Request, false)
	params := strings.Split(string(bytes), "\r\n")
	for _, param := range params {
		if strings.HasPrefix(param, "Authorization:") { // Authorization header
			request = append(request, "Authorization: *")
		} else if param != "" { // other param
			request = append(request, param)
		}
	}
	return request
}

// noinspection GoUnusedExportedFunction
func BuildBasicErrorDto(err interface{}, c *gin.Context) *xdto.ErrorDto {
	return xdto.BuildBasicErrorDto(err, DumpRequest(c))
}

// noinspection GoUnusedExportedFunction
func BuildErrorDto(err interface{}, c *gin.Context, skip int, print bool) *xdto.ErrorDto {
	return xdto.BuildErrorDto(err, DumpRequest(c), skip, print)
}
