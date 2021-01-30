package xserverchan

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/go-serverchan"
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

// loggerParam stores some send-event logger parameters, used in LogSendToLogrus and LogSendToLogger.
type loggerParam struct {
	sckey string
	title string
}

// getLoggerParam returns loggerParam from given sckey and title.
func getLoggerParam(sckey, title string) *loggerParam {
	sckeyLen := len(sckey)
	indices := append(xslice.Range(5, sckeyLen/2-4, 1), xslice.Range(sckeyLen/2+3, sckeyLen-6, 1)...) // xxxxx...xxxxxx...xxxxx
	masked := xstring.MaskToken(sckey, '*', indices...)

	return &loggerParam{
		sckey: masked,
		title: xstring.DefaultMaskToken(title),
	}
}

// LogSendToLogrus logs a send-event message to logrus.Logger using given sckey and title.
func LogSendToLogrus(logger *logrus.Logger, sckey, title string, err error, options ...logop.LoggerOption) {
	param := getLoggerParam(sckey, title)
	extra := logop.NewLoggerExtra(options...)

	if err != nil {
		var msg string
		if err == serverchan.ErrDuplicateMessage {
			msg = fmt.Sprintf("[Serverchan] Send duplicate message to %s failed", param.sckey)
		} else {
			msg = fmt.Sprintf("[Serverchan] Send to %s failed: %v", param.sckey, err)
		}
		logger.Error(msg)
		return
	}

	fields := logrus.Fields{
		"module": "serverchan",
		"title":  param.title,
		"sckey":  param.sckey,
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	entry.Info(msg)
}

// LogSendToLogger logs a send-event message to logrus.StdLogger using given sckey and title.
func LogSendToLogger(logger logrus.StdLogger, sckey, title string, err error, options ...logop.LoggerOption) {
	param := getLoggerParam(sckey, title)
	extra := logop.NewLoggerExtra(options...)

	if err != nil {
		var msg string
		if err == serverchan.ErrDuplicateMessage {
			msg = fmt.Sprintf("[Serverchan] Send duplicate message to %s failed", param.sckey)
		} else {
			msg = fmt.Sprintf("[Serverchan] Send to %s failed: %v", param.sckey, err)
		}
		logger.Println(msg)
		return
	}

	msg := formatLogger(param)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}

// formatLogger formats loggerParam to string for logger.
// Logs like:
// 	[Serverchan] SCU12*******************89abcd*******************56789 | te******le
func formatLogger(param *loggerParam) string {
	return fmt.Sprintf("[Serverchan] %s | %s", param.sckey, param.title)
}
