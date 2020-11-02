package xgin

import (
	"github.com/Aoi-hosizora/ahlib-web/xdto"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"strings"
)

// DumpRequest returns request strings from gin.Context using http.Request.
func DumpRequest(c *gin.Context) []string {
	request := make([]string, 0)
	if c == nil {
		return request
	}

	bs, _ := httputil.DumpRequest(c.Request, false)
	params := strings.Split(string(bs), "\r\n")
	for _, param := range params {
		if strings.HasPrefix(param, "Authorization:") { // Authorization header
			request = append(request, "Authorization: *")
		} else if param != "" { // other param
			request = append(request, param)
		}
	}
	return request
}

// BuildBasicErrorDto builds xdto.ErrorDto from error and gin.Context.
func BuildBasicErrorDto(err interface{}, c *gin.Context, otherKvs ...interface{}) *xdto.ErrorDto {
	return xdto.BuildBasicErrorDto(err, DumpRequest(c), otherKvs...)
}

// BuildErrorDto builds xdto.ErrorDto from error, gin.Context and runtime.
func BuildErrorDto(err interface{}, c *gin.Context, skip int, doPrint bool, otherKvs ...interface{}) *xdto.ErrorDto {
	skip++
	return xdto.BuildErrorDto(err, DumpRequest(c), skip, doPrint, otherKvs...)
}
