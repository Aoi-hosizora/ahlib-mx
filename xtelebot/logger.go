package xtelebot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logopt"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"time"
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

// receiveLoggerParam stores some receive-event logger parameters, used in LogReceiveToLogrus and LogReceiveToLogger.
type receiveLoggerParam struct {
	endpoint     string
	messageID    int
	chatID       int64
	chatUsername string
}

// getReceiveLoggerParamAndFields returns receiveLoggerParam and logrus.Fields from given endpoint and handler's telebot.Message.
func getReceiveLoggerParamAndFields(endpoint string, message *telebot.Message) (*receiveLoggerParam, logrus.Fields) {
	param := &receiveLoggerParam{
		endpoint:     endpoint,
		messageID:    message.ID,
		chatID:       message.Chat.ID,
		chatUsername: message.Chat.Username,
	}
	fields := logrus.Fields{
		"module":        "telebot",
		"action":        "receive",
		"endpoint":      param.endpoint,
		"message_id":    param.messageID,
		"chat_id":       param.chatID,
		"chat_username": param.chatUsername,
	}
	return param, fields
}

// replyLoggerParam stores some reply-event logger parameters, used in LogReplyToLogrus and LogReplyToLogger.
type replyLoggerParam struct {
	receivedMessageID int
	repliedMessageID  int
	repliedType       string
	receivedTime      time.Time
	repliedTime       time.Time
	latency           time.Duration
	chatID            int64
	chatUsername      string
}

// getReplyLoggerParamAndFields returns replyLoggerParam and logrus.Fields from given received, replied telebot.Message and error.
func getReplyLoggerParamAndFields(received, replied *telebot.Message, err error) (*replyLoggerParam, logrus.Fields) {
	var param *replyLoggerParam
	var fields logrus.Fields

	if err == nil {
		param = &replyLoggerParam{
			receivedMessageID: received.ID,
			repliedMessageID:  replied.ID,
			repliedType:       renderMessageType(replied),
			receivedTime:      received.Time(), // <<< UnixTime
			repliedTime:       replied.Time(),
			latency:           replied.Time().Sub(received.Time()),
			chatID:            replied.Chat.ID,
			chatUsername:      replied.Chat.Username,
		}
		fields = logrus.Fields{
			"module":              "telebot",
			"action":              "reply",
			"received_message_id": param.receivedMessageID,
			"replied_message_id":  param.repliedMessageID,
			"replied_type":        param.repliedType,
			"received_time":       param.receivedTime.Format(time.RFC3339),
			"replied_time":        param.repliedTime.Format(time.RFC3339),
			"latency":             param.latency,
			"chat_id":             param.chatID,
			"chat_username":       param.chatUsername,
		}
	} else {
		param = &replyLoggerParam{
			receivedMessageID: received.ID,
			receivedTime:      received.Time(),  // <<< UnixTime
			chatID:            received.Chat.ID, // use received message for nil replied message
			chatUsername:      received.Chat.Username,
		}
		fields = logrus.Fields{
			"module":              "telebot",
			"action":              "reply",
			"received_message_id": param.receivedMessageID,
			"received_time":       param.receivedTime.Format(time.RFC3339),
			"chat_id":             param.chatID,
			"chat_username":       param.chatUsername,
			"error":               err, // <<<
		}
	}

	return param, fields
}

// sendLoggerParam stores some send-event logger parameters, used in LogSendToLogrus and LogSendToLogger.
type sendLoggerParam struct {
	sentMessageID int
	sentType      string
	sentTime      time.Time
	chatID        int64
	chatUsername  string
}

