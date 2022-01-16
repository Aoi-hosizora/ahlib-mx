package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"gopkg.in/tucnak/telebot.v2"
	"sync"
	"testing"
)

func TestMarkups(t *testing.T) {
	t.Run("xxx_btn", func(t *testing.T) {
		data := DataBtn("text", "unique", "data1", "data2", "data3")
		xtesting.Equal(t, data.Text, "text")
		xtesting.Equal(t, data.Unique, "unique")
		xtesting.Equal(t, data.Data, "data1|data2|data3")
		text := TextBtn("text")
		xtesting.Equal(t, text.Text, "text")
		url := URLBtn("text", "url")
		xtesting.Equal(t, url.Text, "text")
		xtesting.Equal(t, url.URL, "url")
	})

	t.Run("xxx_keyboard", func(t *testing.T) {
		inlines := InlineKeyboard(
			InlineRow{DataBtn("text1", "unique1"), DataBtn("text2", "unique2")},
			InlineRow{DataBtn("text3", "unique3")},
		)
		replies := ReplyKeyboard(
			ReplyRow{TextBtn("text1"), TextBtn("text2")},
			ReplyRow{TextBtn("text3")},
		)
		xtesting.Equal(t, inlines[0][0], *DataBtn("text1", "unique1"))
		xtesting.Equal(t, inlines[0][1], *DataBtn("text2", "unique2"))
		xtesting.Equal(t, inlines[1][0], *DataBtn("text3", "unique3"))
		xtesting.Equal(t, replies[0][0], *TextBtn("text1"))
		xtesting.Equal(t, replies[0][1], *TextBtn("text2"))
		xtesting.Equal(t, replies[1][0], *TextBtn("text3"))

		rInlines := SetInlineKeyboard(inlines)
		rReplies := SetReplyKeyboard(replies)
		xtesting.Equal(t, rInlines.InlineKeyboard, inlines)
		xtesting.Equal(t, rReplies.ReplyKeyboard, replies)
		xtesting.Equal(t, rReplies.ResizeReplyKeyboard, true)
	})

	t.Run("mass_functions", func(t *testing.T) {
		xtesting.Equal(t, RemoveInlineKeyboard(), &telebot.ReplyMarkup{InlineKeyboard: nil})
		xtesting.Equal(t, RemoveReplyKeyboard(), &telebot.ReplyMarkup{ReplyKeyboardRemove: true})
		xtesting.Equal(t, CallbackShowAlert("text", true), &telebot.CallbackResponse{Text: "text", ShowAlert: true})
	})
}

func TestStateHandlerSet(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		shs := NewStateHandlerSet()

		// initial
		xtesting.Equal(t, shs.IsRegistered(0), false)
		xtesting.Nil(t, shs.GetHandler(0))

		// register
		xtesting.Panic(t, func() { shs.Register(0, nil) })
		tmp := 0
		shs.Register(0, func(bw *BotWrapper, m *telebot.Message) {})
		xtesting.Equal(t, shs.IsRegistered(0), true)
		xtesting.NotNil(t, shs.GetHandler(0))
		shs.Register(0, func(bw *BotWrapper, m *telebot.Message) { tmp++ })
		xtesting.Equal(t, tmp, 0)
		shs.GetHandler(0)(nil, nil)
		xtesting.Equal(t, tmp, 1)
		shs.Register(1, func(bw *BotWrapper, m *telebot.Message) { tmp += 2 })
		shs.GetHandler(1)(nil, nil)
		xtesting.Equal(t, tmp, 3)

		// unregister
		xtesting.NotPanic(t, func() { shs.Unregister(999) })
		xtesting.True(t, shs.IsRegistered(1))
		xtesting.NotNil(t, shs.GetHandler(1))
		shs.Unregister(1)
		xtesting.False(t, shs.IsRegistered(1))
		xtesting.Nil(t, shs.GetHandler(1))
		xtesting.True(t, shs.IsRegistered(0))
		xtesting.NotNil(t, shs.GetHandler(0))
		shs.Unregister(0)
		xtesting.False(t, shs.IsRegistered(0))
		xtesting.Nil(t, shs.GetHandler(0))
		xtesting.Equal(t, len(shs.handlers), 0)
	})

	t.Run("Concurrency", func(t *testing.T) {
		shs := NewStateHandlerSet()
		wg := sync.WaitGroup{}
		f := func(bw *BotWrapper, m *telebot.Message) {}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xtesting.NotPanic(t, func() {
					shs.IsRegistered(0)
					shs.GetHandler(0)
					shs.Register(0, f)
					shs.Unregister(0)
				})
			}()
		}
		wg.Wait()
		xtesting.Equal(t, len(shs.handlers), 0)
	})
}

