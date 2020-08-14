package xtelebot

import (
	"sync"
)

type UserStatus uint64

type UsersData struct {
	defState UserStatus
	mus      sync.Mutex
	status   map[int64]UserStatus
	muc      sync.Mutex
	cache    map[int64]map[string]interface{}
}

func NewUsersData(defState UserStatus) *UsersData {
	return &UsersData{
		defState: defState,
		status:   make(map[int64]UserStatus),
		cache:    map[int64]map[string]interface{}{},
	}
}

func (u *UsersData) SetStatus(chatID int64, status UserStatus) {
	u.mus.Lock()
	u.status[chatID] = status
	u.mus.Unlock()
}

func (u *UsersData) GetStatus(chatID int64) UserStatus {
	u.mus.Lock()
	s, ok := u.status[chatID]
	if !ok {
		s = u.defState
		u.status[chatID] = s
	}
	u.mus.Unlock()
	return s
}

func (u *UsersData) ResetStatus(chatID int64) {
	u.mus.Lock()
	u.status[chatID] = u.defState
	u.mus.Unlock()
}

func (u *UsersData) SetCache(chatID int64, key string, value interface{}) {
	u.muc.Lock()
	_, ok := u.cache[chatID]
	if !ok {
		u.cache[chatID] = map[string]interface{}{}
	}

	u.cache[chatID][key] = value
	u.muc.Unlock()
}

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

func (u *UsersData) DeleteCache(chatID int64, key string) {
	u.muc.Lock()
	_, ok := u.cache[chatID]
	if !ok {
		u.cache[chatID] = map[string]interface{}{}
	}

	delete(u.cache[chatID], key)
	u.muc.Unlock()
}
