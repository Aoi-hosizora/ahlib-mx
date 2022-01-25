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

// WithExtraText creates a LoggerOption to specific extra text logging in "...extra_text" style, notes that if you use this multiple times, only the last one will be retained.
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

// =======
// receive
// =======

// ReceiveLoggerParam stores some receive-event logger parameters and is used by LogReceiveToLogrus and LogReceiveToLogger.
type ReceiveLoggerParam struct {
	Endpoint     string
	Action       string
	Received     *telebot.Message
	ReceivedID   int
	ReceivedTime time.Time
	ChatID       int64
	ChatName     string
}

var (
	// FormatReceiveFunc is a custom ReceiveLoggerParam's format function for LogReceiveToLogrus and LogReceiveToLogger.
	FormatReceiveFunc func(p *ReceiveLoggerParam) string

	// FieldifyReceiveFunc is a custom ReceiveLoggerParam's fieldify function for LogReceiveToLogrus.
	FieldifyReceiveFunc func(p *ReceiveLoggerParam) logrus.Fields
)

// extractReceiveLoggerParam extracts and returns ReceiveLoggerParam using given parameters.
func extractReceiveLoggerParam(endpoint string, received *telebot.Message) *ReceiveLoggerParam {
	return &ReceiveLoggerParam{
		Endpoint:     endpoint,
		Action:       "received",
		Received:     received,
		ReceivedID:   received.ID,
		ReceivedTime: received.Time(),
		ChatID:       received.Chat.ID,
		ChatName:     received.Chat.Username,
	}
}

// formatReceiveLoggerParam formats given ReceiveLoggerParam to string for LogReceiveToLogrus and LogReceiveToLogger.
//
// The default format logs like:
// 	[Telebot] 3344 |                 /test-endpoint | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |                       $on_text | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |          $rep_btn:reply_button | 12345678 Aoi-hosizora
// 	[Telebot] 3344 |         $inl_btn:inline_button | 12345678 Aoi-hosizora
// 	         |----| |------------------------------| |--------|------------|
// 	           4                    30                  ...        ...
func formatReceiveLoggerParam(p *ReceiveLoggerParam) string {
	if FormatReceiveFunc != nil {
		return FormatReceiveFunc(p)
	}
	return fmt.Sprintf("[Telebot] %4d | %30s | %d %s", p.ReceivedID, p.Endpoint, p.ChatID, p.ChatName)
}

// fieldifyReceiveLoggerParam fieldifies given ReceiveLoggerParam to logrus.Fields for LogReceiveToLogrus.
//
// The default contains the following fields: module, action, endpoint, received_id, chat_id, chat_name.
func fieldifyReceiveLoggerParam(p *ReceiveLoggerParam) logrus.Fields {
	if FieldifyReceiveFunc != nil {
		return FieldifyReceiveFunc(p)
	}
	return logrus.Fields{
		"module":        "telebot",
		"action":        p.Action,
		"endpoint":      p.Endpoint,
		"received_id":   p.ReceivedID,
		"received_time": p.ReceivedTime.Format(time.RFC3339),
		"chat_id":       p.ChatID,
		"chat_name":     p.ChatName,
	}
}

// LogReceiveToLogrus logs a receive-event message to logrus.Logger using given endpoint and handler's received telebot.Message.
func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, received *telebot.Message, options ...LoggerOption) {
	if logger == nil || received == nil {
		return
	}
	epString, ok := formatEndpoint(endpoint)
	if !ok {
		return
	}
	p := extractReceiveLoggerParam(epString, received)
	m := formatReceiveLoggerParam(p)
	f := fieldifyReceiveLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	logger.WithFields(f).Info(m)
}

