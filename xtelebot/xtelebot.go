package xtelebot

import (
	"errors"
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"sync"
)

// ================
// bot data related
// ================

// ChatState is a type of chat states, used in BotData.
type ChatState uint64

// BotData represents a set of chats data in a telegram bot, including states and caches.
type BotData struct {
	initialState ChatState

	states map[int64]ChatState
	mus    sync.RWMutex
	caches map[int64]map[string]interface{}
	muc    sync.RWMutex
}

// NewBotData creates a default BotData.
func NewBotData() *BotData {
	return &BotData{
		initialState: 0,

		states: make(map[int64]ChatState),
		caches: make(map[int64]map[string]interface{}),
	}
}

// ==============
// bot data state
// ==============

// GetStateChats returns all ids from chats which has been set state, the returned slice has no order.
func (b *BotData) GetStateChats() []int64 {
	b.mus.RLock()
	ids := make([]int64, 0, len(b.states))
	for key := range b.states {
		ids = append(ids, key)
	}
	b.mus.RUnlock()
	return ids
}

// GetState returns a chat's state, returns false if no state is set.
func (b *BotData) GetState(chatID int64) (ChatState, bool) {
	b.mus.RLock()
	s, ok := b.states[chatID]
	b.mus.RUnlock()
	return s, ok
}

// GetStateOr returns a chat's state, returns the fallback state if no state is set.
func (b *BotData) GetStateOr(chatID int64, fallbackState ChatState) ChatState {
	s, ok := b.GetState(chatID)
	if !ok {
		return fallbackState
	}
	return s
}

// SetInitialState sets initial ChatState to BotData.
func (b *BotData) SetInitialState(s ChatState) {
	b.mus.Lock()
	b.initialState = s
	b.mus.Unlock()
}

// GetStateOrInit returns a chat's state, sets to the initial state and returns it if no state is set.
func (b *BotData) GetStateOrInit(chatID int64) ChatState {
	b.mus.Lock()
	s, ok := b.states[chatID]
	if !ok {
		s = b.initialState
		b.states[chatID] = s
	}
	b.mus.Unlock()
	return s
}

// SetState sets a chat's state.
func (b *BotData) SetState(chatID int64, state ChatState) {
	b.mus.Lock()
	b.states[chatID] = state
	b.mus.Unlock()
}

// ResetState resets a chat's state to the initial state.
func (b *BotData) ResetState(chatID int64) {
	b.mus.Lock()
	b.states[chatID] = b.initialState
	b.mus.Unlock()
}

// DeleteState deletes a chat's state.
func (b *BotData) DeleteState(chatID int64) {
	b.mus.Lock()
	delete(b.states, chatID)
	b.mus.Unlock()
}

// ==============
// bot data cache
// ==============

// GetCacheChats returns all ids from chats which has been set cache, the returned slice has no order.
func (b *BotData) GetCacheChats() []int64 {
	b.muc.RLock()
	ids := make([]int64, 0, len(b.caches))
	for key := range b.caches {
		ids = append(ids, key)
	}
	b.muc.RUnlock()
	return ids
}

// GetCache returns a chat's cache data, returns false if no cache is set or the key is not found.
func (b *BotData) GetCache(chatID int64, key string) (interface{}, bool) {
	b.muc.RLock()
	if m, ok := b.caches[chatID]; ok {
		if value, ok := m[key]; ok {
			b.muc.RUnlock()
			return value, true
		}
	}
	b.muc.RUnlock()
	return nil, false
}

// GetCacheOr returns a chat's cache data, returns fallback value if no cache is set or the key is not found.
func (b *BotData) GetCacheOr(chatID int64, key string, fallbackValue interface{}) interface{} {
	value, ok := b.GetCache(chatID, key)
	if !ok {
		return fallbackValue
	}
	return value
}

// GetChatCaches returns a chat's all caches data, returns false if no cache is set.
func (b *BotData) GetChatCaches(chatID int64) (map[string]interface{}, bool) {
	b.muc.RLock()
	m, ok := b.caches[chatID]
	if !ok {
		b.muc.RUnlock()
		return nil, false
	}
	out := make(map[string]interface{}) // shallow copy
	for k, v := range m {
		out[k] = v
	}
	b.muc.RUnlock()
	return out, true
}

