package xrecovery

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/sirupsen/logrus"
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
	errorMessage string
	traceStack   xruntime.TraceStack
}

// getLoggerParam returns loggerParam using given error and xruntime.TraceStack.
func getLoggerParam(err interface{}, stack xruntime.TraceStack) *loggerParam {
	return &loggerParam{
		errorMessage: fmt.Sprintf("%v", err),
		traceStack:   stack,
	}
}

// LogToLogrus logs a panic message to logrus.Logger from given error, nil-able xruntime.TraceStack.
func LogToLogrus(logger *logrus.Logger, err interface{}, stack xruntime.TraceStack, options ...logop.LoggerOption) {
	param := getLoggerParam(err, stack)
	extra := logop.NewLoggerOptions(options)

	fields := logrus.Fields{
		"module":        "recovery",
		"error_message": param.errorMessage,
		"trace_stack":   param.traceStack.String(),
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	entry.Error(msg)
}

// LogToLogger logs a panic message to logrus.StdLogger using given error, nil-able xruntime.TraceStack.
func LogToLogger(logger logrus.StdLogger, err interface{}, stack xruntime.TraceStack, options ...logop.LoggerOption) {
	param := getLoggerParam(err, stack)
	extra := logop.NewLoggerOptions(options)

	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	logger.Print(msg)
}

// formatLogger formats loggerParam to logger string.
// Logs like:
// 	[Recovery] panic recovered: test error | xxx.go:12
// 	                           |----------| |---------|
func formatLogger(param *loggerParam) string {
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", param.errorMessage)
	if len(param.traceStack) > 0 {
		s := param.traceStack[0]
		msg += fmt.Sprintf(" | %s:%d", s.Filename, s.LineIndex)
	}
	return msg
}
