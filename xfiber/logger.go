package xfiber

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"time"
)

func LogrusForFiber(logger *logrus.Logger, c *fiber.Ctx) {
	start := time.Now()
	c.Next()
	stop := time.Now()
	latency := stop.Sub(start).String()

	method := c.Method()
	path := c.Path()
	ip := c.IP()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())

	entry := logger.WithFields(logrus.Fields{
		"module":   "gin",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	})
	msg := fmt.Sprintf("[Gin] %3d | %12s | %15s | %6dB | %-7s %s", code, latency, ip, length, method, path)
	if code >= 500 {
		entry.Error(msg)
	} else if code >= 400 {
		entry.Warn(msg)
	} else {
		entry.Info(msg)
	}
}
