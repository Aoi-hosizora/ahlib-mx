package xtelebot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"time"
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

// receiveLoggerParam stores some receive-event logger parameters, used in LogReceiveToLogrus and LogReceiveToLogger.
type receiveLoggerParam struct {
	endpoint     string
	messageID    int
	chatID       int64
	chatTitle    string
	chatUsername string
}

// getReceiveLoggerParam returns receiveLoggerParam from given endpoint and handler's telebot.Message.
func getReceiveLoggerParam(endpoint string, message *telebot.Message) *receiveLoggerParam {
	return &receiveLoggerParam{
		endpoint:     endpoint,
		messageID:    message.ID,
		chatID:       message.Chat.ID,
		chatTitle:    message.Chat.Title,
		chatUsername: message.Chat.Username,
	}
}

// replyLoggerParam stores some reply-event logger parameters, used in LogReplyToLogrus and LogReplyToLogger.
type replyLoggerParam struct {
	receivedMessageID int
	repliedMessageID  int
	repliedType       string
	receivedTime      time.Time // <<<
	repliedTime       time.Time
	latency           time.Duration
	chatID            int64
	chatTitle         string
	chatUsername      string
}

// getReplyLoggerParam returns replyLoggerParam from given received and replied telebot.Message.
func getReplyLoggerParam(received, replied *telebot.Message) *replyLoggerParam {
	return &replyLoggerParam{
		receivedMessageID: received.ID,
		repliedMessageID:  replied.ID,
		repliedType:       renderMessageType(replied),
		receivedTime:      received.Time(), // <<< UnixTime
		repliedTime:       replied.Time(),
		latency:           replied.Time().Sub(received.Time()),
		chatID:            replied.Chat.ID,
		chatTitle:         replied.Chat.Title,
		chatUsername:      replied.Chat.Username,
	}
}

// sendLoggerParam stores some send-event logger parameters, used in LogSendToLogrus and LogSendToLogger.
type sendLoggerParam struct {
	sentMessageID int
	sentType      string
	sentTime      time.Time
	chatID        int64
	chatTitle     string
	chatUsername  string
}

// getReplyLoggerParam returns replyLoggerParam from given sent telebot.Message.
func getSendLoggerParam(sent *telebot.Message) *sendLoggerParam {
	return &sendLoggerParam{
		sentMessageID: sent.ID,
		sentType:      renderMessageType(sent),
		sentTime:      sent.Time(),
		chatID:        sent.Chat.ID,
		chatTitle:     sent.Chat.Title,
		chatUsername:  sent.Chat.Username,
	}
}

