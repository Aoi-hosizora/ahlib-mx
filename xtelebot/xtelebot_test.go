package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"sync"
	"testing"
)

const (
	None ChatStatus = iota
	Status1
	Status2
)

func TestUsersData(t *testing.T) {
	bd := NewBotData(WithInitialChatStatus(None))

	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := int64(0); i < 20; i++ {
		go func(bd *BotData, i int64) {
			xtesting.Equal(t, bd.GetStatus(i), None)
			bd.SetStatus(i, Status1)
			xtesting.Equal(t, bd.GetStatus(i), Status1)
			bd.SetStatus(i, Status2)
			xtesting.Equal(t, bd.GetStatus(i), Status2)
			bd.ResetStatus(i)
			xtesting.Equal(t, bd.GetStatus(i), None)

			xtesting.Equal(t, bd.GetCacheOr(i, "", nil), nil)
			bd.SetCache(i, "", 0)
			xtesting.Equal(t, bd.GetCacheOr(i, "", nil), 0)
			bd.SetCache(i, "", "")
			xtesting.Equal(t, bd.GetCacheOr(i, "", nil), "")
			bd.RemoveCache(i, "")
			xtesting.Equal(t, bd.GetCacheOr(i, "", nil), nil)

			wg.Done()
		}(bd, i)
	}

	wg.Wait()
}
