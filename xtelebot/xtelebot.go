package xtelebot

import (
	"sync"
)

// ChatState is a type of chat states, can be used in fsm.
type ChatState uint64

// botDataOptions is a type of BotData's option, each field can be set by BotDataOption function type.
type botDataOptions struct {
	initialState ChatState
}

// BotDataOption represents an option for BotData's options, can be created by WithXXX functions.
type BotDataOption func(*botDataOptions)

// WithInitialState creates an BotDataOption for initial chat's state, defaults to `ChatState(0)`.
func WithInitialState(initialState ChatState) BotDataOption {
	return func(b *botDataOptions) {
		b.initialState = initialState
	}
}

// BotData represents a set of chats data (including states and caches) in a telegram bot.
type BotData struct {
	option *botDataOptions

	states map[int64]ChatState
	mus    sync.RWMutex
	caches map[int64]map[string]interface{}
	muc    sync.RWMutex
}

// NewBotData creates a new BotData with BotDataOption-s.
func NewBotData(options ...BotDataOption) *BotData {
	opt := &botDataOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}

	return &BotData{
		option: opt,
		states: make(map[int64]ChatState),
		caches: make(map[int64]map[string]interface{}),
	}
}

// =====
// state
// =====

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
	if !ok {
		return 0, false
	}
	return s, true
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
	b.mus.Lock()
	s, ok := b.states[chatID]
	if !ok {
		s = b.option.initialState
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
	b.states[chatID] = b.option.initialState
	b.mus.Unlock()
}

// DeleteState deletes a chat's state.
func (b *BotData) DeleteState(chatID int64) {
	b.mus.Lock()
	delete(b.states, chatID)
	b.mus.Unlock()
}

// =====
// cache
// =====

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
	out := make(map[string]interface{}) // copy
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
