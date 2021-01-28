package xrecovery

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/sirupsen/logrus"
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
	Error        interface{}
	ErrorMessage string
	Stack        xruntime.TraceStack
}

// getLoggerParam returns LoggerParam from given error and xruntime.TraceStack.
func getLoggerParam(err interface{}, stack xruntime.TraceStack) *LoggerParam {
	var errorMessage string
	if e, ok := err.(error); ok {
		errorMessage = e.Error()
	} else {
		errorMessage = fmt.Sprintf("%v", err)
	}

	return &LoggerParam{
		Error:        err,
		ErrorMessage: errorMessage,
		Stack:        stack,
	}
}

// LogToLogrus logs a panic message to logrus.Logger using given error, nil-able xruntime.TraceStack and logop.LoggerOption-s.
func LogToLogrus(logger *logrus.Logger, err interface{}, stack xruntime.TraceStack, options ...logop.LoggerOption) {
	// information
	param := getLoggerParam(err, stack)

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// fields
	fields := logrus.Fields{
		"module":      "recovery",
		"error":       param.ErrorMessage,
		"trace_stack": param.Stack,
	}
	extra.AddToFields(fields)

	// logger
	entry := logger.WithFields(fields)
	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	entry.Error(msg)
}

// LogToLogger logs a panic message to logrus.StdLogger using given error, nil-able xruntime.TraceStack and logop.LoggerOption-s.
func LogToLogger(logger logrus.StdLogger, err interface{}, stack xruntime.TraceStack, options ...logop.LoggerOption) {
	// information
	param := getLoggerParam(err, stack)

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// logger
	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}

// FormatLoggerFunc is a recovery logger format function for LogToLogrus and LogToLogger using given LoggerParam.
var FormatLoggerFunc func(param *LoggerParam) string

// formatLogger represents the inner recovery logger format function for LogToLogrus and LogToLogger using given LoggerParam.
// Logs like:
// 	[Recovery] panic recovered: test error | xxx.go:5
// 	                           |----------| |--------|
func formatLogger(param *LoggerParam) string {
	if FormatLoggerFunc != nil {
		return FormatLoggerFunc(param)
	}

	return fmt.Sprintf("[Recovery] panic recovered: %s", param.ErrorMessage)
}