func TestBotDataState(t *testing.T) {
	const (
		None ChatState = iota
		InitState
		State
	)
	getOk := func(s ChatState, ok bool) bool {
		return ok
	}

	t.Run("Methods", func(t *testing.T) {
		bd := NewBotData()
		xtesting.Equal(t, bd.InitialState(), ChatState(0))
		bd.SetInitialState(InitState)
		xtesting.Equal(t, bd.InitialState(), InitState)

		// initial
		xtesting.Equal(t, bd.GetStateChats(), []int64{})   // GetStateChats
		xtesting.False(t, getOk(bd.GetState(0)))           // GetState
		xtesting.Equal(t, bd.GetStateOr(0, None), None)    // GetStateOr
		xtesting.Equal(t, bd.GetStateOrInit(0), InitState) // GetStateOrInit

		// empty
		bd.DeleteState(0) // DeleteState
		xtesting.Equal(t, bd.GetStateChats(), []int64{})
		xtesting.False(t, getOk(bd.GetState(0)))
		xtesting.Equal(t, bd.GetStateOr(0, None), None)
		xtesting.Equal(t, bd.GetStateOrInit(0), InitState)

		// has state (by init)
		xtesting.Equal(t, bd.GetStateChats(), []int64{0})
		sts, ok := bd.GetState(0)
		xtesting.Equal(t, sts, InitState)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetStateOr(0, None), InitState)
		xtesting.Equal(t, bd.GetStateOrInit(0), InitState)

		// has state (by set)
		bd.SetState(0, State)
		sts, ok = bd.GetState(0)
		xtesting.Equal(t, sts, State)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetStateOr(0, None), State)
		xtesting.Equal(t, bd.GetStateOrInit(0), State)

		// has state (by reset)
		bd.ResetState(0) // ResetState
		xtesting.Equal(t, bd.GetStateChats(), []int64{0})
		sts, ok = bd.GetState(0)
		xtesting.Equal(t, sts, InitState)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetStateOr(0, None), InitState)
		xtesting.Equal(t, bd.GetStateOrInit(0), InitState)
	})

	t.Run("Concurrency", func(t *testing.T) {
		bd := NewBotData()
		wg := sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xtesting.NotPanic(t, func() {
					bd.GetStateChats()
					bd.GetState(0)
					bd.GetStateOr(0, 1)
					bd.InitialState()
					bd.SetInitialState(1)
					bd.GetStateOrInit(0)
					bd.SetState(0, 0)
					bd.ResetState(0)
					bd.DeleteState(0)
				})
			}()
		}
		wg.Wait()
		xtesting.Equal(t, len(bd.states), 0)
	})
}

func TestBotDataCache(t *testing.T) {
	getOk := func(s interface{}, ok bool) bool {
		return ok
	}

	t.Run("Methods", func(t *testing.T) {
		bd := NewBotData()

		// initial
		xtesting.Equal(t, bd.GetCacheChats(), []int64{}) // GetCacheChats
		xtesting.False(t, getOk(bd.GetCache(0, "")))     // GetCache
		xtesting.Equal(t, bd.GetCacheOr(0, "", 0), 0)    // GetCacheOr
		xtesting.False(t, getOk(bd.GetChatCaches(0)))    // GetChatCaches

		// empty
		bd.SetCache(0, "", 0) // SetCache
		bd.ClearCaches(0)     // ClearCaches
		xtesting.Equal(t, bd.GetCacheChats(), []int64{})
		xtesting.False(t, getOk(bd.GetCache(0, "")))
		xtesting.Equal(t, bd.GetCacheOr(0, "", 0), 0)
		xtesting.False(t, getOk(bd.GetChatCaches(0)))

		// has chat, empty cache
		bd.SetCache(0, "", 0)
		bd.RemoveCache(0, "") // RemoveCache
		xtesting.Equal(t, bd.GetCacheChats(), []int64{0})
		xtesting.False(t, getOk(bd.GetCache(0, "")))
		xtesting.Equal(t, bd.GetCacheOr(0, "", 0), 0)
		m, ok := bd.GetChatCaches(0)
		xtesting.Equal(t, m, map[string]interface{}{})
		xtesting.True(t, ok)

		// has chat, has cache
		bd.SetCache(0, "", 0)
		xtesting.Equal(t, bd.GetCacheChats(), []int64{0})
		data, ok := bd.GetCache(0, "")
		xtesting.Equal(t, data, 0)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetCacheOr(0, "", 1), 0)
		m, ok = bd.GetChatCaches(0)
		xtesting.Equal(t, m, map[string]interface{}{"": 0})
		xtesting.True(t, ok)
	})

	t.Run("Concurrency", func(t *testing.T) {
		bd := NewBotData()
		wg := sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xtesting.NotPanic(t, func() {
					bd.GetCacheChats()
					bd.GetCache(0, "")
					bd.GetCacheOr(0, "", 0)
					bd.GetChatCaches(0)
					bd.SetCache(0, "", 0)
					bd.RemoveCache(0, "")
					bd.ClearCaches(0)
				})
			}()
		}
		wg.Wait()
		xtesting.Equal(t, len(bd.caches), 0)
	})
}
