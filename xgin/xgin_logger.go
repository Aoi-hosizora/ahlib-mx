package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// LoggerOption represents an option type for some logger functions' option, can be created by WithXXX functions.
type LoggerOption = internal.LoggerOption

// WithExtraText creates a LoggerOption to specific extra text logging in "...extra_text" style. Note that if you use this multiple times, only the last one will be retained.
func WithExtraText(text string) LoggerOption {
	return internal.WithExtraText(text)
}

// WithExtraFields creates a LoggerOption to specific logging with extra fields. Note that if you use this multiple times, only the last one will be retained.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return internal.WithExtraFields(fields)
}

// WithExtraFieldsV creates a LoggerOption to specific logging with extra fields in variadic. Note that if you use this multiple times, only the last one will be retained.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return internal.WithExtraFieldsV(fields...)
}

// ========
// response
// ========

// ResponseLoggerParam stores some logger parameters and is used by LogResponseToLogrus and LogResponseToLogger.
type ResponseLoggerParam struct {
	// origin
	Context *gin.Context
	Errors  []*gin.Error

	// field
	Method    string
	Path      string
	Query     string
	Status    int
	StartTime time.Time
	EndTime   time.Time
	Latency   time.Duration
	Length    int64
	ClientIP  string
	ErrorMsg  string
}

var (
	// FormatResponseFunc is a custom ResponseLoggerParam's format function for LogResponseToLogrus and LogResponseToLogger.
	FormatResponseFunc func(p *ResponseLoggerParam) string

	// FieldifyResponseFunc is a custom ResponseLoggerParam's fieldify function for LogResponseToLogrus.
	FieldifyResponseFunc func(p *ResponseLoggerParam) logrus.Fields
)

// extractResponseLoggerParam extracts and returns ResponseLoggerParam using given parameters.
func extractResponseLoggerParam(c *gin.Context, start, end time.Time) *ResponseLoggerParam {
	errors := c.Errors.ByType(gin.ErrorTypePublic | gin.ErrorTypePrivate)
	return &ResponseLoggerParam{
		Context:   c,
		Errors:    errors,
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Query:     c.Request.URL.RawQuery,
		Status:    c.Writer.Status(),
		StartTime: start,
		EndTime:   end,
		Latency:   end.Sub(start),
		Length:    int64(c.Writer.Size()),
		ClientIP:  c.ClientIP(),
		ErrorMsg:  errors.String(),
	}
}

// formatResponseLoggerParam formats given ResponseLoggerParam to string for LogResponseToLogrus and LogResponseToLogger.
//
// The default format logs like:
// 	[Gin]      200 |      993.3µs |        11B |             ::1 | GET     /test
// 	[Gin]      500 |       2.34ms |       506B |       127.0.0.1 | POST    /test
// 	[Gin]      404 |        5.67s |         0B |   192.168.1.102 | PUT     /erro | err: test error
// 	     |--------| |------------| |----------| |---------------| |-------|-----|
// 	         8            12            10              15            7     ...
func formatResponseLoggerParam(p *ResponseLoggerParam) string {
	if FormatResponseFunc != nil {
		return FormatResponseFunc(p)
	}
	msg := fmt.Sprintf("[Gin] %8d | %12s | %10s | %15s | %-7s %s", p.Status, p.Latency.String(), xnumber.FormatByteSize(float64(p.Length)), p.ClientIP, p.Method, p.Path)
	if p.ErrorMsg != "" {
		msg += fmt.Sprintf(" | err: %s", p.ErrorMsg)
	}
	return msg
}

// fieldifyResponseLoggerParam fieldifies given ResponseLoggerParam to logrus.Fields for LogResponseToLogrus.
//
// The default contains the following fields: module, method, path, query, status, start_time, end_time, latency, length, client_ip, error_msg.
func fieldifyResponseLoggerParam(p *ResponseLoggerParam) logrus.Fields {
	if FieldifyResponseFunc != nil {
		return FieldifyResponseFunc(p)
	}
	f := logrus.Fields{
		"module":     "gin",
		"method":     p.Method,
		"path":       p.Path,
		"query":      p.Query,
		"status":     p.Status,
		"start_time": p.StartTime,
		"end_time":   p.EndTime,
		"latency":    p.Latency,
		"length":     p.Length,
		"client_ip":  p.ClientIP,
	}
	if p.ErrorMsg != "" {
		f["error_msg"] = p.ErrorMsg
	}
	return f
}

