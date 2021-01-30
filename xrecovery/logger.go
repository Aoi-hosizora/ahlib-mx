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

// WithExtraFields creates an option for logger to log with extra fields.
func WithExtraFields(fields map[string]interface{}) logop.LoggerOption {
	return logop.WithExtraFields(fields)
}

// WithExtraFieldsV creates an option for logger to log with extra fields in vararg.
func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption {
	return logop.WithExtraFieldsV(fields...)
}

// LogToLogrus logs a panic message to logrus.Logger using given error, nil-able xruntime.TraceStack.
// Logs like:
// 	[Recovery] panic recovered: test error
// 	                           |----------|
func LogToLogrus(logger *logrus.Logger, err interface{}, stack xruntime.TraceStack, options ...logop.LoggerOption) {
	errMsg := fmt.Sprintf("%v", err)
	extra := logop.NewLoggerExtra(options...)
	fields := logrus.Fields{
		"module": "recovery",
		"error":  errMsg,
		"stack":  stack.String(),
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errMsg)
	extra.AddToMessage(&msg)
	entry.Error(msg)
}

// LogToLogger logs a panic message to logrus.StdLogger using given error, nil-able xruntime.TraceStack.
func LogToLogger(logger logrus.StdLogger, err interface{}, stack xruntime.TraceStack, options ...logop.LoggerOption) {
	errMsg := fmt.Sprintf("%v", err)
	extra := logop.NewLoggerExtra(options...)

	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errMsg)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}
