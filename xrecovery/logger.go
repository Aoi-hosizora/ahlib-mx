package xrecovery

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/sirupsen/logrus"
	"log"
)

// WithExtraString represents LoggerOption for logging extra string.
func WithExtraString(s string) logop.LoggerOption {
	return logop.WithExtraText(s)
}

// WithExtraFields represents LoggerOption for logging extra fields.
func WithExtraFields(m map[string]interface{}) logop.LoggerOption {
	return logop.WithExtraFields(m)
}

// WithExtraFieldsV represents LoggerOption for logging extra fields (vararg).
func WithExtraFieldsV(m ...interface{}) logop.LoggerOption {
	return logop.WithExtraFieldsV(m...)
}

// WithLogrus logs a panic message with logrus.Logger.
func WithLogrus(logger *logrus.Logger, err interface{}, options ...logop.LoggerOption) {
	// information
	errMessage := ""
	if e, ok := err.(error); ok {
		errMessage = e.Error()
	} else {
		errMessage = fmt.Sprintf("%v", err)
	}

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// fields
	fields := logrus.Fields{
		"module": "panic",
		"error":  errMessage,
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	// logger
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errMessage)
	extra.AddToMessage(&msg)
	entry.Error(msg)
}

// WithLogger logs a panic message with log.Logger.
func WithLogger(logger *log.Logger, err interface{}, options ...logop.LoggerOption) {
	// information
	errMessage := ""
	if e, ok := err.(error); ok {
		errMessage = e.Error()
	} else {
		errMessage = fmt.Sprintf("%v", err)
	}

	// extra
	extra := &logop.LoggerExtraOptions{}
	extra.ApplyOptions(options)

	// logger
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errMessage)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}
