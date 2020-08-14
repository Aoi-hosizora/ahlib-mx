package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
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
			xtesting.Equal(t, ud.GetStatus(i), None)
			ud.SetStatus(i, Status1)
			xtesting.Equal(t, ud.GetStatus(i), Status1)
			ud.SetStatus(i, Status2)
			xtesting.Equal(t, ud.GetStatus(i), Status2)
			ud.ResetStatus(i)
			xtesting.Equal(t, ud.GetStatus(i), None)

			xtesting.Equal(t, ud.GetCache(i, ""), nil)
			ud.SetCache(i, "", 0)
			xtesting.Equal(t, ud.GetCache(i, ""), 0)
			ud.SetCache(i, "", "")
			xtesting.Equal(t, ud.GetCache(i, ""), "")
			ud.DeleteCache(i, "")
			xtesting.Equal(t, ud.GetCache(i, ""), nil)

			wg.Done()
		}(ud, i)
	}

	wg.Wait()
}
