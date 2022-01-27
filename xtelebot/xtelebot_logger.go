package xtelebot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"reflect"
	"runtime"
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
	// origin
	Endpoint interface{}
	Chat     *telebot.Chat
	Message  *telebot.Message

	// field
	FormattedEp string
	ChatID      int64
	ChatName    string
	MessageID   int
	MessageTime time.Time
}

var (
	// FormatReceiveFunc is a custom ReceiveLoggerParam's format function for LogReceiveToLogrus and LogReceiveToLogger.
	FormatReceiveFunc func(p *ReceiveLoggerParam) string

	// FieldifyReceiveFunc is a custom ReceiveLoggerParam's fieldify function for LogReceiveToLogrus.
	FieldifyReceiveFunc func(p *ReceiveLoggerParam) logrus.Fields
)

// extractReceiveLoggerParam extracts and returns ReceiveLoggerParam using given parameters.
func extractReceiveLoggerParam(endpoint interface{}, received *telebot.Message) *ReceiveLoggerParam {
	formatted, ok := formatEndpoint(endpoint)
	if !ok {
		return nil
	}
	return &ReceiveLoggerParam{
		Endpoint: endpoint,
		Chat:     received.Chat,
		Message:  received,

		FormattedEp: formatted,
		ChatID:      received.Chat.ID,
		ChatName:    received.Chat.Username,
		MessageID:   received.ID,
		MessageTime: received.Time(),
	}
}

// formatReceiveLoggerParam formats given ReceiveLoggerParam to string for LogReceiveToLogrus and LogReceiveToLogger.
//
// The default format logs like:
// 	[Telebot] recv | msg#    3344 |               /test-endpoint | 12345678 Aoi-hosizora
// 	[Telebot] recv | msg#    3345 |                     $on_text | 12345678 Aoi-hosizora
// 	[Telebot] recv | msg#    3346 |         $rep_btn:button_text | 12345678 Aoi-hosizora
// 	[Telebot] recv | msg#    3347 |       $inl_btn:button_unique | 12345678 Aoi-hosizora
// 	                     |-------| |----------------------------| |---------------------|
// 	                         7                   28                         ...
func formatReceiveLoggerParam(p *ReceiveLoggerParam) string {
	if FormatReceiveFunc != nil {
		return FormatReceiveFunc(p)
	}
	return fmt.Sprintf("[Telebot] recv | msg# %7s | %28s | %d %s", xnumber.Itoa(p.MessageID), p.FormattedEp, p.ChatID, p.ChatName)
}

// fieldifyReceiveLoggerParam fieldifies given ReceiveLoggerParam to logrus.Fields for LogReceiveToLogrus.
//
// The default contains the following fields: module, action, endpoint, received_id, chat_id, chat_name.
func fieldifyReceiveLoggerParam(p *ReceiveLoggerParam) logrus.Fields {
	if FieldifyReceiveFunc != nil {
		return FieldifyReceiveFunc(p)
	}
	return logrus.Fields{
		"module":       "telebot",
		"action":       "receive",
		"endpoint":     p.FormattedEp,
		"chat_id":      p.ChatID,
		"chat_name":    p.ChatName,
		"message_id":   p.MessageID,
		"message_time": p.MessageTime.Format(time.RFC3339),
	}
}

// LogReceiveToLogrus logs a receive-event message to logrus.Logger using given endpoint and handler's received telebot.Message.
func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, received *telebot.Message, options ...LoggerOption) {
	if logger == nil || endpoint == nil || received == nil || received.Chat == nil {
		return
	}
	p := extractReceiveLoggerParam(endpoint, received)
	if p == nil {
		return
	}
	m := formatReceiveLoggerParam(p)
	f := fieldifyReceiveLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	logger.WithFields(f).Info(m)
}

