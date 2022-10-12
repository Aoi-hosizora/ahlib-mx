package xtelebot

import (
	"context"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"reflect"
)

// ================
// bot wrapper type
// ================

// BotWrapper represents a telebot.Bot wrapper type with some custom handling and sending behaviors. For more details about telegram bot api,
// please visit https://core.telegram.org/bots/api.
type BotWrapper struct {
	bot  *telebot.Bot
	data *BotData
	shs  *StateHandlerSet

	handledCallback   func(endpoint interface{}, formattedEndpoint string, handlerName string)
	receivedCallback  func(endpoint interface{}, received *telebot.Message)
	respondedCallback func(typ RespondEventType, event *RespondEvent)
	panicHandler      func(endpoint, messageOrCallback, value interface{})
}

const (
	panicNilTelebot = "xtelebot: nil telebot"
)

// NewBotWrapper creates a new BotWrapper with given telebot.Bot and new BotData and StateHandlerSet, panics when using nil telebot.Bot.
func NewBotWrapper(bot *telebot.Bot) *BotWrapper {
	if bot == nil {
		panic(panicNilTelebot)
	}
	return &BotWrapper{
		bot:  bot,
		data: NewBotData(),
		shs:  NewStateHandlerSet(),

		handledCallback:   DefaultHandledCallback,
		receivedCallback:  nil, // defaults to do nothing
		respondedCallback: nil, // defaults to do nothing
		panicHandler: func(_, _, value interface{}) {
			log.Printf("xtelebot warning: Handler panicked with `%v`", value)
		},
	}
}

// Bot returns the inner telebot.Bot from BotWrapper.
func (b *BotWrapper) Bot() *telebot.Bot {
	return b.bot
}

// Data returns the inner BotData from BotWrapper.
func (b *BotWrapper) Data() *BotData {
	return b.data
}

// Shs returns the inner StateHandlerSet from BotWrapper.
func (b *BotWrapper) Shs() *StateHandlerSet {
	return b.shs
}

// ==================
// bot wrapper handle
// ==================

type (
	// MessageHandler represents a handler type for string command and telebot.ReplyButton.
	MessageHandler func(*BotWrapper, *telebot.Message)

	// CallbackHandler represents a handler type for telebot.InlineButton.
	CallbackHandler func(*BotWrapper, *telebot.Callback)
)

const (
	panicInvalidEndpoint   = "xtelebot: invalid endpoint"
	panicInvalidCommand    = "xtelebot: invalid command"
	panicNilButton         = "xtelebot: nil button"
	panicEmptyButtonText   = "xtelebot: empty button text"
	panicEmptyButtonUnique = "xtelebot: empty button unique"
)

// IsHandled checks whether given endpoint's handler has been handled or registered, panics when using invalid endpoint, that is neither string
// nor telebot.CallbackEndpoint.
func (b *BotWrapper) IsHandled(endpoint interface{}) bool {
	handlerMap := xreflect.GetUnexportedField(xreflect.FieldValueOf(b.bot, "handlers"))
	switch ep := endpoint.(type) {
	case string:
		return handlerMap.MapIndex(reflect.ValueOf(ep)).IsValid()
	case telebot.CallbackEndpoint:
		return handlerMap.MapIndex(reflect.ValueOf(ep.CallbackUnique())).IsValid()
	default:
		panic(panicInvalidEndpoint)
	}
}

// RemoveHandler removes the handler of given endpoint, panics when using invalid endpoint, that is neither string nor telebot.CallbackEndpoint.
func (b *BotWrapper) RemoveHandler(endpoint interface{}) {
	handlerMap := xreflect.GetUnexportedField(xreflect.FieldValueOf(b.bot, "handlers"))
	switch ep := endpoint.(type) {
	case string:
		handlerMap.SetMapIndex(reflect.ValueOf(ep), reflect.Value{})
	case telebot.CallbackEndpoint:
		handlerMap.SetMapIndex(reflect.ValueOf(ep.CallbackUnique()), reflect.Value{})
	default:
		panic(panicInvalidEndpoint)
	}
}

