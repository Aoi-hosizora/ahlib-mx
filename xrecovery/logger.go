package xrecovery

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
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

// LogToLogrus logs a panic message to logrus.Logger.
func LogToLogrus(logger *logrus.Logger, err interface{}, options ...logop.LoggerOption) {
	// information
	var errorMessage string
	if e, ok := err.(error); ok {
		errorMessage = e.Error()
	} else {
		errorMessage = fmt.Sprintf("%v", err)
	}

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// fields
	fields := logrus.Fields{
		"module": "recovery",
		"error":  errorMessage,
	}
	extra.AddToFields(fields)

	// logger
	entry := logger.WithFields(fields)
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errorMessage)
	extra.AddToMessage(&msg)
	// [Recovery] panic recovered: test error
	//                            |----------|
	entry.Error(msg)
}

// LogToLogger logs a panic message to logrus.StdLogger such as log.Logger.
func LogToLogger(logger logrus.StdLogger, err interface{}, options ...logop.LoggerOption) {
	// information
	var errorMessage string
	if e, ok := err.(error); ok {
		errorMessage = e.Error()
	} else {
		errorMessage = fmt.Sprintf("%v", err)
	}

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// logger
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errorMessage)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}
