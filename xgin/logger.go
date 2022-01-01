package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// LoggerOption represents an option type for some logger functions' option, can be created by WithXXX functions.
type LoggerOption = internal.LoggerOption

// WithExtraText creates a LoggerOption to specific extra text logging in "... | extra_text" style, notes that if you use this multiple times, only the last one will be retained.
func WithExtraText(text string) LoggerOption {
	return internal.WithExtraText(text)
}

// WithExtraFields creates a LoggerOption to specific logging with extra fields, notes that if you use this multiple times, only the last one will be retained.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return internal.WithExtraFields(fields)
}

// WithExtraFieldsV creates a LoggerOption to specific logging with extra fields in variadic, notes that if you use this multiple times, only the last one will be retained.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return internal.WithExtraFieldsV(fields...)
}

// responseLoggerParam stores some logger parameters used by LogToLogrus and LogToLogger.
type responseLoggerParam struct {
	method   string
	path     string
	status   int
	latency  string
	length   string
	clientIP string
	errorMsg string
}

// extractResponseLoggerData extracts and returns responseLoggerParam and logrus.Fields using given parameters.
func extractResponseLoggerData(c *gin.Context, start, end time.Time) (*responseLoggerParam, logrus.Fields) {
	path := c.Request.URL.Path
	if raw := c.Request.URL.RawQuery; raw != "" {
		path = path + "?" + raw
	}
	latency := end.Sub(start)
	length := c.Writer.Size()
	if length < 0 {
		length = 0
	}
	errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

	param := &responseLoggerParam{
		method:   c.Request.Method,
		path:     path,
		status:   c.Writer.Status(),
		latency:  latency.String(),
		length:   xnumber.RenderByte(float64(length)),
		clientIP: c.ClientIP(),
		errorMsg: strings.TrimSpace(errorMsg),
	}
	fields := logrus.Fields{
		"module":     "gin",
		"method":     param.method,
		"path":       param.path,
		"status":     param.status,
		"start_time": start.Format(time.RFC3339),
		"end_time":   end.Format(time.RFC3339),
		"latency":    latency,
		"length":     length,
		"client_ip":  param.clientIP,
		"error_msg":  param.errorMsg,
	}
	return param, fields
}

// formatResponseLogger formats given responseLoggerParam to string for LogToLogrus and LogToLogger.
//
// Logs like:
// 	[Gin]      200 |      993.3µs |             ::1 |        11B | GET     /test
// 	     |--------| |------------| |---------------| |----------| |-------|-----|
// 	         8            12               15             10          7     ...
func formatResponseLogger(p *responseLoggerParam) string {
	msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s", p.status, p.latency, p.clientIP, p.length, p.method, p.path)
	if p.errorMsg != "" {
		msg += fmt.Sprintf(" | err: %s", p.errorMsg)
	}
	return msg
}

// LogToLogrus logs gin's request and response information to logrus.Logger using given gin.Context and other arguments.
func LogToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...LoggerOption) {
	if logger == nil || c == nil {
		return
	}
	p, f := extractResponseLoggerData(c, start, end)
	m := formatResponseLogger(p)

	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	switch {
	case p.status >= 500:
		logger.WithFields(f).Error(m)
	case p.status >= 400:
		logger.WithFields(f).Warn(m)
	default:
		logger.WithFields(f).Info(m)
	}
}

// LogToLogger logs gin's request and response information to logrus.StdLogger using given gin.Context and other arguments.
func LogToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...LoggerOption) {
	if logger == nil || c == nil {
		return
	}
	p, _ := extractResponseLoggerData(c, start, end)
	m := formatResponseLogger(p)

	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// recoveryLoggerParam stores some logger parameters used by LogRecoveryToLogrus and LogRecoveryToLogger.
type recoveryLoggerParam struct {
	panicMsg  string
	filename  string
	lineIndex int
}

// extractRecoveryLoggerData extracts and returns responseLoggerParam and logrus.Fields using given parameters.
func extractRecoveryLoggerData(v interface{}, stack xruntime.TraceStack) (*recoveryLoggerParam, logrus.Fields) {
	param := &recoveryLoggerParam{panicMsg: fmt.Sprintf("%v", v)}
	if len(stack) > 0 {
		param.filename = stack[0].Filename
		param.lineIndex = stack[0].LineIndex
	}
	fields := logrus.Fields{
		"module":      "recovery",
		"panic_msg":   param.panicMsg,
		"trace_stack": stack.String(),
	}
	return param, fields
}

// formatResponseLogger formats given recoveryLoggerParam to string for LogRecoveryToLogrus and LogRecoveryToLogger.
//
// Logs like:
// 	[Recovery] panic recovered: test error | xxx.go:12
// 	                           |----------| |---------|
// 	                                ...         ...
func formatRecoveryLogger(p *recoveryLoggerParam) string {
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", p.panicMsg)
	if p.filename != "" {
		msg += fmt.Sprintf(" | %s:%d", p.filename, p.lineIndex)
	}
	return msg
}

// LogRecoveryToLogrus logs panic message to logrus.Logger using given value returned from recover and nil-able xruntime.TraceStack.
func LogRecoveryToLogrus(logger *logrus.Logger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption) {
	if logger == nil || v == nil {
		return
	}
	p, f := extractRecoveryLoggerData(v, stack)
	m := formatRecoveryLogger(p)

	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	logger.WithFields(f).Error(m)
}

// LogRecoveryToLogger logs panic message to logrus.StdLogger using given value returned from recover and nil-able xruntime.TraceStack.
func LogRecoveryToLogger(logger logrus.StdLogger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption) {
	if logger == nil || v == nil {
		return
	}
	p, _ := extractRecoveryLoggerData(v, stack)
	m := formatRecoveryLogger(p)

	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}
