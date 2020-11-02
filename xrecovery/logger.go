package xrecovery

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/xwlogger"
	"github.com/sirupsen/logrus"
	"log"
)

// WithExtraString represents LoggerOption for logging extra string.
func WithExtraString(s string) xwlogger.LoggerOption {
	return func(ex *xwlogger.LoggerExtra) {
		ex.ExtraString = &s
	}
}

// WithExtraFields represents LoggerOption for logging extra fields.
func WithExtraFields(m map[string]interface{}) xwlogger.LoggerOption {
	return func(ex *xwlogger.LoggerExtra) {
		ex.ExtraFields = &m
	}
}

// WithLogrus logs a panic message with logrus.Logger.
func WithLogrus(logger *logrus.Logger, err interface{}, options ...xwlogger.LoggerOption) {
	// information
	errMessage := ""
	if e, ok := err.(error); ok {
		errMessage = e.Error()
	} else {
		errMessage = fmt.Sprintf("%v", err)
	}

	// extra
	extra := &xwlogger.LoggerExtra{}
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
	extra.AddToString(&msg)
	entry.Error(msg)
}

// WithLogger logs a panic message with log.Logger.
func WithLogger(logger *log.Logger, err interface{}, options ...xwlogger.LoggerOption) {
	// information
	errMessage := ""
	if e, ok := err.(error); ok {
		errMessage = e.Error()
	} else {
		errMessage = fmt.Sprintf("%v", err)
	}

	// extra
	extra := &xwlogger.LoggerExtra{}
	extra.ApplyOptions(options)

	// logger
	msg := fmt.Sprintf("[Recovery] panic recovered: %s", errMessage)
	extra.AddToString(&msg)
	logger.Println(msg)
}