// LogReceiveToLogger logs a receive-event message to logrus.StdLogger using given endpoint and handler's received telebot.Message.
func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, received *telebot.Message, options ...LoggerOption) {
	if logger == nil || endpoint == nil || received == nil || received.Chat == nil {
		return
	}
	p := extractReceiveLoggerParam(endpoint, received)
	if p == nil {
		return
	}
	m := formatReceiveLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// =======
// respond
// =======

// RespondLoggerParam stores some respond-event (RespondEvent) logger parameters and is used by LogRespondToLogrus and LogRespondToLogger.
type RespondLoggerParam struct {
	// origin
	EventType     RespondEventType
	Event         *RespondEvent
	SourceChat    *telebot.Chat    // fixed
	SourceMessage *telebot.Message // fixed (not send)
	ResultMessage *telebot.Message // fixed (no error)
	ReturnedError error

	// field
	SourceChatID       int64
	SourceChatName     string
	SourceMessageID    int
	SourceMessageTime  time.Time
	ResultMessageID    int
	ResultMessageChars int
	ResultMessageTime  time.Time
	ReplyLatency       *time.Duration // only for reply
	ReturnedErrorMsg   string
}

var (
	// FormatRespondFunc is a custom RespondLoggerParam's format function for LogRespondToLogrus and LogRespondToLogger.
	FormatRespondFunc func(p *RespondLoggerParam) string

	// FieldifyRespondFunc is a custom RespondLoggerParam's fieldify function for LogRespondToLogrus.
	FieldifyRespondFunc func(p *RespondLoggerParam) logrus.Fields
)

// extractRespondLoggerParam extracts and returns RespondLoggerParam using given parameters.
func extractRespondLoggerParam(typ RespondEventType, ev *RespondEvent) *RespondLoggerParam {
	p := &RespondLoggerParam{EventType: typ, Event: ev}
	var sc *telebot.Chat
	var sm, rm *telebot.Message

	switch typ {
	case RespondSendEvent:
		sc = ev.SendSource
		rm = ev.SendResult
	case RespondReplyEvent:
		sc = ev.ReplySource.Chat
		sm = ev.ReplySource
		rm = ev.ReplyResult
	case RespondEditEvent:
		sc = ev.EditSource.Chat
		sm = ev.EditSource
		rm = ev.EditResult
	case RespondDeleteEvent:
		sc = ev.DeleteSource.Chat
		sm = ev.DeleteSource
		rm = ev.DeleteResult
	default:
		return nil
	}

	if sc != nil {
		p.SourceChat = sc
		p.SourceChatID = sc.ID
		p.SourceChatName = sc.Username
	}
	if sm != nil {
		p.SourceMessage = sm
		p.SourceMessageID = sm.ID
		p.SourceMessageTime = sm.Time()
	}
	if rm != nil {
		p.ResultMessage = rm
		p.ResultMessageID = rm.ID
		p.ResultMessageChars = len([]rune(rm.Text))
		p.ResultMessageTime = rm.Time()
	}
	if typ == RespondReplyEvent && sm != nil && rm != nil {
		latency := rm.Time().Sub(sm.Time())
		p.ReplyLatency = &latency
	}
	if err := ev.ReturnedError; err != nil {
		p.ReturnedError = err
		p.ReturnedErrorMsg = err.Error()
	}
	return p
}

// formatRespondLoggerParam formats given RespondLoggerParam to string for LogRespondToLogrus and LogRespondToLogger.
//
// The default format logs like:
// 	RespondReplyEvent:
// 	[Telebot]  rep | msg#    3348 |   4096 chr | rep_to#    3346 |      993.3µs | 12345678 Aoi-hosizora
// 	         |----|      |-------| |------|      |-------| |------------| |---------------------|
// 	           4             7        6              7           12                  ...
// 	RespondSendEvent, RespondEditEvent, RespondDeleteEvent:
// 	[Telebot] send | msg#    3348 |   4096 chr | 12345678 Aoi-hosizora
// 	[Telebot] edit | msg#    3348 |   4096 chr | 12345678 Aoi-hosizora
// 	[Telebot]  del | msg#    3348 |   4096 chr | 12345678 Aoi-hosizora
// 	         |----|      |-------| |------|   |---------------------|
// 	           4             7        6                  ...
// 	Error cases:
// 	[Telebot]  rep | 12345678 Aoi-hosizora | err: telegram: bot was blocked by the user (401)
// 	[Telebot] send | 12345678 Aoi-hosizora | err: test error
// 	         |----| |---------------------|
// 	           4              ...
func formatRespondLoggerParam(p *RespondLoggerParam) string {
	if FormatRespondFunc != nil {
		return FormatRespondFunc(p)
	}
	if p.ReturnedErrorMsg != "" {
		return fmt.Sprintf("[Telebot] %4s | %d %s | err: %s", string(p.EventType), p.SourceChatID, p.SourceChatName, p.ReturnedErrorMsg)
	}
	switch p.EventType {
	case RespondReplyEvent:
		return fmt.Sprintf("[Telebot]  rep | msg# %7s | %6s chr | rep_to# %7s | %12s | %d %s",
			xnumber.Itoa(p.ResultMessageID), xnumber.Itoa(p.ResultMessageChars), xnumber.Itoa(p.SourceMessageID), p.ReplyLatency.String(), p.SourceChatID, p.SourceChatName)
	case RespondSendEvent, RespondEditEvent, RespondDeleteEvent:
		return fmt.Sprintf("[Telebot] %4s | msg# %7s | %6s chr | %d %s", string(p.EventType),
			xnumber.Itoa(p.ResultMessageID), xnumber.Itoa(p.ResultMessageChars), p.SourceChatID, p.SourceChatName)
	default:
		return ""
	}
}

// fieldifyRespondLoggerParam fieldifies given ReplyLoggerParam to logrus.Fields.
//
// The default contains the following fields: module, action, sent_id, sent_type, sent_time, chat_id, chat_name.
func fieldifyRespondLoggerParam(p *RespondLoggerParam) logrus.Fields {
	if FieldifyRespondFunc != nil {
		return FieldifyRespondFunc(p)
	}
	l := logrus.Fields{"module": "telebot", "action": "respond_" + string(p.EventType)}
	if p.SourceChat != nil {
		l["source_chat_id"] = p.SourceChatID
		l["source_chat_name"] = p.SourceChatName
	}
	if p.SourceMessage != nil {
		l["source_message_id"] = p.SourceMessageID
		l["source_message_time"] = p.SourceMessageTime.Format(time.RFC3339)
	}
	if p.ResultMessage != nil {
		l["result_message_id"] = p.ResultMessageID
		l["result_message_chars"] = p.ResultMessageChars
		l["result_message_time"] = p.ResultMessageTime.Format(time.RFC3339)
	}
	if p.ReplyLatency != nil {
		l["reply_latency"] = p.ReplyLatency
	}
	if p.ReturnedErrorMsg != "" {
		l["returned_error_msg"] = p.ReturnedErrorMsg
	}
	return l
}

// LogRespondToLogrus logs a respond-event (RespondEvent) message to logrus.Logger using xxx.
func LogRespondToLogrus(logger *logrus.Logger, typ RespondEventType, ev *RespondEvent, options ...LoggerOption) {
	if logger == nil || ev == nil {
		return
	}
	p := extractRespondLoggerParam(typ, ev)
	if p == nil {
		return
	}
	m := formatRespondLoggerParam(p)
	f := fieldifyRespondLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	extra.ApplyToFields(f)
	logger.WithFields(f).Info(m)
}

// LogRespondToLogger logs a respond-event (RespondEvent) message to logrus.StdLogger using xxx.
func LogRespondToLogger(logger logrus.StdLogger, typ RespondEventType, ev *RespondEvent, options ...LoggerOption) {
	if logger == nil || ev == nil {
		return
	}
	p := extractRespondLoggerParam(typ, ev)
	if p == nil {
		return
	}
	m := formatRespondLoggerParam(p)
	extra := internal.BuildLoggerOptions(options)
	extra.ApplyToMessage(&m)
	logger.Print(m)
}

// ========
// internal
// ========

// handlerFuncName returns the name of handler function, used for handledCallback.
func handlerFuncName(handler interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
}

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
		unique := ep.Text // CallbackUnique
		if unique == "" {
			return "", false // empty
		}
		return "$rep_btn:" + unique, true
	case *telebot.InlineButton:
		unique := ep.Unique // CallbackUnique - '\f'
		if unique == "" {
			return "", false // empty
		}
		return "$inl_btn:" + unique, true
	default:
		return "", false // unsupported endpoint
	}
}
