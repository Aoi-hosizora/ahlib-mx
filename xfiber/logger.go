package xfiber

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

// LoggerExtra represents the extra strings and fields.
type LoggerExtra struct {
	OtherString string
	OtherFields map[string]interface{}
}

// WithLogrus logs every fiber's requests and errors with logrus.Logger.
func WithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx, extra *LoggerExtra) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))
	latencyStr := latency.String()
	ip := c.IP()

	fields := logrus.Fields{
		"module":   "gin",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	}
	if extra != nil && extra.OtherFields != nil {
		for k, v := range extra.OtherFields {
			fields[k] = v
		}
	}
	entry := logger.WithFields(fields)

	if c.Error() != nil {
		msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		if extra != nil && extra.OtherString != "" {
			msg += fmt.Sprintf(" | %s", extra.OtherString)
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

// WithLogger logs every fiber's requests and errors with log.Logger.
func WithLogger(logger *log.Logger, start time.Time, c *fiber.Ctx, other string) {
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))
	latencyStr := latency.String()
	ip := c.IP()

	if c.Error() != nil {
		msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		if other != "" {
			msg += fmt.Sprintf(" | %s", other)
		}
		logger.Println(msg)
	} else {
		msg := fmt.Sprintf("[Fiber] %s", c.Error().Error())
		logger.Println(msg)
	}
}