// LogReceiveToLogger logs a receive-event message to logrus.StdLogger using given endpoint and handler's received telebot.Message.
func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, received *telebot.Message, options ...LoggerOption) {
	if logger == nil || received == nil {
		return
	}
	epString, ok := formatEndpoint(endpoint)
	if !ok {
		return
	}
	p := extractReceiveLoggerParam(epString, received)
	m := formatReceiveLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// =====
// reply
// =====

// ReplyLoggerParam stores some reply-event logger parameters and is used by LogReplyToLogrus and LogReplyToLogger.
type ReplyLoggerParam struct {
	Action       string
	Received     *telebot.Message
	Replied      *telebot.Message
	ReceivedID   int
	RepliedID    int
	RepliedType  string
	ReceivedTime time.Time
	RepliedTime  time.Time
	Latency      time.Duration
	ChatID       int64
	ChatName     string
	ErrorMsg     string
}

var (
	// FormatReplyFunc is a custom ReplyLoggerParam's format function for LogReplyToLogrus and LogReplyToLogger.
	FormatReplyFunc func(p *ReplyLoggerParam) string

	// FieldifyReplyFunc is a custom ReplyLoggerParam's fieldify function for LogReplyToLogrus.
	FieldifyReplyFunc func(p *ReplyLoggerParam) logrus.Fields
)

// extractReplyLoggerParam extracts and returns ReplyLoggerParam using given parameters.
func extractReplyLoggerParam(received, replied *telebot.Message, err error) *ReplyLoggerParam {
	if err != nil {
		return &ReplyLoggerParam{
			Action:       "reply",
			Received:     received,
			ReceivedID:   received.ID,
			ReceivedTime: received.Time(),
			ChatID:       received.Chat.ID,
			ChatName:     received.Chat.Username,
			ErrorMsg:     err.Error(), // <<<
		}
	}
	return &ReplyLoggerParam{
		Action:       "reply",
		Received:     received,
		Replied:      replied,
		ReceivedID:   received.ID,
		RepliedID:    replied.ID,
		RepliedType:  formatMessageType(replied),
		ReceivedTime: received.Time(),
		RepliedTime:  replied.Time(),
		Latency:      replied.Time().Sub(received.Time()),
		ChatID:       replied.Chat.ID,
		ChatName:     replied.Chat.Username,
	}
}

// formatReplyLoggerParam formats given ReplyLoggerParam to string for LogReplyToLogrus and LogReplyToLogger.
//
// The default format logs like:
// 	[Telebot] Reply to message 3344 from chat '12345678 Aoi-hosizora' failed | telegram: bot was blocked by the user (401)
// 	[Telebot] 3345 |           2s |   t:text | 3344 | 12345678 Aoi-hosizora
// 	         |----| |------------| |--------| |----| |--------|------------|
// 	           4          12            8       4       ...        ...
func formatReplyLoggerParam(p *ReplyLoggerParam) string {
	if FormatReplyFunc != nil {
		return FormatReplyFunc(p)
	}
	if p.ErrorMsg != "" {
		return fmt.Sprintf("[Telebot] Reply to message %d from chat '%d %s' failed | %s", p.ReceivedID, p.ChatID, p.ChatName, p.ErrorMsg)
	}
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4d | %d %s", p.RepliedID, p.Latency.String(), p.RepliedType, p.ReceivedID, p.ChatID, p.ChatName)
}

// fieldifyReplyLoggerParam fieldifies given ReplyLoggerParam to logrus.Fields for LogReplyToLogrus.
//
// The default contains the following fields: module, action, received_id, replied_id, replied_type, received_time, replied_time, latency, chat_id, chat_name.
func fieldifyReplyLoggerParam(p *ReplyLoggerParam) logrus.Fields {
	if FieldifyReplyFunc != nil {
		return FieldifyReplyFunc(p)
	}
	if p.ErrorMsg != "" {
		return logrus.Fields{
			"module":        "telebot",
			"action":        p.Action,
			"received_id":   p.ReceivedID,
			"received_time": p.RepliedTime.Format(time.RFC3339),
			"chat_id":       p.ChatID,
			"chat_name":     p.ChatName,
			"error_msg":     p.ErrorMsg, // <<<
		}
	}
	return logrus.Fields{
		"module":        "telebot",
		"action":        p.Action,
		"received_id":   p.ReceivedID,
		"replied_id":    p.RepliedID,
		"replied_type":  p.RepliedType,
		"received_time": p.ReceivedTime.Format(time.RFC3339),
		"replied_time":  p.RepliedTime.Format(time.RFC3339),
		"latency":       p.Latency.String(),
		"chat_id":       p.ChatID,
		"chat_name":     p.ChatName,
	}
}

// LogReplyToLogrus logs a reply-event message to logrus.Logger using given received, replied telebot.Message and error, `replied` is expected to be nii when error aroused.
func LogReplyToLogrus(logger *logrus.Logger, received, replied *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || received == nil || (replied == nil && err == nil) {
		return
	}
	p := extractReplyLoggerParam(received, replied, err)
	m := formatReplyLoggerParam(p)
	f := fieldifyReplyLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	if p.ErrorMsg != "" {
		logger.WithFields(f).Error(m)
	} else {
		logger.WithFields(f).Info(m)
	}
}

// LogReplyToLogger logs a reply-event message to logrus.StdLogger using given received, replied telebot.Message and error, `replied` is expected to be nii when error aroused.
func LogReplyToLogger(logger logrus.StdLogger, received, replied *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || received == nil || (err == nil && replied == nil) {
		return
	}
	p := extractReplyLoggerParam(received, replied, err)
	m := formatReplyLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// ====
// send
// ====

// SendLoggerParam stores some send-event logger parameters and is used by LogSendToLogrus and LogSendToLogger.
type SendLoggerParam struct {
	Action   string
	Chat     *telebot.Chat
	Sent     *telebot.Message
	SentID   int
	SentType string
	SentTime time.Time
	ChatID   int64
	ChatName string
	ErrorMsg string
}

var (
	// FormatSendFunc is a custom SendLoggerParam's format function for LogSendToLogrus and LogSendToLogger.
	FormatSendFunc func(p *SendLoggerParam) string

	// FieldifySendFunc is a custom SendLoggerParam's fieldify function for LogSendToLogrus.
	FieldifySendFunc func(p *SendLoggerParam) logrus.Fields
)

// extractSendLoggerParam extracts and returns SendLoggerParam using given parameters.
func extractSendLoggerParam(chat *telebot.Chat, sent *telebot.Message, err error) *SendLoggerParam {
	if err != nil {
		return &SendLoggerParam{
			Action:   "send",
			Chat:     chat, // use `chat`
			ChatID:   chat.ID,
			ChatName: chat.Username,
			ErrorMsg: err.Error(), // <<<
		}
	}
	return &SendLoggerParam{
		Action:   "send",
		Chat:     sent.Chat, // use `sent.Chat`
		Sent:     sent,
		SentID:   sent.ID,
		SentType: formatMessageType(sent),
		SentTime: sent.Time(),
		ChatID:   sent.Chat.ID,
		ChatName: sent.Chat.Username,
	}
}

// formatSendLoggerParam formats given SendLoggerParam to string for LogSendToLogrus and LogSendToLogger.
//
// The default format logs like:
// 	[Telebot] Send message to chat '12345678 Aoi-hosizora' failed | telegram: bot was blocked by the user (401)
// 	[Telebot] 3346 |            x |   t:text |    x | 12345678 Aoi-hosizora
// 	         |----| |------------| |--------| |----| |--------|------------|
// 	           4          12            8       4       ...        ...
func formatSendLoggerParam(p *SendLoggerParam) string {
	if FormatSendFunc != nil {
		return FormatSendFunc(p)
	}
	if p.ErrorMsg != "" {
		return fmt.Sprintf("[Telebot] Send message to chat '%d %s' failed | %s", p.ChatID, p.ChatName, p.ErrorMsg)
	}
	return fmt.Sprintf("[Telebot] %4d | %12s | %8s | %4s | %d %s", p.SentID, "x", p.SentType, "x", p.ChatID, p.ChatName)
}

// fieldifySendLoggerParam fieldifies given ReplyLoggerParam to logrus.Fields.
//
// The default contains the following fields: module, action, sent_id, sent_type, sent_time, chat_id, chat_name.
func fieldifySendLoggerParam(p *SendLoggerParam) logrus.Fields {
	if FieldifySendFunc != nil {
		return FieldifySendFunc(p)
	}
	if p.ErrorMsg != "" {
		return logrus.Fields{
			"module":    "telebot",
			"action":    p.Action,
			"chat_id":   p.ChatID,
			"chat_name": p.ChatName,
			"error_msg": p.ErrorMsg, // <<<
		}
	}
	return logrus.Fields{
		"module":    "telebot",
		"action":    p.Action,
		"sent_id":   p.SentID,
		"sent_type": p.SentType,
		"sent_time": p.SentTime.Format(time.RFC3339),
		"chat_id":   p.ChatID,
		"chat_name": p.ChatName,
	}
}

// LogSendToLogrus logs a send-event message to logrus.Logger using given telebot.Chat, sent telebot.Message and error, `sent` is expected to be nii when error aroused.
func LogSendToLogrus(logger *logrus.Logger, chat *telebot.Chat, sent *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || chat == nil || (sent == nil && err == nil) {
		return
	}
	p := extractSendLoggerParam(chat, sent, err)
	m := formatSendLoggerParam(p)
	f := fieldifySendLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	if p.ErrorMsg != "" {
		logger.WithFields(f).Error(m)
	} else {
		logger.WithFields(f).Info(m)
	}
}

// LogSendToLogger logs a send-event message to logrus.StdLogger using given telebot.Chat, sent telebot.Message and error, `sent` is expected to be nii when error aroused.
func LogSendToLogger(logger logrus.StdLogger, chat *telebot.Chat, sent *telebot.Message, err error, options ...LoggerOption) {
	if logger == nil || chat == nil || (sent == nil && err == nil) {
		return
	}
	p := extractSendLoggerParam(chat, sent, err)
	m := formatSendLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// ========
// internal
// ========

// formatEndpoint formats an endpoint interface{} to string, only supported string, telebot.InlineButton and telebot.ReplyButton endpoint types.
func formatEndpoint(endpoint interface{}) (string, bool) {
	switch ep := endpoint.(type) {
	case string:
		if len(ep) <= 1 || (ep[0] != '/' && ep[0] != '\a') {
			return "", false // empty
		}
		out := ""
		if ep[0] == '\a' {
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

// formatMessageType formats a telebot.Message's type, visit telebot.Sendable.
func formatMessageType(m *telebot.Message) string {
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
