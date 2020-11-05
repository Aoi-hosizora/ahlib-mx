package xfiber

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/xmap"
	"github.com/Aoi-hosizora/ahlib-web/internal/xwlogger"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

// WithExtraString represents LoggerOption for logging extra string.
func WithExtraString(s string) xwlogger.LoggerOption {
	return func(ex *xwlogger.LoggerExtra) {
		ex.ExtraString = &s
	}
}

// WithExtraFields represents LoggerOption for logging extra fields.
func WithExtraFields(m map[string]interface{}) xwlogger.LoggerOption {
	return func(ex *xwlogger.LoggerExtra) {
		ex.ExtraFields = &m
	}
}

// WithExtraFieldsV represents LoggerOption for logging extra fields (vararg).
func WithExtraFieldsV(m ...interface{}) xwlogger.LoggerOption {
	return func(ex *xwlogger.LoggerExtra) {
		m := xmap.SliceToStringMap(m)
		ex.ExtraFields = &m
	}
}

// WithLogrus logs fiber's request and error with logrus.Logger.
func WithLogrus(logger *logrus.Logger, start time.Time, c *fiber.Ctx, options ...xwlogger.LoggerOption) {
	// information
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))
	latencyStr := latency.String()
	ip := c.IP()

	// extra
	extra := &xwlogger.LoggerExtra{}
	extra.ApplyOptions(options)

	// field
	fields := logrus.Fields{
		"module":   "fiber",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	// logger
	if c.Error() != nil {
		msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		extra.AddToString(&msg)
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

// WithLogger logs fiber's request and error with log.Logger.
func WithLogger(logger *log.Logger, start time.Time, c *fiber.Ctx, options ...xwlogger.LoggerOption) {
	// information
	latency := time.Now().Sub(start)
	method := c.Method()
	path := c.Path()
	code := c.Fasthttp.Response.StatusCode()
	length := len(c.Fasthttp.Response.Body())
	lengthStr := xnumber.RenderByte(float64(length))
	latencyStr := latency.String()
	ip := c.IP()

	// extra
	extra := &xwlogger.LoggerExtra{}
	extra.ApplyOptions(options)

	// logger
	if c.Error() != nil {
		msg := fmt.Sprintf("[Fiber] %6d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		extra.AddToString(&msg)
		logger.Println(msg)
	} else {
		msg := fmt.Sprintf("[Fiber] %s", c.Error().Error())
		logger.Println(msg)
	}
}
