package xgin

import (
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