// getSendLoggerParamAndFields returns sendLoggerParam and logrus.Fields from given telebot.Chat, sent telebot.Message and error.
func getSendLoggerParamAndFields(chat *telebot.Chat, sent *telebot.Message, err error) (*sendLoggerParam, logrus.Fields) {
	var param *sendLoggerParam
	var fields logrus.Fields

	if err == nil {
		param = &sendLoggerParam{
			sentMessageID: sent.ID,
			sentType:      renderMessageType(sent),
			sentTime:      sent.Time(),
			chatID:        sent.Chat.ID, // no use of given chat
			chatUsername:  sent.Chat.Username,
		}
		fields = logrus.Fields{
			"module":          "telebot",
			"action":          "send",
			"sent_message_id": param.sentMessageID,
			"sent_type":       param.sentType,
			"sent_time":       param.sentTime,
			"chat_id":         param.chatID,
			"chat_username":   param.chatUsername,
		}
	} else {
		param = &sendLoggerParam{
			chatID:       chat.ID, // use given chat for nil sent message
			chatUsername: chat.Username,
		}
		fields = logrus.Fields{
			"module":        "telebot",
			"action":        "send",
			"chat_id":       param.chatID,
			"chat_username": param.chatUsername,
			"error":         err, // <<<
		}
	}

	return param, fields
}

// LogReceiveToLogrus logs a receive-event message to logrus.Logger using given endpoint and handler's telebot.Message.
func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, message *telebot.Message, options ...logopt.LoggerOption) {
	endpointString, ok := renderEndpoint(endpoint)
	if !ok || message == nil {
		return
	}

	param, fields := getReceiveLoggerParamAndFields(endpointString, message)
	extra := logopt.NewLoggerOptions(options)
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatReceiveLogger(param)
	extra.AddToMessage(&msg)
	entry.Info(msg)
}

// LogReplyToLogrus logs a reply-event message to logrus.Logger using given received, replied telebot.Message and error.
func LogReplyToLogrus(logger *logrus.Logger, received, replied *telebot.Message, err error, options ...logopt.LoggerOption) {
	if received == nil || (replied == nil && err == nil) {
		return
	}

	param, fields := getReplyLoggerParamAndFields(received, replied, err)
	extra := logopt.NewLoggerOptions(options)
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	if err != nil {
		msg := formatReplyErrorLogger(received, err)
		extra.AddToMessage(&msg)
		entry.Error(msg)
	} else {
		msg := formatReplyLogger(param)
		extra.AddToMessage(&msg)
		entry.Info(msg)
	}
}

// LogSendToLogrus logs a send-event message to logrus.Logger using given telebot.Chat, sent telebot.Message and error.
func LogSendToLogrus(logger *logrus.Logger, chat *telebot.Chat, sent *telebot.Message, err error, options ...logopt.LoggerOption) {
	if chat == nil || (sent == nil && err == nil) {
		return
	}

	param, fields := getSendLoggerParamAndFields(chat, sent, err)
	extra := logopt.NewLoggerOptions(options)
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	if err != nil {
		msg := formatSendErrorLogger(chat, err)
		extra.AddToMessage(&msg)
		entry.Error(msg)
	} else {
		msg := formatSendLogger(param)
		extra.AddToMessage(&msg)
		entry.Info(msg)
	}
}

// LogReceiveToLogger logs a receive-event message to logrus.StdLogger using given endpoint and handler's telebot.Message.
func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, message *telebot.Message, options ...logopt.LoggerOption) {
	endpointString, ok := renderEndpoint(endpoint)
	if !ok || message == nil {
		return
	}

	param, _ := getReceiveLoggerParamAndFields(endpointString, message)
	extra := logopt.NewLoggerOptions(options)

	msg := formatReceiveLogger(param)
	extra.AddToMessage(&msg)
	logger.Print(msg)
}

// LogReplyToLogger logs a reply-event message to logrus.StdLogger using given received, replied telebot.Message and error.
func LogReplyToLogger(logger logrus.StdLogger, received, replied *telebot.Message, err error, options ...logopt.LoggerOption) {
	if received == nil || (err == nil && replied == nil) {
		return
	}

	param, _ := getReplyLoggerParamAndFields(received, replied, err)
	extra := logopt.NewLoggerOptions(options)

	if err != nil {
		msg := formatReplyErrorLogger(received, err)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	} else {
		msg := formatReplyLogger(param)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	}
}

