package xtelebot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

// LoggerOption represents an option type for some logger functions' option, can be created by WithXXX functions.
type LoggerOption = internal.LoggerOption

// WithExtraText creates a LoggerOption to specific extra text logging in "... | extra_text" style, notes that if you use this multiple times, only the last one will be retained.
func WithExtraText(text string) LoggerOption {
	return internal.WithExtraText(text)
}

// WithExtraFields creates a LoggerOption to specific logging with extra fields, notes that if you use this multiple times, only the last one will be retained.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return internal.WithExtraFields(fields)
}

// WithExtraFieldsV creates a LoggerOption to specific logging with extra fields in variadic, notes that if you use this multiple times, only the last one will be retained.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return internal.WithExtraFieldsV(fields...)
}

// receiveLoggerParam stores some receive-event logger parameters used by LogReceiveToLogrus and LogReceiveToLogger.
type receiveLoggerParam struct {
	endpoint   string
	receivedID int
	chatID     int64
	chatName   string
}

// extractReceiveLoggerData extracts and returns receiveLoggerParam and logrus.Fields using given parameters.
func extractReceiveLoggerData(endpoint string, message *telebot.Message) (*receiveLoggerParam, logrus.Fields) {
	param := &receiveLoggerParam{
		endpoint:   endpoint,
		receivedID: message.ID,
		chatID:     message.Chat.ID,
		chatName:   message.Chat.Username,
	}
	fields := logrus.Fields{
		"module":      "telebot",
		"action":      "receive",
		"endpoint":    param.endpoint,
		"received_id": param.receivedID,
		"chat_id":     param.chatID,
		"chat_name":   param.chatName,
	}
	return param, fields
}

// formatReceiveLogger formats given receiveLoggerParam to string for LogReceiveToLogrus and LogReceiveToLogger.
//
// Logs like:
// 	[Telebot] 3344 |                 /test-endpoint | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |                       $on_text | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |                 $rep_btn:reply | 12345678 Aoi-hosizora
// 	         |----| |------------------------------| |--------|------------|
// 	           4                    30                  ...        ...
func formatReceiveLogger(param *receiveLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %30s | %d %s", param.receivedID, param.endpoint, param.chatID, param.chatName)
}

// replyLoggerParam stores some reply-event logger parameters used by LogReplyToLogrus and LogReplyToLogger.
type replyLoggerParam struct {
	receivedID  int
	repliedID   int
	repliedType string
	latency     string
	chatID      int64
	chatName    string
	errorMsg    string
}

// extractReplyLoggerData extracts and returns replyLoggerParam and logrus.Fields using given parameters.
func extractReplyLoggerData(received, replied *telebot.Message, err error) (*replyLoggerParam, logrus.Fields) {
	var param *replyLoggerParam
	var fields logrus.Fields

	if err == nil {
		latency := replied.Time().Sub(received.Time())
		param = &replyLoggerParam{
			receivedID:  received.ID,
			repliedID:   replied.ID,
			repliedType: renderMessageType(replied),
			latency:     latency.String(),
			chatID:      replied.Chat.ID,
			chatName:    replied.Chat.Username,
		}
		fields = logrus.Fields{
			"module":        "telebot",
			"action":        "reply",
			"received_id":   param.receivedID,
			"replied_id":    param.repliedID,
			"replied_type":  param.repliedType,
			"received_time": received.Time().Format(time.RFC3339),
			"replied_time":  replied.Time().Format(time.RFC3339),
			"latency":       param.latency,
			"chat_id":       param.chatID,
			"chat_name":     param.chatName,
		}
	} else {
		param = &replyLoggerParam{
			receivedID: received.ID,
			chatID:     received.Chat.ID,
			chatName:   received.Chat.Username,
			errorMsg:   err.Error(), // <<<
		}
		fields = logrus.Fields{
			"module":        "telebot",
			"action":        "reply",
			"received_id":   param.receivedID,
			"received_time": received.Time().Format(time.RFC3339),
			"chat_id":       param.chatID,
			"chat_name":     param.chatName,
			"error_msg":     param.errorMsg, // <<<
		}
	}

	return param, fields
}

