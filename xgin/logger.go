package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

func LogrusForGin(logger *logrus.Logger, c *gin.Context) {
	start := time.Now()
	c.Next()
	stop := time.Now()
	latency := stop.Sub(start).String()

	method := c.Request.Method
	path := c.Request.URL.Path
	ip := c.ClientIP()
	code := c.Writer.Status()
	length := math.Abs(float64(c.Writer.Size()))

	entry := logger.WithFields(logrus.Fields{
		"module":   "gin",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	})
	if len(c.Errors) == 0 {
		lengthStr := xnumber.RenderByte(length)
		msg := fmt.Sprintf("[Gin] %3d | %12s | %15s | %8s | %-7s %s", code, latency, ip, lengthStr, method, path)
		if code >= 500 {
			entry.Error(msg)
		} else if code >= 400 {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	} else {
		entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
	}
}