// LogSendToLogger logs a send-event message to logrus.StdLogger using given telebot.Chat, sent telebot.Message and error.
func LogSendToLogger(logger logrus.StdLogger, chat *telebot.Chat, sent *telebot.Message, err error, options ...logopt.LoggerOption) {
	if chat == nil || (sent == nil && err == nil) {
		return
	}

	param, _ := getSendLoggerParamAndFields(chat, sent, err)
	extra := logopt.NewLoggerOptions(options)

	if err != nil {
		msg := formatSendErrorLogger(chat, err)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	} else {
		msg := formatSendLogger(param)
		extra.AddToMessage(&msg)
		logger.Print(msg)
	}
}

// formatReceiveLogger formats receiveLoggerParam to logger string.
// Logs like:
// 	[Telebot] 3344 |                 /test-endpoint | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |                       $on_text | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |                 $rep_btn:reply | 12345678 Aoi-hosizora
// 	         |----| |------------------------------| |--------|------------|
// 	           4                    30                   ...       ...
func formatReceiveLogger(param *receiveLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %30s | %d %s",
		param.messageID, param.endpoint, param.chatID, param.chatUsername)
}

// formatReplyLogger formats replyLoggerParam to logger string.
// Logs like:
// 	[Telebot] 3345 |           2s |   t:text | 3344 | 12345678 Aoi-hosizora
// 	         |----| |------------| |--------| |----| |--------|------------|
// 	           4          12            8       4        ...       ...
func formatReplyLogger(param *replyLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4d | %d %s",
		param.repliedMessageID, param.latency.String(), param.repliedType, param.receivedMessageID, param.chatID, param.chatUsername)
}

// formatSendLogger formats sendLoggerParam to logger string.
// Logs like:
// 	[Telebot] 3346 |            x |   t:text |    x | 12345678 Aoi-hosizora
// 	         |----| |------------| |--------| |----| |--------|------------|
// 	           4          12            8       4        ...       ...
func formatSendLogger(param *sendLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4s | %d %s",
		param.sentMessageID, "x", param.sentType, "x", param.chatID, param.chatUsername)
}

// formatReplyErrorLogger formats received telebot.Message and error to logger string.
// Logs like:
// 	[Telebot] Reply to '12345678 Aoi-hosizora' failed | telegram: bot was blocked by the user (401)
func formatReplyErrorLogger(received *telebot.Message, err error) string {
	return fmt.Sprintf("[Telebot] Reply to '%d %s' failed | %v", received.Chat.ID, received.Chat.Username, err)
}

// formatSendErrorLogger formats sent telebot.Chat and error to logger string.
// Logs like:
// 	[Telebot] Send to '12345678 Aoi-hosizora' failed | telegram: bot was blocked by the user (401)
func formatSendErrorLogger(chat *telebot.Chat, err error) string {
	return fmt.Sprintf("[Telebot] Send to '%d %s' failed | %v", chat.ID, chat.Username, err)
}

// renderEndpoint renders an endpoint interface{} to string.
//
// Support endpoint types:
// 	string, telebot.InlineButton, telebot.ReplyButton
func renderEndpoint(endpoint interface{}) (string, bool) {
	epStr := ""
	switch ep := endpoint.(type) {
	case string:
		if ep == "" || ep == "\a" {
			return "", false // empty
		}
		if len(ep) > 1 && ep[0] == '\a' {
			epStr = "$on_" + ep[1:] // OnXXX
		} else {
			epStr = ep // string
		}
	case *telebot.ReplyButton:
		unique := ep.Text
		if unique == "" {
			return "", false // empty
		}
		epStr = "$rep_btn:" + unique // CallbackUnique
	case *telebot.InlineButton:
		unique := ep.Unique
		if unique == "" {
			return "", false // empty
		}
		epStr = "$inl_btn:" + unique // CallbackUnique
	default:
		return "", false // unsupported endpoint
	}

	return epStr, true
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
