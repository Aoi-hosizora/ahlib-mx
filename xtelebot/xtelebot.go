package xtelebot

import (
	"sync"
)

// ======================
// UserStatus & UsersData
// ======================

// UserStatus represents a user's status, can be used in fsm.
type UserStatus uint64

// UsersData represents some data belongs to all users.
type UsersData struct {
	defStatus UserStatus                       // user's default status
	mus       sync.Mutex                       // status mutex
	status    map[int64]UserStatus             // store all user's status
	muc       sync.Mutex                       // cache mutex
	cache     map[int64]map[string]interface{} // store all user's cache
}

// UsersData creates a new UsersData with a default UserStatus.
func NewUsersData(defStatus UserStatus) *UsersData {
	return &UsersData{
		defStatus: defStatus,
		status:    make(map[int64]UserStatus),
		cache:     map[int64]map[string]interface{}{},
	}
}

// SetStatus sets a user's status.
func (u *UsersData) SetStatus(chatID int64, status UserStatus) {
	u.mus.Lock()
	u.status[chatID] = status
	u.mus.Unlock()
}

// GetStatus gets and returns a user's status.
func (u *UsersData) GetStatus(chatID int64) UserStatus {
	u.mus.Lock()
	s, ok := u.status[chatID]
	if !ok {
		s = u.defStatus
		u.status[chatID] = s
	}
	u.mus.Unlock()
	return s
}

// ResetStatus resets a user's status to the default status.
func (u *UsersData) ResetStatus(chatID int64) {
	u.mus.Lock()
	u.status[chatID] = u.defStatus
	u.mus.Unlock()
}

// SetStatus sets a user's cache data from the given key.
func (u *UsersData) SetCache(chatID int64, key string, value interface{}) {
	u.muc.Lock()
	_, ok := u.cache[chatID]
	if !ok {
		u.cache[chatID] = map[string]interface{}{}
	}

	u.cache[chatID][key] = value
	u.muc.Unlock()
}

// GetCache gets and returns a user's cache data from the given key.
func (u *UsersData) GetCache(chatID int64, key string) interface{} {
	u.muc.Lock()
	_, ok := u.cache[chatID]
	if !ok {
		u.cache[chatID] = map[string]interface{}{}
	}

	value, ok := u.cache[chatID][key]
	u.muc.Unlock()
	if !ok {
		return nil
	}
	return value
}

// DeleteCache deletes a user's cache data from the given key.
func (u *UsersData) DeleteCache(chatID int64, key string) {
	u.muc.Lock()
	_, ok := u.cache[chatID]
	if !ok {
		u.cache[chatID] = map[string]interface{}{}
	}

	delete(u.cache[chatID], key)
	u.muc.Unlock()
}
