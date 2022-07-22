package xtelebot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// LoggerOption represents an option type for some logger functions' option, can be created by WithXXX functions.
type LoggerOption = internal.LoggerOption

// WithExtraText creates a LoggerOption to specify extra text logging in "...extra_text" style. Note that if you use this multiple times, only the last one will be retained.
func WithExtraText(text string) LoggerOption {
	return internal.WithExtraText(text)
}

// WithExtraFields creates a LoggerOption to specify logging with extra fields. Note that if you use this multiple times, only the last one will be retained.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return internal.WithExtraFields(fields)
}

// WithExtraFieldsV creates a LoggerOption to specify logging with extra fields in variadic. Note that if you use this multiple times, only the last one will be retained.
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
		Chat:     received.Chat, // non-nillable
		Message:  received,      // non-nillable

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
// 	         |----|      |-------| |----------------------------| |---------------------|
// 	           4             7                   28                         ...
func formatReceiveLoggerParam(p *ReceiveLoggerParam) string {
	if FormatReceiveFunc != nil {
		return FormatReceiveFunc(p)
	}
	return fmt.Sprintf("[Telebot] %s | msg# %7s | %28s | %d %s", colorizeEventType(""), // "recv"
		xnumber.Itoa(p.MessageID), p.FormattedEp, p.ChatID, p.ChatName)
}

// fieldifyReceiveLoggerParam fieldifies given ReceiveLoggerParam to logrus.Fields for LogReceiveToLogrus.
//
// The default contains the following fields: module, action, endpoint, chat_id, chat_name, message_id, message_time.
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

// LogReceiveToLogrus logs a receive-event message to logrus.Logger using given endpoint and telebot.Message received from handler.
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

// LogReceiveToLogger logs a receive-event message to logrus.StdLogger using given endpoint and telebot.Message received from handler.
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
	EventType RespondEventType
	Event     *RespondEvent

	// extracted origin
	SourceChat     *telebot.Chat             // sc: send, rep, edit, del, call
	SourceMessage  *telebot.Message          // sm:       rep, edit, del, call
	SourceCallback *telebot.Callback         // sl:                       call
	ResultMessage  *telebot.Message          // rm: send, rep, edit, del
	ResultAnswer   *telebot.CallbackResponse // ra:                       call
	ReturnedError  error                     // -

	// field
	SourceChatID       int64          // sc: send, rep, edit, del, call
	SourceChatName     string         // sc: send, rep, edit, del, call
	SourceMessageID    int            // sm:       rep, edit, del, call
	SourceMessageTime  time.Time      // sm:       rep, edit, del, call
	SourceCallbackID   string         // sl:                       call
	ResultMessageID    int            // rm: send, rep, edit, del
	ResultMessageChars int            // rm: send, rep, edit, del
	ResultMessageTime  time.Time      // rm: send, rep, edit, del
	ReplyLatency       *time.Duration // *:        rep
	CallbackAlert      *string        // *:                        call
	ReturnedErrorMsg   string         // -
}

var (
	// FormatRespondFunc is a custom RespondLoggerParam's format function for LogRespondToLogrus and LogRespondToLogger.
	FormatRespondFunc func(p *RespondLoggerParam) string

	// FieldifyRespondFunc is a custom RespondLoggerParam's fieldify function for LogRespondToLogrus.
	FieldifyRespondFunc func(p *RespondLoggerParam) logrus.Fields
)

