package xtelebot

import (
	"sync"
)

// ChatStatus represents a status of a chat, can be used in fsm.
type ChatStatus uint64

// botDataConfig represents some configs for BotData, set by BotDataOption.
type botDataConfig struct {
	initialStatus ChatStatus // initial status of a chat
}

// BotDataOption represents an option for BotData, created by WithXXX functions.
type BotDataOption func(*botDataConfig)

// WithInitialStatus creates an option for initial chat's status, this defaults to `ChatStatus(0)`.
func WithInitialStatus(initialStatus ChatStatus) BotDataOption {
	return func(b *botDataConfig) {
		b.initialStatus = initialStatus
	}
}

// BotData represents a set of chats data in a bot.
type BotData struct {
	config *botDataConfig // BotData config

	statuses map[int64]ChatStatus             // store statuses of all chats
	mus      sync.RWMutex                     // locks statuses
	caches   map[int64]map[string]interface{} // store caches data of all chats
	muc      sync.RWMutex                     // locks caches
}

// NewBotData creates a new BotData with BotDataOption-s.
func NewBotData(options ...BotDataOption) *BotData {
	config := &botDataConfig{
		initialStatus: ChatStatus(0), // initial status is 0
	}
	for _, option := range options {
		if option != nil {
			option(config)
		}
	}

	return &BotData{
		config:   config,
		statuses: make(map[int64]ChatStatus),
		caches:   make(map[int64]map[string]interface{}),
	}
}

// ======
// status
// ======

// GetStatusChats returns all ids from chats which has been set status, the returned slice has no order.
func (b *BotData) GetStatusChats() []int64 {
	b.mus.RLock()
	ids := make([]int64, 0, len(b.statuses))
	for key := range b.statuses {
		ids = append(ids, key)
	}
	b.mus.RUnlock()
	return ids
}

// GetStatus returns a chat's status, returns false if no status is set.
func (b *BotData) GetStatus(chatID int64) (ChatStatus, bool) {
	b.mus.RLock()
	s, ok := b.statuses[chatID]
	b.mus.RUnlock()
	if !ok {
		return 0, false
	}
	return s, true
}

// GetStatusOr returns a chat's status, returns the fallback status if no status is set.
func (b *BotData) GetStatusOr(chatID int64, fallbackStatus ChatStatus) ChatStatus {
	s, ok := b.GetStatus(chatID)
	if !ok {
		return fallbackStatus
	}
	return s
}

// GetStatusOrInit returns a chat's status, sets to the initial status and returns it if no status is set.
func (b *BotData) GetStatusOrInit(chatID int64) ChatStatus {
	s, ok := b.GetStatus(chatID)
	if !ok {
		s = b.config.initialStatus
		b.mus.Lock()
		b.statuses[chatID] = s
		b.mus.Unlock()
	}
	return s
}

// SetStatus sets a chat's status.
func (b *BotData) SetStatus(chatID int64, status ChatStatus) {
	b.mus.Lock()
	b.statuses[chatID] = status
	b.mus.Unlock()
}

// ResetStatus resets a chat's status to the initial status.
func (b *BotData) ResetStatus(chatID int64) {
	b.mus.Lock()
	b.statuses[chatID] = b.config.initialStatus
	b.mus.Unlock()
}

// DeleteStatus deletes a chat's status.
func (b *BotData) DeleteStatus(chatID int64) {
	b.mus.Lock()
	delete(b.statuses, chatID)
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
	defer b.muc.RUnlock()
	if m, ok := b.caches[chatID]; ok {
		if value, ok := m[key]; ok {
			return value, true
		}
	}
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
	defer b.muc.RUnlock()
	m, ok := b.caches[chatID]
	if !ok {
		return nil, false
	}

	out := make(map[string]interface{}) // copy
	for k, v := range m {
		out[k] = v
	}
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
