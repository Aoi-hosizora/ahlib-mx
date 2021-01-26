package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// WithExtraText creates an option for logger to log with extra text.
func WithExtraText(text string) logop.LoggerOption {
	return logop.WithExtraText(text)
}

// WithExtraText creates an option for logger to log with extra fields.
func WithExtraFields(fields map[string]interface{}) logop.LoggerOption {
	return logop.WithExtraFields(fields)
}

// WithExtraText creates an option for logger to log with extra fields in vararg.
func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption {
	return logop.WithExtraFieldsV(fields...)
}

// loggerParam stores the logger parameters, used in LogToLogrus and LogToLogger.
type loggerParam struct {
	// request
	method   string
	path     string
	start    time.Time
	clientIP string

	// response
	status     int
	stop       time.Time
	latency    time.Duration
	latencyStr string
	length     int
	lengthStr  string

	// error
	errorMessage string
}

// getLoggerParam returns loggerParam from given gin.Context and time.
func getLoggerParam(c *gin.Context, start, stop time.Time) *loggerParam {
	path := c.Request.URL.Path
	if raw := c.Request.URL.RawQuery; raw != "" {
		path = path + "?" + raw
	}
	latency := stop.Sub(start)
	length := c.Writer.Size()
	if length < 0 {
		length = 0
	}

	return &loggerParam{
		method:   c.Request.Method,
		path:     path,
		start:    start,
		clientIP: c.ClientIP(),

		status:     c.Writer.Status(),
		stop:       stop,
		latency:    latency,
		latencyStr: latency.String(),
		length:     length,
		lengthStr:  xnumber.RenderByte(float64(length)),

		errorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
	}
}

// LogToLogrus logs gin's request and response information to logrus.Logger using given gin.Context and logop.LoggerOption-s.
func LogToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption) {
	// information
	param := getLoggerParam(c, start, end)

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// fields
	fields := logrus.Fields{
		"module":   "gin",
		"method":   param.method,
		"path":     param.path,
		"status":   param.status,
		"start":    param.start,
		"stop":     param.stop,
		"latency":  param.latency,
		"length":   param.length,
		"clientIP": param.clientIP,
	}
	extra.AddToFields(fields)

	// logger
	entry := logger.WithFields(fields)
	if len(c.Errors) != 0 {
		msg := fmt.Sprintf("[Gin] %s", param.errorMessage)
		entry.Error(msg)
	} else {
		msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s",
			param.status, param.latencyStr, param.clientIP, param.lengthStr, param.method, param.path)
		extra.AddToMessage(&msg)
		// [Gin]      200 |      993.3µs |             ::1 |        11B | GET     /test
		//      |--------| |------------| |---------------| |----------| |-------|-----|
		//          8            12               15             10          7     ...
		switch {
		case param.status >= 500:
			entry.Error(msg)
		case param.status >= 400:
			entry.Warn(msg)
		default:
			entry.Info(msg)
		}
	}
}

// LogToLogrus logs gin's request and response information to logrus.StdLogger such as log.Logger using given gin.Context and logop.LoggerOption-s.
func LogToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption) {
	// information
	param := getLoggerParam(c, start, end)

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// logger
	if len(c.Errors) != 0 {
		msg := fmt.Sprintf("[Gin] %s", param.errorMessage)
		logger.Print(msg)
	} else {
		msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s",
			param.status, param.latencyStr, param.clientIP, param.lengthStr, param.method, param.path)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	}
}