// extractRespondLoggerParam extracts and returns RespondLoggerParam using given parameters, this unexported never panic.
func extractRespondLoggerParam(typ RespondEventType, ev *RespondEvent) *RespondLoggerParam {
	p := &RespondLoggerParam{EventType: typ, Event: ev}
	var sc *telebot.Chat             // source
	var sm *telebot.Message          // source
	var sl *telebot.Callback         // source
	var rm *telebot.Message          // result
	var ra *telebot.CallbackResponse // result

	switch typ {
	case RespondSendEvent:
		sc = ev.SendSource
		rm = ev.SendResult
	case RespondReplyEvent:
		sm = ev.ReplySource
		if ev.ReplySource != nil {
			sc = ev.ReplySource.Chat
		}
		rm = ev.ReplyResult
	case RespondEditEvent:
		sm = ev.EditSource
		if ev.EditSource != nil {
			sc = ev.EditSource.Chat
		}
		rm = ev.EditResult
	case RespondDeleteEvent:
		sm = ev.DeleteSource
		if ev.DeleteSource != nil {
			sc = ev.DeleteSource.Chat
		}
		rm = ev.DeleteResult
	case RespondCallbackEvent:
		sl = ev.CallbackSource
		if ev.CallbackSource != nil {
			sm = ev.CallbackSource.Message
			if ev.CallbackSource.Message != nil {
				sc = ev.CallbackSource.Message.Chat
			}
		}
		ra = ev.CallbackResult
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
	if sl != nil {
		p.SourceCallback = sl
		p.SourceCallbackID = sl.ID
	}
	if rm != nil {
		p.ResultMessage = rm
		p.ResultMessageID = rm.ID
		p.ResultMessageChars = len([]rune(rm.Text))
		p.ResultMessageTime = rm.Time()
	}
	if ra != nil {
		p.ResultAnswer = ra
	}
	if typ == RespondReplyEvent {
		latency := time.Duration(-1)
		if sm != nil && rm != nil {
			latency = rm.Time().Sub(sm.Time())
		}
		p.ReplyLatency = &latency
	}
	if typ == RespondCallbackEvent {
		s := "-"
		if ra != nil && strings.TrimSpace(ra.Text) != "" {
			if ra.ShowAlert {
				s = "with_alert"
			} else {
				s = "with_text"
			}
		}
		p.CallbackAlert = &s
	}
	if err := ev.ReturnedError; err != nil {
		p.ReturnedError = err
		p.ReturnedErrorMsg = err.Error()
	}
	return p
}

// colorizeEventType colorizes and truncates to 4 characters given RespondEventType to string.
func colorizeEventType(typ RespondEventType) string {
	switch typ {
	case "": // trick for receive log
		return xcolor.Blue.Sprintf("recv")
	case RespondSendEvent:
		return xcolor.Green.Sprint("send")
	case RespondReplyEvent:
		return xcolor.Green.Sprint(" rep")
	case RespondEditEvent:
		return xcolor.Yellow.Sprint("edit")
	case RespondDeleteEvent:
		return xcolor.Red.Sprint(" del")
	case RespondCallbackEvent:
		return xcolor.Cyan.Sprint("call")
	default:
		return " ???"
	}
}

// formatRespondLoggerParam formats given RespondLoggerParam to string for LogRespondToLogrus and LogRespondToLogger.
//
// The default format logs like:
// 	RespondReplyEvent:
// 	[Telebot]  rep | msg#    3348 |   4096 chr | rep_to#    3346 |      993.3µs | 12345678 Aoi-hosizora
// 	         |----|      |-------| |------|             |-------| |------------| |---------------------|
// 	           4             7        6                     7           12                  ...
// 	RespondSendEvent, RespondEditEvent, RespondDeleteEvent:
// 	[Telebot] send | msg#    3348 |   4096 chr | 12345678 Aoi-hosizora
// 	[Telebot] edit | msg#    3348 |   4096 chr | 12345678 Aoi-hosizora
// 	[Telebot]  del | msg#    3348 |   4096 chr | 12345678 Aoi-hosizora
// 	         |----|      |-------| |------|     |---------------------|
// 	           4             7        6                    ...
// 	RespondCallbackEvent:
// 	[Telebot] call | msg#    3348 | with_alert | 12345678 Aoi-hosizora
// 	[Telebot] call | msg#    3348 |  with_text | 12345678 Aoi-hosizora
// 	         |----|      |-------| |----------| |---------------------|
// 	           4             7          10                ...
// 	Error cases:
// 	[Telebot]  rep | msg#    3348 | 12345678 Aoi-hosizora | err: telegram: bot was blocked by the user (401)
// 	[Telebot] edit | msg#    3348 | 12345678 Aoi-hosizora | err: test error
// 	[Telebot] send | msg#       0 | 12345678 Aoi-hosizora | err: test error
// 	         |----|      |-------| |---------------------|
// 	           4             7               ...
func formatRespondLoggerParam(p *RespondLoggerParam) string {
	if FormatRespondFunc != nil {
		return FormatRespondFunc(p)
	}
	if p.ReturnedErrorMsg != "" {
		return fmt.Sprintf("[Telebot] %s | msg# %7s | %d %s | err: %s", colorizeEventType(p.EventType),
			xnumber.Itoa(p.ResultMessageID), p.SourceChatID, p.SourceChatName, p.ReturnedErrorMsg)
	}
	switch p.EventType {
	case RespondReplyEvent:
		return fmt.Sprintf("[Telebot] %s | msg# %7s | %6s chr | rep_to# %7s | %12s | %d %s", colorizeEventType(p.EventType), // " rep"
			xnumber.Itoa(p.ResultMessageID), xnumber.Itoa(p.ResultMessageChars), xnumber.Itoa(p.SourceMessageID), p.ReplyLatency.String(), p.SourceChatID, p.SourceChatName)
	case RespondSendEvent, RespondEditEvent, RespondDeleteEvent:
		return fmt.Sprintf("[Telebot] %s | msg# %7s | %6s chr | %d %s", colorizeEventType(p.EventType), // "send" / "edit" / " del"
			xnumber.Itoa(p.ResultMessageID), xnumber.Itoa(p.ResultMessageChars), p.SourceChatID, p.SourceChatName)
	case RespondCallbackEvent:
		return fmt.Sprintf("[Telebot] %s | msg# %7s | %10s | %d %s", colorizeEventType(p.EventType), // "call"
			xnumber.Itoa(p.SourceMessageID), *p.CallbackAlert, p.SourceChatID, p.SourceChatName)
	default:
		return ""
	}
}

// fieldifyRespondLoggerParam fieldifies given RespondLoggerParam to logrus.Fields.
//
// The default contains the following fields: module, action, source_chat_id, source_chat_name, source_message_id, source_message_time, source_callback_id,
// result_message_id, result_message_chars, result_message_time, reply_latency, callback_alert, returned_error_msg.
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
	if p.SourceCallback != nil {
		l["source_callback_id"] = p.SourceCallbackID
	}
	if p.ResultMessage != nil {
		l["result_message_id"] = p.ResultMessageID
		l["result_message_chars"] = p.ResultMessageChars
		l["result_message_time"] = p.ResultMessageTime.Format(time.RFC3339)
	}
	if p.ReplyLatency != nil {
		l["reply_latency"] = *p.ReplyLatency
	}
	if p.CallbackAlert != nil {
		l["callback_alert"] = *p.CallbackAlert
	}
	if p.ReturnedErrorMsg != "" {
		l["returned_error_msg"] = p.ReturnedErrorMsg
	}
	return l
}

// LogRespondToLogrus logs a respond-event message to logrus.Logger using given RespondEventType and RespondEvent.
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
	if p.ReturnedErrorMsg != "" {
		logger.WithFields(f).Error(m)
	} else {
		logger.WithFields(f).Info(m)
	}
}

// LogRespondToLogger logs a respond-event message to logrus.StdLogger using given RespondEventType and RespondEvent.
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

// isCommandValid checks whether given command is valid.
func isCommandValid(command string) bool {
	return len(command) > 1 && (command[0] == '/' || command[0] == '\a')
}

// handlerFuncName returns the name of handler function, used for handledCallback.
func handlerFuncName(handler interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
}

// formatEndpoint formats given endpoint type to string, only supports string, telebot.ReplyButton and telebot.InlineButton types.
func formatEndpoint(endpoint interface{}) (formatted string, supported bool) {
	switch ep := endpoint.(type) {
	case string:
		if !isCommandValid(ep) {
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
