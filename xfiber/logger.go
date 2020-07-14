package xfiber

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"time"
)

// Log request and response for fiber, no need for `c.Next()`
func LogrusForFiber(logger *logrus.Logger, start time.Time, c *fiber.Ctx) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())

	entry := logger.WithFields(logrus.Fields{
		"module":   "fiber",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	})
	msg := fmt.Sprintf("[Fiber] %3d | %12s | %15s | %6dB | %-7s %s", code, latency.String(), ip, length, method, path)
	if code >= 500 {
		entry.Error(msg)
	} else if code >= 400 {
		entry.Warn(msg)
	} else {
		entry.Info(msg)
	}
}
