package xfiber

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func WithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx, other string, otherFields map[string]interface{}) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))

	fields := logrus.Fields{
		"module":   "gin",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	}
	if otherFields != nil {
		for k, v := range otherFields {
			fields[k] = v
		}
	}
	entry := logger.WithFields(fields)

	if c.Error() != nil {
		msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latency.String(), ip, lengthStr, method, path)
		if other != "" {
			msg += fmt.Sprintf(" | %s", other)
		}

		if code >= 500 {
			entry.Error(msg)
		} else if code >= 400 {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	} else {
		msg := fmt.Sprintf("[Fiber] %s", c.Error().Error())
		entry.Error(msg)
	}
}

func WithLogger(logger *log.Logger, start time.Time, c *fiber.Ctx, other string) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))

	if c.Error() != nil {
		msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latency.String(), ip, lengthStr, method, path)
		if other != "" {
			msg += fmt.Sprintf(" | %s", other)
		}
		logger.Println(msg)
	} else {
		msg := fmt.Sprintf("[Fiber] %s", c.Error().Error())
		logger.Println(msg)
	}
}
