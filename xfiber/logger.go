package xfiber

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"time"
)

// Log request and response for fiber, no need for `c.Next()`.
func LoggerWithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))

	entry := logger.WithFields(logrus.Fields{
		"module":   "fiber",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	})
	msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latency.String(), ip, lengthStr, method, path)
	if code >= 500 {
		entry.Error(msg)
	} else if code >= 400 {
		entry.Warn(msg)
	} else {
		entry.Info(msg)
	}
}