// HandleCommand handles string command with MessageHandler to telebot.Bot, panics when using invalid command or nil handler, visit
// https://github.com/tucnak/telebot/tree/v2#commands for more details.
func (b *BotWrapper) HandleCommand(command string, handler MessageHandler) {
	if !isCommandValid(command) {
		panic(panicInvalidCommand)
	}
	if handler == nil {
		panic(panicNilHandler)
	}

	b.bot.Handle(command, func(m *telebot.Message) {
		defer func() {
			v := recover()
			if v != nil && b.panicHandler != nil {
				b.panicHandler(command, m, v)
			}
		}()
		if b.receivedCallback != nil {
			b.receivedCallback(command, m)
		}
		handler(b, m)
	})

	formatted, _ := formatEndpoint(command)
	if b.handledCallback != nil {
		b.handledCallback(command, formatted, handlerFuncName(handler))
	}
}

// HandleReplyButton handles telebot.ReplyButton with MessageHandler to telebot.Bot, panics when using nil button, invalid button or
// nil handler, visit https://github.com/tucnak/telebot/tree/v2#keyboards for more details.
func (b *BotWrapper) HandleReplyButton(button *telebot.ReplyButton, handler MessageHandler) {
	if button == nil {
		panic(panicNilButton)
	}
	if button.Text /* CallbackUnique => Text */ == "" {
		panic(panicEmptyButtonText)
	}
	if handler == nil {
		panic(panicNilHandler)
	}

	b.bot.Handle(button, func(m *telebot.Message) {
		defer func() {
			v := recover()
			if v != nil && b.panicHandler != nil {
				b.panicHandler(button, m, v)
			}
		}()
		if b.receivedCallback != nil {
			b.receivedCallback(button, m)
		}
		handler(b, m)
	})

	formatted, _ := formatEndpoint(button)
	if b.handledCallback != nil {
		b.handledCallback(button, formatted, handlerFuncName(handler))
	}
}

// HandleInlineButton handles telebot.InlineButton with CallbackHandler to telebot.Bot, panics when using nil button, invalid button or
// nil handler, visit https://github.com/tucnak/telebot/tree/v2#keyboards for more details.
func (b *BotWrapper) HandleInlineButton(button *telebot.InlineButton, handler CallbackHandler) {
	if button == nil {
		panic(panicNilButton)
	}
	if button.Unique /* CallbackUnique => \fUnique */ == "" {
		panic(panicEmptyButtonUnique)
	}
	if handler == nil {
		panic(panicNilHandler)
	}

	b.bot.Handle(button, func(c *telebot.Callback) {
		defer func() {
			v := recover()
			if v != nil && b.panicHandler != nil {
				b.panicHandler(button, c, v)
			}
		}()
		if b.receivedCallback != nil {
			b.receivedCallback(button, c.Message)
		}
		handler(b, c)
	})

	formatted, _ := formatEndpoint(button)
	if b.handledCallback != nil {
		b.handledCallback(button, formatted, handlerFuncName(handler))
	}
}

// ===================
// bot wrapper respond
// ===================

// RespondEventType is a type of respond event type (such as "send", "reply", "edit", "delete" and "callback"), will be used in respondedCallback,
// LogRespondToLogrus and LogRespondToLogger.
type RespondEventType string

const (
	RespondSendEvent     RespondEventType = "send"     // RespondEventType for BotWrapper.RespondSend.
	RespondReplyEvent    RespondEventType = "reply"    // RespondEventType for BotWrapper.RespondReply.
	RespondEditEvent     RespondEventType = "edit"     // RespondEventType for BotWrapper.RespondEdit.
	RespondDeleteEvent   RespondEventType = "delete"   // RespondEventType for BotWrapper.RespondDelete.
	RespondCallbackEvent RespondEventType = "callback" // RespondEventType for BotWrapper.RespondCallback.
)

// RespondEvent is a type of respond event, containing arguments of respond method (such as BotWrapper.RespondSend) and responded result and error,
// this will be used in responded callback, LogRespondToLogrus and LogRespondToLogger.
type RespondEvent struct {
	// for RespondSendEvent
	SendSource  *telebot.Chat
	SendWhat    interface{}
	SendOptions []interface{}
	SendResult  *telebot.Message

	// for RespondReplyEvent
	ReplySource   *telebot.Message
	ReplyExplicit bool
	ReplyWhat     interface{}
	ReplyOptions  []interface{}
	ReplyResult   *telebot.Message

	// for RespondEditEvent
	EditSource  *telebot.Message
	EditWhat    interface{}
	EditOptions []interface{}
	EditResult  *telebot.Message

	// for RespondDeleteEvent
	DeleteSource *telebot.Message
	DeleteResult *telebot.Message // fake

	// for RespondCallbackEvent
	CallbackSource *telebot.Callback
	CallbackAnswer *telebot.CallbackResponse
	CallbackResult *telebot.CallbackResponse // fake

	// ctx and error
	RespondContext context.Context
	ReturnedError  error
}