// formatReplyLogger formats given replyLoggerParam to string for LogReplyToLogrus and LogReplyToLogger.
//
// Logs like:
// 	[Telebot] 3345 |           2s |   t:text | 3344 | 12345678 Aoi-hosizora
// 	         |----| |------------| |--------| |----| |--------|------------|
// 	           4          12            8       4       ...        ...
func formatReplyLogger(param *replyLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4d | %d %s", param.repliedID, param.latency, param.repliedType, param.receivedID, param.chatID, param.chatName)
}

// formatReplyErrorLogger formats given replyLoggerParam with error to string for LogReplyToLogrus and LogReplyToLogger.
//
// Logs like:
// 	[Telebot] Reply to message 3344 from chat '12345678 Aoi-hosizora' failed | telegram: bot was blocked by the user (401)
func formatReplyErrorLogger(param *replyLoggerParam) string {
	return fmt.Sprintf("[Telebot] Reply to message %d from chat '%d %s' failed | %s", param.receivedID, param.chatID, param.chatName, param.errorMsg)
}

// sendLoggerParam stores some send-event logger parameters used bu LogSendToLogrus and LogSendToLogger.
type sendLoggerParam struct {
	sentID   int
	sentType string
	chatID   int64
	chatName string
	errorMsg string
}

// extractSendLoggerData extracts and returns sendLoggerParam and logrus.Fields using given parameters.
func extractSendLoggerData(sent *telebot.Message, chat *telebot.Chat, err error) (*sendLoggerParam, logrus.Fields) {
	var param *sendLoggerParam
	var fields logrus.Fields

	if err == nil {
		param = &sendLoggerParam{
			sentID:   sent.ID,
			sentType: renderMessageType(sent),
			chatID:   sent.Chat.ID, // use chat from `sent` rather than `chat`
			chatName: sent.Chat.Username,
		}
		fields = logrus.Fields{
			"module":    "telebot",
			"action":    "send",
			"sent_id":   param.sentID,
			"sent_type": param.sentType,
			"sent_time": sent.Time().Format(time.RFC3339),
			"chat_id":   param.chatID,
			"chat_name": param.chatName,
		}
	} else {
		param = &sendLoggerParam{
			chatID:   chat.ID,
			chatName: chat.Username,
			errorMsg: err.Error(), // <<<
		}
		fields = logrus.Fields{
			"module":    "telebot",
			"action":    "send",
			"chat_id":   param.chatID,
			"chat_name": param.chatName,
			"error_msg": param.errorMsg, // <<<
		}
	}

	return param, fields
}

// formatSendLogger formats given sendLoggerParam to string for LogSendToLogrus and LogSendToLogger.
//
// Logs like:
// 	[Telebot] 3346 |            x |   t:text |    x | 12345678 Aoi-hosizora
// 	         |----| |------------| |--------| |----| |--------|------------|
// 	           4          12            8       4       ...        ...
func formatSendLogger(param *sendLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4s | %d %s", param.sentID, "x", param.sentType, "x", param.chatID, param.chatName)
}

// formatSendLogger formats given sendLoggerParam with error to string for LogSendToLogrus and LogSendToLogger.
//
// Logs like:
// 	[Telebot] Send message to chat '12345678 Aoi-hosizora' failed | telegram: bot was blocked by the user (401)
func formatSendErrorLogger(param *sendLoggerParam) string {
	return fmt.Sprintf("[Telebot] Send message to chat '%d %s' failed | %s", param.chatID, param.chatName, param.errorMsg)
}

// LogReceiveToLogrus logs a receive-event message to logrus.Logger using given endpoint and handler's telebot.Message.
func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, message *telebot.Message, options ...LoggerOption) {
	if logger == nil || message == nil {
		return
	}
	endpointString, ok := renderEndpoint(endpoint)
	if !ok {
		return
	}
	p, f := extractReceiveLoggerData(endpointString, message)
	m := formatReceiveLogger(p)

	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	logger.WithFields(f).Info(m)
}