// LogReceiveToLogrus logs a receive-event message to logrus.Logger using given endpoint and telebot.Message.
//
// Support endpoint types:
// 	string, telebot.InlineButton, telebot.ReplyButton
// Not support handler's type:
// 	telebot.Query
func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, message *telebot.Message, options ...logop.LoggerOption) {
	unique, ok := renderEndpoint(endpoint)
	if !ok {
		return // unsupported endpoint
	}

	param := getReceiveLoggerParam(unique, message)
	extra := logop.NewLoggerExtra(options...)

	fields := logrus.Fields{
		"module":        "telebot",
		"action":        "receive",
		"endpoint":      param.endpoint,
		"message_id":    param.messageID,
		"chat_id":       param.chatID,
		"chat_title":    param.chatTitle,
		"chat_username": param.chatUsername,
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatReceiveLogger(param)
	extra.AddToMessage(&msg)
	entry.Info(msg)
}

// LogReplyToLogrus logs a reply-event message to logrus.Logger using given received telebot.Message and replied telebot.Message with error.
func LogReplyToLogrus(logger *logrus.Logger, received *telebot.Message, sent *telebot.Message, err error, options ...logop.LoggerOption) {
	if received == nil || (err == nil && sent == nil) {
		return // ignore
	}
	if err != nil {
		msg := fmt.Sprintf("[Telebot] Reply to %d %s error: %v", sent.Chat.ID, sent.Chat.Username, err)
		logger.Error(msg)
		return
	}

	param := getReplyLoggerParam(received, sent)
	extra := logop.NewLoggerExtra(options...)

	fields := logrus.Fields{
		"module":              "telebot",
		"action":              "reply",
		"received_message_id": param.receivedMessageID,
		"replied_message_id":  param.repliedMessageID,
		"replied_type":        param.repliedType,
		"received_time":       param.receivedTime,
		"replied_time":        param.repliedTime,
		"latency":             param.latency,
		"chat_id":             param.chatID,
		"chat_title":          param.chatTitle,
		"chat_username":       param.chatUsername,
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatReplyLogger(param)
	extra.AddToMessage(&msg)
	entry.Info(msg)
}

// LogSendToLogrus logs a send-event message to logrus.Logger using given sent telebot.Message with error.
func LogSendToLogrus(logger *logrus.Logger, sent *telebot.Message, err error, options ...logop.LoggerOption) {
	if err == nil && sent == nil {
		return // ignore
	}
	if err != nil {
		msg := fmt.Sprintf("[Telebot] Send to %d %s error: %v", sent.Chat.ID, sent.Chat.Username, err)
		logger.Error(msg)
		return
	}

	param := getSendLoggerParam(sent)
	extra := logop.NewLoggerExtra(options...)

	fields := logrus.Fields{
		"module":          "telebot",
		"action":          "send",
		"sent_message_id": param.sentMessageID,
		"sent_type":       param.sentType,
		"sent_time":       param.sentTime,
		"chat_id":         param.chatID,
		"chat_title":      param.chatTitle,
		"chat_username":   param.chatUsername,
	}
	extra.AddToFields(fields)
	entry := logger.WithFields(fields)

	msg := formatSendLogger(param)
	extra.AddToMessage(&msg)
	entry.Info(msg)
}

// LogReceiveToLogger logs a receive-event message to logrus.StdLogger using given endpoint and telebot.Message.
func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, message *telebot.Message, options ...logop.LoggerOption) {
	unique, ok := renderEndpoint(endpoint)
	if !ok {
		return // unsupported endpoint
	}

	param := getReceiveLoggerParam(unique, message)
	extra := logop.NewLoggerExtra(options...)

	msg := formatReceiveLogger(param)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}

// LogReplyToLogger logs a reply-event message to logrus.StdLogger using given received telebot.Message and replied telebot.Message with error.
func LogReplyToLogger(logger logrus.StdLogger, received, replied *telebot.Message, err error, options ...logop.LoggerOption) {
	if received == nil || (err == nil && replied == nil) {
		return // ignore
	}
	if err != nil {
		msg := fmt.Sprintf("[Telebot] Reply to %d %s failed: %v", replied.Chat.ID, replied.Chat.Username, err)
		logger.Println(msg)
		return
	}

	param := getReplyLoggerParam(received, replied)
	extra := logop.NewLoggerExtra(options...)

	msg := formatReplyLogger(param)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}

// LogSendToLogger logs a send-event message to logrus.StdLogger using given sent telebot.Message with error.
func LogSendToLogger(logger logrus.StdLogger, sent *telebot.Message, err error, options ...logop.LoggerOption) {
	if err == nil && sent == nil {
		return // ignore
	}
	if err != nil {
		msg := fmt.Sprintf("[Telebot] Send to %d %s failed: %v", sent.Chat.ID, sent.Chat.Username, err)
		logger.Println(msg)
		return
	}

	param := getSendLoggerParam(sent)
	extra := logop.NewLoggerExtra(options...)

	msg := formatSendLogger(param)
	extra.AddToMessage(&msg)
	logger.Println(msg)
}

// formatReceiveLogger formats receiveLoggerParam to string for logger.
// Logs like:
// 	[Telebot] 3344 |                       $on_text | 12345678 Aoi-hosizora
func formatReceiveLogger(param *receiveLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %30s | %d %s",
		param.messageID, param.endpoint, param.chatID, param.chatUsername)
}

// formatReplyLogger formats replyLoggerParam to string for logger.
// Logs like:
// 	[Telebot] 3345 |           2s |   t:text | 3344 | 12345678 Aoi-hosizora
func formatReplyLogger(param *replyLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4d | %d %s",
		param.repliedMessageID, param.latency.String(), param.repliedType, param.receivedMessageID, param.chatID, param.chatUsername)
}

// formatSendLogger formats sendLoggerParam to string for logger.
// Logs like:
// 	[Telebot] 3346 |            x |   t:text |    x | 12345678 Aoi-hosizora
func formatSendLogger(param *sendLoggerParam) string {
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4s | %d %s",
		param.sentMessageID, "x", param.sentType, "x", param.chatID, param.chatUsername)
}

// renderEndpoint renders an endpoint interface{} to string.
//
// Support endpoint types:
// 	string, telebot.InlineButton, telebot.ReplyButton
func renderEndpoint(endpoint interface{}) (string, bool) {
	unique := ""
	switch ep := endpoint.(type) {
	case string:
		if ep == "" || ep == "\a" {
			return "", false // empty
		}
		if len(ep) > 1 && ep[0] == '\a' {
			unique = "$on_" + unique[1:] // OnXXX string
		} else {
			unique = ep // string
		}
	case telebot.ReplyButton:
		unique = "$rep_btn:" + ep.Text // CallbackUnique
	case telebot.InlineButton:
		unique = "$inl_btn:" + ep.Unique // CallbackUnique
	default:
		return "", false // unsupported
	}

	return unique, true
}

// renderMessageType renders a telebot.Message's type, see telebot.Sendable.
func renderMessageType(m *telebot.Message) string {
	typ := "text"
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
	return "t:" + typ
}
