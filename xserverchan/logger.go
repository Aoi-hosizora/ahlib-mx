package xserverchan

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/go-serverchan"
	"github.com/sirupsen/logrus"
	"regexp"
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
	maskedSckey string
	maskedTitle string
	titleLength int
	bodyLength  int
}

var (
	// _asteriskRe is used to replace "*" from masked title.
	_asteriskRe = regexp.MustCompile(`\*+`)
)

// getLoggerParamAndFields returns loggerParam and logrus.Fields from given sckey, title, body and error.
func getLoggerParamAndFields(sckey, title, body string, err error) (*loggerParam, logrus.Fields) {
	indices := append(xslice.Range(6, 15, 1), append(xslice.Range(22, 31, 1), xslice.Range(38, 47, 1)...)...)
	maskedSckey := xstring.MaskToken(sckey, '*', indices...) // XXXXXX**********XXXXXX**********XXXXXX**********XXXXXX
	maskedTitle := xstring.DefaultMaskToken(title)
	maskedTitle = _asteriskRe.ReplaceAllString(maskedTitle, "**")

	param := &loggerParam{
		maskedSckey: maskedSckey,
		maskedTitle: maskedTitle,
		titleLength: len(title),
		bodyLength:  len(body),
	}
	var fields logrus.Fields
	if err == nil {
		fields = logrus.Fields{
			"module":       "serverchan",
			"masked_sckey": param.maskedSckey,
			"masked_title": param.maskedTitle,
			"title_length": param.titleLength, // <<<
			"body_length":  param.bodyLength,  // <<<
		}
	} else {
		fields = logrus.Fields{
			"module":       "serverchan",
			"masked_sckey": param.maskedSckey,
			"masked_title": param.maskedTitle,
			"error":        err, // <<<
		}
	}
	return param, fields
}

// LogToLogrus logs a send-event message to logrus.Logger using given sckey, title, body and error.
func LogToLogrus(logger *logrus.Logger, sckey, title, body string, err error, options ...logop.LoggerOption) {
	param, fields := getLoggerParamAndFields(sckey, title, body, err)
	extra := logop.NewLoggerOptions(options)
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	if err != nil {
		msg := formatErrorLogger(param, err)
		extra.AddToMessage(&msg)
		entry.Error(msg)
	} else {
		msg := formatLogger(param)
		extra.AddToMessage(&msg)
		entry.Info(msg)
	}
}

// LogToLogger logs a send-event message to logrus.StdLogger using given sckey, title and body.
func LogToLogger(logger logrus.StdLogger, sckey, title, body string, err error, options ...logop.LoggerOption) {
	param, _ := getLoggerParamAndFields(sckey, title, body, err)
	extra := logop.NewLoggerOptions(options)

	if err != nil {
		msg := formatErrorLogger(param, err)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	} else {
		msg := formatLogger(param)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	}
}

// formatLogger formats loggerParam to logger string.
// Logs like:
// 	[Serverchan]   256B + 64.00KB | SCU0ab**********jklmn4**********9uvwxy**********g456ab | te**st
// 	            |----------------| |------------------------------------------------------| |-------|
// 	                    16                                    54                               ...
func formatLogger(param *loggerParam) string {
	length := fmt.Sprintf("%dB + %s", param.titleLength, xnumber.RenderByte(float64(param.bodyLength)))
	return fmt.Sprintf("[Serverchan] %16s | %54s | %s", length, param.maskedSckey, param.maskedTitle)
}

// formatErrorLogger formats loggerParam and error to logger string.
// Logs like:
// 	[Serverchan] Send message 't**e' to bad push token 'KEY'
// 	[Serverchan] Send duplicate message 'te**st' to 'SCU0ab**********jklmn4**********9uvwxy**********g456ab'
// 	[Serverchan] Send message 't**e' to 'SCU0ab**********jklmn4**********9uvwxy**********g456ab' failed | serverchan: respond with not success
func formatErrorLogger(param *loggerParam, err error) string {
	var msg string
	switch {
	case err == serverchan.ErrBadPushToken:
		msg = fmt.Sprintf("[Serverchan] Send message '%s' to bad push token '%s'", param.maskedTitle, param.maskedSckey)
	case err == serverchan.ErrDuplicateMessage:
		msg = fmt.Sprintf("[Serverchan] Send duplicate message '%s' to '%s'", param.maskedTitle, param.maskedSckey)
	default:
		msg = fmt.Sprintf("[Serverchan] Send message '%s' to '%s' failed | %v", param.maskedTitle, param.maskedSckey, err)
	}
	return msg
}
