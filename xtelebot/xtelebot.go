package xtelebot

import (
	"gopkg.in/tucnak/telebot.v2"
	"sync"
)

// ====================
// reply markup helpers
// ====================

// _markup is a global telebot.ReplyMarkup for telebot.InlineButton and telebot.ReplyButton helpers.
var _markup = &telebot.ReplyMarkup{}

// DataBtn creates a telebot.InlineButton using given button text, callback unique and callback data.
func DataBtn(text, unique string, data ...string) *telebot.InlineButton {
	return _markup.Data(text, unique, data...).Inline()
}

// TextBtn creates a telebot.ReplyButton using given button text.
func TextBtn(text string) *telebot.ReplyButton {
	return _markup.Text(text).Reply()
}

// URLBtn creates a telebot.InlineButton using given button text and url.
func URLBtn(text, url string) *telebot.InlineButton {
	return _markup.URL(text, url).Inline()
}

type (
	// InlineRow is a collection type represents a row of telebot.InlineButton, used in InlineKeyboard.
	InlineRow []*telebot.InlineButton

	// ReplyRow is a collection type represents a row of telebot.ReplyButton, used in ReplyKeyboard.
	ReplyRow []*telebot.ReplyButton
)

// InlineKeyboard creates a telebot.InlineButton keyboard with given InlineRow-s.
//
// Example:
// 	markup := &telebot.ReplyMarkup{
// 		InlineKeyboard: xtelebot.InlineKeyboard(
// 			xtelebot.InlineRow{button.InlineBtn1},
// 			xtelebot.InlineRow{button.InlineBtn2, button.InlineBtn3},
// 		),
// 	}
func InlineKeyboard(rows ...InlineRow) [][]telebot.InlineButton {
	out := make([][]telebot.InlineButton, 0, len(rows))
	for _, row := range rows {
		columns := make([]telebot.InlineButton, 0, len(row))
		for _, btn := range row {
			columns = append(columns, *btn)
		}
		out = append(out, columns)
	}
	return out
}

// ReplyKeyboard creates a telebot.ReplyButton keyboard with given ReplyRow-s.
//
// Example:
// 	markup := &telebot.ReplyMarkup{
// 		ResizeReplyKeyboard: true,
// 		ReplyKeyboard: xtelebot.ReplyKeyboard(
// 			xtelebot.ReplyRow{button.ReplyBtn1, button.ReplyBtn2},
// 			xtelebot.ReplyRow{button.ReplyBtn3},
// 		),
// 	}
func ReplyKeyboard(rows ...ReplyRow) [][]telebot.ReplyButton {
	out := make([][]telebot.ReplyButton, 0, len(rows))
	for _, row := range rows {
		columns := make([]telebot.ReplyButton, 0, len(row))
		for _, btn := range row {
			columns = append(columns, *btn)
		}
		out = append(out, columns)
	}
	return out
}

// RemoveInlineKeyboard creates a telebot.ReplyMarkup for removing telebot.InlineButton keyboard.
func RemoveInlineKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{InlineKeyboard: nil /* dummy */}
}

// RemoveReplyKeyboard creates a telebot.ReplyMarkup for removing telebot.ReplyButton keyboard.
func RemoveReplyKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{ReplyKeyboardRemove: true}
}

// CallbackShowAlert creates a telebot.CallbackResponse for showing alert in telebot.Callback.
func CallbackShowAlert(text string, showAlert bool) *telebot.CallbackResponse {
	return &telebot.CallbackResponse{Text: text, ShowAlert: showAlert}
}

// ==================
// chat state related
// ==================

// ChatState is a type of chat state, is used in StateHandlerSet and BotData, also can be used in custom finite state machine.
type ChatState uint64

// StateHandlerSet is a container for ChatState and MessageHandler pairs, can be used to get MessageHandler for specific ChatState,
// with using finite state machine, in telebot.OnText handler.
type StateHandlerSet struct {
	handlers map[ChatState]MessageHandler
	muH      sync.RWMutex
}

// NewStateHandlerSet create an empty StateHandlerSet.
func NewStateHandlerSet() *StateHandlerSet {
	return &StateHandlerSet{handlers: make(map[ChatState]MessageHandler)}
}

// IsRegistered checks whether given ChatState has registered a MessageHandler.
func (s *StateHandlerSet) IsRegistered(state ChatState) bool {
	s.muH.RLock()
	_, ok := s.handlers[state]
	s.muH.RUnlock()
	return ok
}

// GetHandler returns the registered MessageHandler for given ChatState, returns nil if handler for given ChatState has not registered yet.
func (s *StateHandlerSet) GetHandler(state ChatState) MessageHandler {
	s.muH.RLock()
	handler, ok := s.handlers[state]
	s.muH.RUnlock()
	if !ok {
		return nil
	}
	return handler
}

const (
	panicNilHandler = "xtelebot: nil handler"
)

