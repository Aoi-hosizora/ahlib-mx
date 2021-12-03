package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logopt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// WithExtraText creates a logger option to log with extra text.
func WithExtraText(text string) logopt.LoggerOption {
	return logopt.WithExtraText(text)
}

// WithExtraFields creates a logger option to log with extra fields.
func WithExtraFields(fields map[string]interface{}) logopt.LoggerOption {
	return logopt.WithExtraFields(fields)
}

// WithExtraFieldsV creates a logger option to log with extra fields in vararg.
func WithExtraFieldsV(fields ...interface{}) logopt.LoggerOption {
	return logopt.WithExtraFieldsV(fields...)
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

// getLoggerParamAndFields returns loggerParam and logrus.Fields using given parameters.
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
func LogToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...logopt.LoggerOption) {
	p, f := getLoggerParamAndFields(c, start, end)
	m := formatLogger(p)

	extra := logopt.NewLoggerOptions(options)
	extra.AddToMessage(&m)
	extra.AddToFields(f)
	switch {
	case p.status >= 500:
		logger.WithFields(f).Error(m)
	case p.status >= 400:
		logger.WithFields(f).Warn(m)
	default:
		logger.WithFields(f).Info(m)
	}
}

// LogToLogger logs gin's request and response information to logrus.StdLogger using given gin.Context and times.
func LogToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...logopt.LoggerOption) {
	p, _ := getLoggerParamAndFields(c, start, end)
	m := formatLogger(p)

	extra := logopt.NewLoggerOptions(options)
	extra.AddToMessage(&m)
	logger.Print(m)
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
		msg += fmt.Sprintf(" | err: %s", param.contextError)
	}
	return msg
}