// LogResponseToLogrus logs gin's request and response information to logrus.Logger using given gin.Context and other arguments.
func LogResponseToLogrus(logger *logrus.Logger, c *gin.Context, start, end time.Time, options ...LoggerOption) {
	if logger == nil || c == nil {
		return
	}
	p := extractResponseLoggerParam(c, start, end)
	m := formatResponseLoggerParam(p)
	f := fieldifyResponseLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	switch {
	case p.Status >= 500:
		logger.WithFields(f).Error(m)
	case p.Status >= 400:
		logger.WithFields(f).Warn(m)
	default:
		logger.WithFields(f).Info(m)
	}
}

// LogResponseToLogger logs gin's request and response information to logrus.StdLogger using given gin.Context and other arguments.
func LogResponseToLogger(logger logrus.StdLogger, c *gin.Context, start, end time.Time, options ...LoggerOption) {
	if logger == nil || c == nil {
		return
	}
	p := extractResponseLoggerParam(c, start, end)
	m := formatResponseLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// ========
// recovery
// ========

// RecoveryLoggerParam stores some logger parameters used by LogRecoveryToLogrus and LogRecoveryToLogger.
type RecoveryLoggerParam struct {
	// origin
	PanicValue interface{}
	TraceStack xruntime.TraceStack

	// field
	PanicMsg     string
	FullFilename string
	FullFuncname string
	LineIndex    int
	LineText     string
}

var (
	// FormatRecoveryFunc is a custom RecoveryLoggerParam's format function for LogRecoveryToLogrus and LogRecoveryToLogger.
	FormatRecoveryFunc func(p *RecoveryLoggerParam) string

	// FieldifyRecoveryFunc is a custom RecoveryLoggerParam's fieldify function for LogRecoveryToLogrus.
	FieldifyRecoveryFunc func(p *RecoveryLoggerParam) logrus.Fields
)

// extractRecoveryLoggerParam extracts and returns ResponseLoggerParam using given parameters.
func extractRecoveryLoggerParam(v interface{}, stack xruntime.TraceStack) *RecoveryLoggerParam {
	msg := fmt.Sprintf("%v", v)
	filename, funcname, lineIndex, lineText := "", "", 0, ""
	if len(stack) > 0 {
		filename = stack[0].Filename
		funcname = stack[0].FuncFullName
		lineIndex = stack[0].LineIndex
		lineText = stack[0].LineText
	}
	return &RecoveryLoggerParam{
		PanicValue:   v,
		TraceStack:   stack,
		PanicMsg:     msg,
		FullFilename: filename,
		FullFuncname: funcname,
		LineIndex:    lineIndex,
		LineText:     lineText,
	}
}

// formatResponseLoggerParam formats given RecoveryLoggerParam to string for LogRecoveryToLogrus and LogRecoveryToLogger.
//
// The default format logs like:
// 	[Recovery] panic recovered: test error | xxx.go:12
// 	                           |----------| |---------|
// 	                                ...         ...
func formatRecoveryLoggerParam(p *RecoveryLoggerParam) string {
	if FormatRecoveryFunc != nil {
		return FormatRecoveryFunc(p)
	}
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", p.PanicMsg)
	if p.FullFilename != "" {
		msg += fmt.Sprintf(" | %s:%d", p.FullFilename, p.LineIndex)
	}
	return msg
}

// fieldifyRecoveryLoggerParam fieldifies given RecoveryLoggerParam to logrus.Fields for LogRecoveryToLogrus.
//
// The default contains the following fields: module, panic_msg, trace_stack, filename, funcname, line_index, line_text.
func fieldifyRecoveryLoggerParam(p *RecoveryLoggerParam) logrus.Fields {
	if FieldifyRecoveryFunc != nil {
		return FieldifyRecoveryFunc(p)
	}
	return logrus.Fields{
		"module":      "recovery",
		"panic_msg":   p.PanicMsg,
		"trace_stack": p.TraceStack.String(),
		"filename":    p.FullFilename,
		"funcname":    p.FullFuncname,
		"line_index":  p.LineIndex,
		"line_text":   p.LineText,
	}
}

// LogRecoveryToLogrus logs panic message to logrus.Logger using given recover returned value and nil-able xruntime.TraceStack.
func LogRecoveryToLogrus(logger *logrus.Logger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption) {
	if logger == nil || v == nil {
		return
	}
	p := extractRecoveryLoggerParam(v, stack)
	m := formatRecoveryLoggerParam(p)
	f := fieldifyRecoveryLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	logger.WithFields(f).Error(m)
}

// LogRecoveryToLogger logs panic message to logrus.StdLogger using given recover returned value and nil-able xruntime.TraceStack.
func LogRecoveryToLogger(logger logrus.StdLogger, v interface{}, stack xruntime.TraceStack, options ...LoggerOption) {
	if logger == nil || v == nil {
		return
	}
	p := extractRecoveryLoggerParam(v, stack)
	m := formatRecoveryLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}