var (
	errNilChat     = errors.New("xtelebot: nil source chat")
	errNilMessage  = errors.New("xtelebot: nil source message")
	errNilCallback = errors.New("xtelebot: nil source callback")
	errNilWhat     = errors.New("xtelebot: nil respond value")
)

// RespondSend is a convenient form of RespondSendCtx without context.Context.
func (b *BotWrapper) RespondSend(source *telebot.Chat, what interface{}, options ...interface{}) (*telebot.Message, error) {
	return b.RespondSendCtx(context.Background(), source, what, options...)
}

// RespondSendCtx responds and sends message to given telebot.Chat, if error returned is not caused by arguments, it will also invoke responded callback.
func (b *BotWrapper) RespondSendCtx(ctx context.Context, source *telebot.Chat, what interface{}, options ...interface{}) (*telebot.Message, error) {
	if source == nil {
		return nil, errNilChat
	}
	if what == nil {
		return nil, errNilWhat
	}

	msg, err := b.bot.Send(source, what, options...)
	if err == telebot.ErrUnsupportedWhat {
		return nil, err // no need to invoke respondedCallback
	}

	if b.respondedCallback != nil {
		b.respondedCallback(RespondSendEvent, &RespondEvent{
			SendSource: source, SendWhat: what, SendOptions: options,
			SendResult: msg, RespondContext: ctx, ReturnedError: err,
		})
	}
	return msg, err
}

// RespondReply is a convenient form of RespondReplyCtx without context.Context.
func (b *BotWrapper) RespondReply(source *telebot.Message, explicit bool, what interface{}, options ...interface{}) (*telebot.Message, error) {
	return b.RespondReplyCtx(context.Background(), source, explicit, what, options...)
}

// RespondReplyCtx responds and replies message to given telebot.Message explicitly or implicitly, if error returned is not caused by arguments, it will
// also invoke responded callback.
func (b *BotWrapper) RespondReplyCtx(ctx context.Context, source *telebot.Message, explicit bool, what interface{}, options ...interface{}) (*telebot.Message, error) {
	if source == nil {
		return nil, errNilMessage
	}
	if what == nil {
		return nil, errNilWhat
	}

	var msg *telebot.Message
	var err error
	if !explicit {
		msg, err = b.bot.Send(source.Chat, what, options...)
	} else {
		msg, err = b.bot.Reply(source, what, options...) // send with ReplyTo option
	}
	if err == telebot.ErrUnsupportedWhat {
		return nil, err
	}

	if b.respondedCallback != nil {
		b.respondedCallback(RespondReplyEvent, &RespondEvent{
			ReplySource: source, ReplyExplicit: explicit, ReplyWhat: what, ReplyOptions: options,
			ReplyResult: msg, RespondContext: ctx, ReturnedError: err,
		})
	}
	return msg, err
}

// RespondEdit is a convenient form of RespondEditCtx without context.Context.
func (b *BotWrapper) RespondEdit(source *telebot.Message, what interface{}, options ...interface{}) (*telebot.Message, error) {
	return b.RespondEditCtx(context.Background(), source, what, options...)
}

// RespondEditCtx responds and edits given telebot.Message with value, if error returned is not caused by arguments, it will also invoke responded callback.
func (b *BotWrapper) RespondEditCtx(ctx context.Context, source *telebot.Message, what interface{}, options ...interface{}) (*telebot.Message, error) {
	if source == nil {
		return nil, errNilMessage
	}
	if what == nil {
		return nil, errNilWhat
	}

	msg, err := b.bot.Edit(source, what, options...)
	if err == telebot.ErrUnsupportedWhat {
		return nil, err
	}

	if b.respondedCallback != nil {
		b.respondedCallback(RespondEditEvent, &RespondEvent{
			EditSource: source, EditWhat: what, EditOptions: options,
			EditResult: msg, RespondContext: ctx, ReturnedError: err,
		})
	}
	return msg, err
}

// RespondDelete is a convenient form of RespondDeleteCtx without context.Context.
func (b *BotWrapper) RespondDelete(source *telebot.Message) error {
	return b.RespondDeleteCtx(context.Background(), source)
}

