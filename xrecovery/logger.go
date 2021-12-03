package xrecovery

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logopt"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/sirupsen/logrus"
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
	errorMessage string
	traceStack   xruntime.TraceStack
}

// getLoggerParamAndFields returns loggerParam and logrus.Fields using given parameters.
func getLoggerParamAndFields(err interface{}, stack xruntime.TraceStack) (*loggerParam, logrus.Fields) {
	param := &loggerParam{
		errorMessage: fmt.Sprintf("%v", err),
		traceStack:   stack,
	}
	fields := logrus.Fields{
		"module":        "recovery",
		"error_message": param.errorMessage,
		"trace_stack":   param.traceStack.String(),
	}
	return param, fields
}

// LogToLogrus logs a panic message to logrus.Logger from given error, nil-able xruntime.TraceStack.
func LogToLogrus(logger *logrus.Logger, err interface{}, stack xruntime.TraceStack, options ...logopt.LoggerOption) {
	p, f := getLoggerParamAndFields(err, stack)
	m := formatLogger(p)

	extra := logopt.NewLoggerOptions(options)
	extra.AddToMessage(&m)
	extra.AddToFields(f)
	logger.WithFields(f).Error(m)
}

// LogToLogger logs a panic message to logrus.StdLogger using given error, nil-able xruntime.TraceStack.
func LogToLogger(logger logrus.StdLogger, err interface{}, stack xruntime.TraceStack, options ...logopt.LoggerOption) {
	p, _ := getLoggerParamAndFields(err, stack)
	m := formatLogger(p)

	extra := logopt.NewLoggerOptions(options)
	extra.AddToMessage(&m)
	logger.Print(m)
}

// formatLogger formats loggerParam to logger string.
// Logs like:
// 	[Recovery] panic recovered: test error | xxx.go:12
// 	                           |----------| |---------|
// 	                                ...         ...
func formatLogger(param *loggerParam) string {
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", param.errorMessage)
	if len(param.traceStack) > 0 {
		s := param.traceStack[0]
		msg += fmt.Sprintf(" | %s:%d", s.Filename, s.LineIndex)
	}
	return msg
}