// SetCache sets a chat's cache data using the given key and value.
func (b *BotData) SetCache(chatID int64, key string, value interface{}) {
	b.muc.Lock()
	m, ok := b.caches[chatID]
	if !ok {
		m = make(map[string]interface{})
		b.caches[chatID] = m
	}
	m[key] = value
	b.muc.Unlock()
}

// RemoveCache removes a chat's cache data.
func (b *BotData) RemoveCache(chatID int64, key string) {
	b.muc.Lock()
	if m, ok := b.caches[chatID]; ok {
		delete(m, key)
	}
	b.muc.Unlock()
}

// ClearCaches clears a chat's all caches.
func (b *BotData) ClearCaches(chatID int64) {
	b.muc.Lock()
	delete(b.caches, chatID)
	b.muc.Unlock()
}

// ===================
// bot wrapper related
// ===================

// BotWrapper represents a telebot.Bot wrapper type with some custom handling and sending behaviors.
type BotWrapper struct {
	bot  *telebot.Bot
	data *BotData

	handledEndpointCallback func(endpoint string, handlerName string)
	panicHandler            func(v interface{})
	receivedCallback        func(endpoint interface{}, received *telebot.Message)
	afterRepliedCallback    func(received *telebot.Message, replied *telebot.Message, err error)
	afterSentCallback       func(chat *telebot.Chat, sent *telebot.Message, err error)
}

const (
	panicNilBot = "xtelebot: nil telebot.Bot"
)

