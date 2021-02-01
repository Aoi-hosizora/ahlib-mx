package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// WithExtraText creates a logger option to log with extra text.
func WithExtraText(text string) logop.LoggerOption {
	return logop.WithExtraText(text)
}

// WithExtraFields creates a logger option to log with extra fields.
func WithExtraFields(fields map[string]interface{}) logop.LoggerOption {
	return logop.WithExtraFields(fields)
}

// WithExtraFieldsV creates a logger option to log with extra fields in vararg.
func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption {
	return logop.WithExtraFieldsV(fields...)
}

// loggerParam stores some logger parameters, used in LogToLogrus and LogToLogger.
type loggerParam struct {
	method       string
	path         string
	status       int
	startTime    time.Time
	endTime      time.Time
	latency      time.Duration
	length       int
	clientIP     string
	contextError string
}

// getLoggerParamAndFields returns loggerParam and logrus.Fields from given gin.Context and times.
func getLoggerParamAndFields(c *gin.Context, start, end time.Time) (*loggerParam, logrus.Fields) {
	path := c.Request.URL.Path
	if raw := c.Request.URL.RawQuery; raw != "" {
		path = path + "?" + raw
	}
	length := c.Writer.Size()
	if length < 0 {
		length = 0
	}
	errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
	errorMessage = strings.TrimSpace(errorMessage)

	param := &loggerParam{
		method:       c.Request.Method,
		path:         path,
		status:       c.Writer.Status(),
		startTime:    start,
		endTime:      end,
		latency:      end.Sub(start),
		length:       length,
		clientIP:     c.ClientIP(),
		contextError: errorMessage,
	}
	fields := logrus.Fields{
		"module":     "gin",
		"method":     param.method,
		"path":       param.path,
		"status":     param.status,
		"start_time": param.startTime.Format(time.RFC3339),
		"end_time":   param.endTime.Format(time.RFC3339),
		"latency":    param.latency,
		"length":     param.length,
		"client_ip":  param.clientIP,
		"ctx_error":  param.contextError,
	}
	return param, fields
}

// LogToLogrus logs gin's request and response information to logrus.Logger using given gin.Context and times.
func LogToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption) {
	param, fields := getLoggerParamAndFields(c, start, end)
	extra := logop.NewLoggerOptions(options)
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	switch {
	case param.status >= 500:
		entry.Error(msg)
	case param.status >= 400:
		entry.Warn(msg)
	default:
		entry.Info(msg)
	}
}

// LogToLogrus logs gin's request and response information to logrus.StdLogger using given gin.Context and times.
func LogToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...logop.LoggerOption) {
	param, _ := getLoggerParamAndFields(c, start, end)
	extra := logop.NewLoggerOptions(options)

	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	logger.Print(msg)
}

// formatLogger formats loggerParam to logger string.
// Logs like:
// 	[Gin]      200 |      993.3µs |             ::1 |        11B | GET     /test
// 	     |--------| |------------| |---------------| |----------| |-------|-----|
// 	         8            12               15             10          7     ...
func formatLogger(param *loggerParam) string {
	msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s",
		param.status, param.latency.String(), param.clientIP, xnumber.RenderByte(float64(param.length)), param.method, param.path)
	if param.contextError != "" {
		msg = fmt.Sprintf("%s | (%s)", msg, param.contextError)
	}
	return msg
}
