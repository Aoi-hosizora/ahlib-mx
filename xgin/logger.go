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

// LoggerParam stores the logger parameters, is used in LogToLogrus and LogToLogger.
type LoggerParam struct {
	// request
	Method   string
	Path     string
	Start    time.Time
	ClientIP string

	// response
	Status     int
	Stop       time.Time
	Latency    time.Duration
	LatencyStr string
	Length     int
	LengthStr  string

	// error
	ErrorMessage string
}

// getLoggerParam returns LoggerParam from given gin.Context and time.
func getLoggerParam(c *gin.Context, start, stop time.Time) *LoggerParam {
	path := c.Request.URL.Path
	if raw := c.Request.URL.RawQuery; raw != "" {
		path = path + "?" + raw
	}
	latency := stop.Sub(start)
	length := c.Writer.Size()
	if length < 0 {
		length = 0
	}

	return &LoggerParam{
		Method:   c.Request.Method,
		Path:     path,
		Start:    start,
		ClientIP: c.ClientIP(),

		Status:     c.Writer.Status(),
		Stop:       stop,
		Latency:    latency,
		LatencyStr: latency.String(),
		Length:     length,
		LengthStr:  xnumber.RenderByte(float64(length)),

		ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
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
		"module":    "gin",
		"method":    param.Method,
		"path":      param.Path,
		"status":    param.Status,
		"start":     param.Start,
		"stop":      param.Stop,
		"latency":   param.Latency,
		"length":    param.Length,
		"client_ip": param.ClientIP,
	}
	extra.AddToFields(fields)

	// logger
	entry := logger.WithFields(fields)
	if len(c.Errors) != 0 {
		msg := fmt.Sprintf("[Gin] %s", param.ErrorMessage)
		entry.Error(msg)
	} else {
		msg := formatGinLogger(param)
		extra.AddToMessage(&msg)
		switch {
		case param.Status >= 500:
			entry.Error(msg)
		case param.Status >= 400:
			entry.Warn(msg)
		default:
			entry.Info(msg)
		}
	}
}

// LogToLogrus logs gin's request and response information to logrus.StdLogger using given gin.Context and logop.LoggerOption-s.
func LogToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption) {
	// information
	param := getLoggerParam(c, start, end)

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// logger
	if len(c.Errors) != 0 {
		msg := fmt.Sprintf("[Gin] %s", param.ErrorMessage)
		logger.Print(msg)
	} else {
		msg := formatGinLogger(param)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	}
}

// FormatGinLoggerFunc is a gin logger format function for LogToLogrus and LogToLogger using given LoggerParam.
var FormatGinLoggerFunc func(param *LoggerParam) string

// formatGinLogger represents the inner gin logger format function for LogToLogrus and LogToLogger using given LoggerParam.
// Logs like:
// 	[Gin]      200 |      993.3µs |             ::1 |        11B | GET     /test
// 	     |--------| |------------| |---------------| |----------| |-------|-----|
// 	         8            12               15             10          7     ...
func formatGinLogger(param *LoggerParam) string {
	if FormatGinLoggerFunc != nil {
		return FormatGinLoggerFunc(param)
	}

	return fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s",
		param.Status, param.LatencyStr, param.ClientIP, param.LengthStr, param.Method, param.Path)
}