// RespondDeleteCtx responds and deletes given telebot.Message, if error returned is not caused by arguments, it will also invoke responded callback.
func (b *BotWrapper) RespondDeleteCtx(ctx context.Context, source *telebot.Message) error {
	if source == nil {
		return errNilMessage
	}

	err := b.bot.Delete(source)

	if b.respondedCallback != nil {
		b.respondedCallback(RespondDeleteEvent, &RespondEvent{
			DeleteSource: source,
			DeleteResult: source /* for unifying only */, RespondContext: ctx, ReturnedError: err,
		})
	}
	return err
}

// RespondCallback is a convenient form of RespondCallbackCtx without context.Context.
func (b *BotWrapper) RespondCallback(source *telebot.Callback, answer *telebot.CallbackResponse) error {
	return b.RespondCallbackCtx(context.Background(), source, answer)
}

// RespondCallbackCtx responds and answers to given telebot.Callback, if error returned is not caused by arguments, it will also invoke responded callback.
func (b *BotWrapper) RespondCallbackCtx(ctx context.Context, source *telebot.Callback, answer *telebot.CallbackResponse) error {
	if source == nil {
		return errNilCallback
	}

	var err error
	if answer == nil {
		err = b.bot.Respond(source)
	} else {
		err = b.bot.Respond(source, answer)
	}

	if b.respondedCallback != nil {
		b.respondedCallback(RespondCallbackEvent, &RespondEvent{
			CallbackSource: source, CallbackAnswer: answer,
			CallbackResult: answer /* for unifying only */, RespondContext: ctx, ReturnedError: err,
		})
	}
	return err
}

// ====================
// bot wrapper callback
// ====================

// DefaultHandledCallback is the default BotWrapper's handledCallback, can be modified by BotWrapper.SetHandledCallback.
//
// The default callback logs like (just like gin.DebugPrintRouteFunc):
// 	[Telebot] /test-endpoint               --> ...
// 	[Telebot] $on_text                     --> ...
// 	[Telebot] $rep_btn:button_text         --> ...
// 	[Telebot] $inl_btn:button_unique       --> ...
// 	         |----------------------------|   |---|
// 	                       28                  ...
func DefaultHandledCallback(_ interface{}, formattedEndpoint string, handlerName string) {
	fmt.Printf("[Telebot] %-28s --> %s\n", formattedEndpoint, handlerName)
}

// DefaultColorizedHandledCallback is the DefaultAddedCallback (BotWrapper's handledCallback) in color.
//
// The default callback logs like (just like gin.DebugPrintRouteFunc):
// 	[Telebot] /test-endpoint               --> ...
// 	[Telebot] $on_text                     --> ...
// 	[Telebot] $rep_btn:button_text         --> ...
// 	[Telebot] $inl_btn:button_unique       --> ...
// 	         |----------------------------|   |---|
// 	                   28 (blue)               ...
func DefaultColorizedHandledCallback(_ interface{}, formattedEndpoint string, handlerName string) {
	fmt.Printf("[Telebot] %s --> %s\n", xcolor.Blue.Sprintf(fmt.Sprintf("%-28s", formattedEndpoint)), handlerName)
}

// SetHandledCallback sets endpoint handled callback, callback will be invoked in handling methods, defaults to DefaultHandledCallback.
func (b *BotWrapper) SetHandledCallback(f func(endpoint interface{}, formattedEndpoint string, handlerName string)) {
	b.handledCallback = f
}

// SetReceivedCallback sets received callback, callback will be invoked after receiving consumed messages which has been handled, defaults to do nothing.
func (b *BotWrapper) SetReceivedCallback(cb func(endpoint interface{}, received *telebot.Message)) {
	b.receivedCallback = cb
}

// SetRespondedCallback sets responded callback, callback will be invoked in respond methods, defaults to do nothing.
func (b *BotWrapper) SetRespondedCallback(cb func(typ RespondEventType, event *RespondEvent)) {
	b.respondedCallback = cb
}

// SetPanicHandler sets panic handler to all handlers. Note that the `messageOrCallback` parameter means handler's parameter, that is telebot.Message
// for MessageHandler and telebot.Callback for CallbackHandler, defaults to print warning message with given panicked value.
func (b *BotWrapper) SetPanicHandler(handler func(endpoint, messageOrCallback, value interface{})) {
	b.panicHandler = handler
}