// LogReceiveToLogger logs a receive-event message to logrus.StdLogger using given endpoint and handler's telebot.Message.
func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, received *telebot.Message, options ...LoggerOption) {
	if logger == nil || received == nil {
		return
	}
	ep, ok := renderEndpoint(endpoint)
	if !ok {
		return
	}
	p, _ := extractReceiveLoggerData(ep, received)
	m := formatReceiveLogger(p)

	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// LogReplyToLogrus logs a reply-event message to logrus.Logger using given received, replied telebot.Message and error, `replied` is expected to be nii when error aroused.
func LogReplyToLogrus(logger *logrus.Logger, received, replied *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || received == nil || (replied == nil && err == nil) {
		return
	}
	p, f := extractReplyLoggerData(received, replied, err)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToFields(f)

	if err != nil {
		m := formatReplyErrorLogger(p)
		extra.ApplyToMessage(&m)
		logger.WithFields(f).Error(m)
	} else {
		m := formatReplyLogger(p)
		extra.ApplyToMessage(&m)
		logger.WithFields(f).Info(m)
	}
}

// LogReplyToLogger logs a reply-event message to logrus.StdLogger using given received, replied telebot.Message and error, `replied` is expected to be nii when error aroused.
func LogReplyToLogger(logger logrus.StdLogger, received, replied *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || received == nil || (err == nil && replied == nil) {
		return
	}
	p, _ := extractReplyLoggerData(received, replied, err)
	extra := internal.BuildLoggerOptions(options)

	if err != nil {
		m := formatReplyErrorLogger(p)
		extra.ApplyToMessage(&m)
		logger.Print(m)
	} else {
		m := formatReplyLogger(p)
		extra.ApplyToMessage(&m)
		logger.Print(m)
	}
}

// LogSendToLogrus logs a send-event message to logrus.Logger using given telebot.Chat, sent telebot.Message and error, `sent` is expected to be nii when error aroused.
func LogSendToLogrus(logger *logrus.Logger, chat *telebot.Chat, sent *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || chat == nil || (sent == nil && err == nil) {
		return
	}
	p, f := extractSendLoggerData(sent, chat, err)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToFields(f)

	if err != nil {
		m := formatSendErrorLogger(p)
		extra.ApplyToMessage(&m)
		logger.WithFields(f).Error(m)
	} else {
		m := formatSendLogger(p)
		extra.ApplyToMessage(&m)
		logger.WithFields(f).Info(m)
	}
}

// LogSendToLogger logs a send-event message to logrus.StdLogger using given telebot.Chat, sent telebot.Message and error, `sent` is expected to be nii when error aroused.
func LogSendToLogger(logger logrus.StdLogger, chat *telebot.Chat, sent *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || chat == nil || (sent == nil && err == nil) {
		return
	}
	p, _ := extractSendLoggerData(sent, chat, err)
	extra := internal.BuildLoggerOptions(options)

	if err != nil {
		m := formatSendErrorLogger(p)
		extra.ApplyToMessage(&m)
		logger.Print(m)
	} else {
		m := formatSendLogger(p)
		extra.ApplyToMessage(&m)
		logger.Print(m)
	}
}

// renderEndpoint renders an endpoint interface{} to string, only supported string, telebot.InlineButton and telebot.ReplyButton endpoint types.
func renderEndpoint(endpoint interface{}) (string, bool) {
	switch ep := endpoint.(type) {
	case string:
		if ep == "" || ep == "\a" {
			return "", false // empty
		}
		out := ""
		if len(ep) > 1 && ep[0] == '\a' {
			out = "$on_" + ep[1:] // OnXXX
		} else {
			out = ep // string
		}
		return out, true
	case *telebot.ReplyButton:
		unique := ep.Text
		if unique == "" {
			return "", false // empty
		}
		return "$rep_btn:" + unique, true // CallbackUnique
	case *telebot.InlineButton:
		unique := ep.Unique
		if unique == "" {
			return "", false // empty
		}
		return "$inl_btn:" + unique, true // CallbackUnique
	default:
		return "", false // unsupported endpoint
	}
}

// renderMessageType renders a telebot.Message's type, see telebot.Sendable.
func renderMessageType(m *telebot.Message) string {
	typ := "text" // default
	switch {
	case m.Photo != nil:
		typ = "photo"
	case m.Sticker != nil:
		typ = "stk"
	case m.Video != nil:
		typ = "video"
	case m.Audio != nil:
		typ = "audio"
	case m.Voice != nil:
		typ = "voice"
	case m.Location != nil:
		typ = "loc"

	case m.Animation != nil:
		typ = "anime"
	case m.Dice != nil:
		typ = "dice"
	case m.Document != nil:
		typ = "doc"
	case m.Invoice != nil:
		typ = "inv"
	case m.Poll != nil:
		typ = "poll"
	case m.Venue != nil:
		typ = "venue"
	case m.VideoNote != nil:
		typ = "vnote"
	}

	return "t:" + typ // t:xxx
}