// Register registers a MessageHandler for given ChatState, panics when using nil handler.
func (s *StateHandlerSet) Register(state ChatState, handler MessageHandler) {
	if handler == nil {
		panic(panicNilHandler)
	}
	s.muH.Lock()
	s.handlers[state] = handler
	s.muH.Unlock()
}

// Unregister unregisters the MessageHandler for given ChatState.
func (s *StateHandlerSet) Unregister(state ChatState) {
	s.muH.Lock()
	delete(s.handlers, state)
	s.muH.Unlock()
}

// =============
// bot data type
// =============

// BotData represents a set of chats data in a telegram bot, including states and caches.
type BotData struct {
	states map[int64]ChatState
	muS    sync.RWMutex
	caches map[int64]map[string]interface{}
	muC    sync.RWMutex

	initialState ChatState
}

// NewBotData creates a default BotData, with zero initial ChatState.
func NewBotData() *BotData {
	return &BotData{
		states: make(map[int64]ChatState),
		caches: make(map[int64]map[string]interface{}),

		initialState: 0,
	}
}

// ==============
// bot data state
// ==============

// GetStateChats returns all ids from chats which has been set state, the returned slice has no order.
func (b *BotData) GetStateChats() []int64 {
	b.muS.RLock()
	ids := make([]int64, 0, len(b.states))
	for key := range b.states {
		ids = append(ids, key)
	}
	b.muS.RUnlock()
	return ids
}

// GetState returns a chat's state, returns false if no state is set.
func (b *BotData) GetState(chatID int64) (ChatState, bool) {
	b.muS.RLock()
	s, ok := b.states[chatID]
	b.muS.RUnlock()
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

// GetStateOrInit returns a chat's state, sets to the initial state and returns it if no state is set.
func (b *BotData) GetStateOrInit(chatID int64) ChatState {
	b.muS.Lock()
	s, ok := b.states[chatID]
	if !ok {
		s = b.initialState
		b.states[chatID] = s
	}
	b.muS.Unlock()
	return s
}

// InitialState returns the initial ChatState from BotData/
func (b *BotData) InitialState() ChatState {
	b.muS.RLock()
	s := b.initialState
	b.muS.RUnlock()
	return s
}

// SetInitialState sets initial ChatState to BotData.
func (b *BotData) SetInitialState(s ChatState) {
	b.muS.Lock()
	b.initialState = s
	b.muS.Unlock()
}

// SetState sets a chat's state.
func (b *BotData) SetState(chatID int64, state ChatState) {
	b.muS.Lock()
	b.states[chatID] = state
	b.muS.Unlock()
}

// ResetState resets a chat's state to the initial state.
func (b *BotData) ResetState(chatID int64) {
	b.muS.Lock()
	b.states[chatID] = b.initialState
	b.muS.Unlock()
}

// DeleteState deletes a chat's state.
func (b *BotData) DeleteState(chatID int64) {
	b.muS.Lock()
	delete(b.states, chatID)
	b.muS.Unlock()
}

// ==============
// bot data cache
// ==============

// GetCacheChats returns all ids from chats which has been set cache, the returned slice has no order.
func (b *BotData) GetCacheChats() []int64 {
	b.muC.RLock()
	ids := make([]int64, 0, len(b.caches))
	for key := range b.caches {
		ids = append(ids, key)
	}
	b.muC.RUnlock()
	return ids
}

// GetCache returns a chat's cache data, returns false if no cache is set or the key is not found.
func (b *BotData) GetCache(chatID int64, key string) (interface{}, bool) {
	b.muC.RLock()
	if m, ok := b.caches[chatID]; ok {
		if value, ok := m[key]; ok {
			b.muC.RUnlock()
			return value, true
		}
	}
	b.muC.RUnlock()
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

// GetChatCaches returns a chat's all caches data (without any copy), returns false if no cache is set.
func (b *BotData) GetChatCaches(chatID int64) (map[string]interface{}, bool) {
	b.muC.RLock()
	m, ok := b.caches[chatID]
	if !ok {
		b.muC.RUnlock()
		return nil, false
	}
	b.muC.RUnlock()
	return m, true
}

// SetCache sets a chat's cache data using the given key and value.
func (b *BotData) SetCache(chatID int64, key string, value interface{}) {
	b.muC.Lock()
	m, ok := b.caches[chatID]
	if !ok {
		m = make(map[string]interface{})
		b.caches[chatID] = m
	}
	m[key] = value
	b.muC.Unlock()
}

// RemoveCache removes a key from chat's cache.
func (b *BotData) RemoveCache(chatID int64, key string) {
	b.muC.Lock()
	if m, ok := b.caches[chatID]; ok {
		delete(m, key)
	}
	b.muC.Unlock()
}

// ClearCaches clears a chat's all caches.
func (b *BotData) ClearCaches(chatID int64) {
	b.muC.Lock()
	delete(b.caches, chatID)
	b.muC.Unlock()
}