// NewBotWrapper creates a new BotWrapper with given telebot.Bot, panics when using nil telebot.Bot.
func NewBotWrapper(bot *telebot.Bot) *BotWrapper {
	if bot == nil {
		panic(panicNilBot)
	}
	return &BotWrapper{
		bot:  bot,
		data: NewBotData(),

		handledEndpointCallback: func(endpoint string, handlerName string) {
			fmt.Printf("[Bot-debug] %-30s --> %s\n", endpoint, handlerName)
		},
		panicHandler: func(v interface{}) { log.Printf("Warning: Panic with `%v`", v) },
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

type (
	// MessageHandler represents a handler type for string command and telebot.ReplyButton.
	MessageHandler func(*BotWrapper, *telebot.Message)

	// CallbackHandler represents a handler type for telebot.InlineButton.
	CallbackHandler func(*BotWrapper, *telebot.Callback)
)

const (
	panicInvalidCommand = "xtelebot: invalid command"
	panicNilHandler     = "xtelebot: nil handler"
	panicNilButton      = "xtelebot: nil button"
	panicEmptyUnique    = "xtelebot: empty button unique"
)

// HandleCommand adds string command and MessageHandler to telebot.Bot.
func (b *BotWrapper) HandleCommand(command string, handler MessageHandler) {
	if len(command) <= 1 || (command[0] != '/' && command[0] != '\a') {
		panic(panicInvalidCommand)
	}
	if handler == nil {
		panic(panicNilHandler)
	}
	b.bot.Handle(command, func(m *telebot.Message) {
		defer func() {
			if v := recover(); v != nil && b.panicHandler != nil {
				b.panicHandler(v)
			}
		}()
		if b.receivedCallback != nil {
			b.receivedCallback(command, m)
		}
		handler(b, m)
	})
	handledEndpointCallback(b.handledEndpointCallback, command, handler)
}

// HandleReplyButton adds telebot.ReplyButton and MessageHandler to telebot.Bot.
func (b *BotWrapper) HandleReplyButton(button *telebot.ReplyButton, handler MessageHandler) {
	if button == nil {
		panic(panicNilButton)
	}
	if button.CallbackUnique() /* Text */ == "" {
		panic(panicEmptyUnique)
	}
	if handler == nil {
		panic(panicNilHandler)
	}
	b.bot.Handle(button, func(m *telebot.Message) {
		defer func() {
			if v := recover(); v != nil && b.panicHandler != nil {
				b.panicHandler(v)
			}
		}()
		if b.receivedCallback != nil {
			b.receivedCallback(button, m)
		}
		handler(b, m)
	})
	handledEndpointCallback(b.handledEndpointCallback, button, handler)
}

// HandleInlineButton adds telebot.InlineButton and CallbackHandler to telebot.Bot.
func (b *BotWrapper) HandleInlineButton(button *telebot.InlineButton, handler CallbackHandler) {
	if button == nil {
		panic(panicNilButton)
	}
	if button.Unique == "" {
		panic(panicEmptyUnique)
	}
	if handler == nil {
		panic(panicNilHandler)
	}
	b.bot.Handle(button, func(c *telebot.Callback) {
		defer func() {
			if v := recover(); v != nil && b.panicHandler != nil {
				b.panicHandler(v)
			}
		}()
		if b.receivedCallback != nil {
			b.receivedCallback(button, c.Message)
		}
		handler(b, c)
	})
	handledEndpointCallback(b.handledEndpointCallback, button, handler)
}

var (
	errNilMessage = errors.New("xtelebot: nil message")
	errNilWhat    = errors.New("xtelebot: nil what")
	errNilChat    = errors.New("xtelebot: nil chat")
)

// ReplyTo sends message to chat from given telebot.Message (means replying to message), and invokes repliedCallback.
func (b *BotWrapper) ReplyTo(received *telebot.Message, what interface{}, options ...interface{}) (*telebot.Message, error) {
	if received == nil {
		return nil, errNilMessage
	}
	if what == nil {
		return nil, errNilWhat
	}
	msg, err := b.bot.Send(received.Chat, what, options...)
	if b.afterRepliedCallback != nil {
		b.afterRepliedCallback(received, msg, err)
	}
	return msg, err
}

// SendTo sends message to given telebot.Chat, and invokes sentCallback.
func (b *BotWrapper) SendTo(chat *telebot.Chat, what interface{}, options ...interface{}) (*telebot.Message, error) {
	if chat == nil {
		return nil, errNilChat
	}
	if what == nil {
		return nil, errNilWhat
	}
	msg, err := b.bot.Send(chat, what, options...)
	if b.afterSentCallback != nil {
		b.afterSentCallback(chat, msg, err)
	}
	return msg, err
}

// SetHandledEndpointCallback sets endpoint handled callback, defaults to print the debug message.
//
// The default handler logs like:
// 	[Bot-debug] /test-endpoint                 --> ...
// 	[Bot-debug] $on_text                       --> ...
// 	[Bot-debug] $rep_btn:reply_button          --> ...
// 	[Bot-debug] $inl_btn:inline_button         --> ...
// 	           |------------------------------|   |---|
// 	                          30                   ...
func (b *BotWrapper) SetHandledEndpointCallback(f func(endpoint string, handlerName string)) {
	b.handledEndpointCallback = f
}

// SetPanicHandler sets panic handler for all endpoint handlers.
func (b *BotWrapper) SetPanicHandler(handler func(v interface{})) {
	b.panicHandler = handler
}

// SetReceivedCallback sets received callback, that will be invoked in BotWrapper.HandleCommand, BotWrapper.HandleInlineButton and BotWrapper.HandleReplyButton.
func (b *BotWrapper) SetReceivedCallback(cb func(endpoint interface{}, received *telebot.Message)) {
	b.receivedCallback = cb
}

// SetAfterRepliedCallback sets replied callback, that will be invoked in BotWrapper.ReplyTo after telebot.Bot Send() invoked.
func (b *BotWrapper) SetAfterRepliedCallback(cb func(received *telebot.Message, replied *telebot.Message, err error)) {
	b.afterRepliedCallback = cb
}

// SetAfterSentCallback sets sent callback, that will be invoked in BotWrapper.SendTo after telebot.Bot Send() invoked.
func (b *BotWrapper) SetAfterSentCallback(cb func(chat *telebot.Chat, sent *telebot.Message, err error)) {
	b.afterSentCallback = cb
}
