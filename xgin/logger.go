package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/xwlogger"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"math"
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

// WithLogrus logs gin's request and error with logrus.Logger.
func WithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context, options ...xwlogger.LoggerOption) {
	// information
	latency := time.Now().Sub(start)
	method := c.Request.Method
	path := c.Request.URL.Path
	code := c.Writer.Status()
	length := math.Abs(float64(c.Writer.Size()))
	lengthStr := xnumber.RenderByte(length)
	latencyStr := latency.String()
	ip := c.ClientIP()

	// extra
	extra := &xwlogger.LoggerExtra{}
	extra.ApplyOptions(options)

	// fields
	fields := logrus.Fields{
		"module":   "gin",
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
	if len(c.Errors) == 0 {
		msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		extra.AddToString(&msg)
		if code >= 500 {
			entry.Error(msg)
		} else if code >= 400 {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	} else {
		msg := fmt.Sprintf("[Gin] %s", c.Errors.ByType(gin.ErrorTypePrivate).String())
		entry.Error(msg)
	}
}

// WithLogger logs gin's request and error with log.Logger.
func WithLogger(logger *log.Logger, start time.Time, c *gin.Context, options ...xwlogger.LoggerOption) {
	// information
	latency := time.Now().Sub(start)
	method := c.Request.Method
	path := c.Request.URL.Path
	code := c.Writer.Status()
	length := math.Abs(float64(c.Writer.Size()))
	lengthStr := xnumber.RenderByte(length)
	latencyStr := latency.String()
	ip := c.ClientIP()

	// extra
	extra := &xwlogger.LoggerExtra{}
	extra.ApplyOptions(options)

	// logger
	if len(c.Errors) == 0 {
		msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		extra.AddToString(&msg)
		logger.Println(msg)
	} else {
		msg := c.Errors.ByType(gin.ErrorTypePrivate).String()
		logger.Println(fmt.Sprintf("[Gin] %s", msg))
	}
}
