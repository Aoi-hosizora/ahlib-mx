package xtelebot

import (
	"github.com/go-playground/assert/v2"
	"sync"
	"testing"
)

const (
	None UserStatus = iota
	Status1
	Status2
)

func TestUsersData(t *testing.T) {
	ud := NewUsersData(None)

	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := int64(0); i < 20; i++ {
		go func(ud *UsersData, i int64) {
			assert.Equal(t, ud.GetStatus(i), None)
			ud.SetStatus(i, Status1)
			assert.Equal(t, ud.GetStatus(i), Status1)
			ud.SetStatus(i, Status2)
			assert.Equal(t, ud.GetStatus(i), Status2)
			ud.ResetStatus(i)
			assert.Equal(t, ud.GetStatus(i), None)

			assert.Equal(t, ud.GetCache(i, ""), nil)
			ud.SetCache(i, "", 0)
			assert.Equal(t, ud.GetCache(i, ""), 0)
			ud.SetCache(i, "", "")
			assert.Equal(t, ud.GetCache(i, ""), "")
			ud.DeleteCache(i, "")
			assert.Equal(t, ud.GetCache(i, ""), nil)

			wg.Done()
		}(ud, i)
	}

	wg.Wait()
}
