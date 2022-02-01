package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"gopkg.in/tucnak/telebot.v2"
	"sync"
	"testing"
	"github.com/sirupsen/logrus"
	"log"
	"fmt"
	"time"
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
		xtesting.Equal(t, shs.IsRegistered(0), false) // IsRegistered
		xtesting.Nil(t, shs.GetHandler(0)) // GetHandler

		// register
		xtesting.Panic(t, func() { shs.Register(0, nil) }) // Register
		tmp := 0
		shs.Register(0, func(bw *BotWrapper, m *telebot.Message) { })
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
		xtesting.NotPanic(t, func() { shs.Unregister(999) }) // Unregister
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
		f := func(bw *BotWrapper, m *telebot.Message) { }
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

func TestReceiveLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.Default()

	ep := "/test-endpoint"
	on := "\atext"
	rep := &telebot.ReplyButton{Text: "test-button1"}
	inl := &telebot.InlineButton{Unique: "test-button2"}

	for _, std := range []bool{false, true} {
		for _, custom := range []bool{false, true} {
			for _, tc := range []struct {
				giveEndpoint interface{}
				giveMessage  *telebot.Message
				giveOptions  []LoggerOption
			}{
				{nil, nil, nil},                      // x
				{"x", nil, nil},                      // x
				{"x", text, nil},                     // x
				{&telebot.InlineButton{}, text, nil}, // x
				{&telebot.ReplyButton{}, text, nil},  // x

				{ep, text, nil},
				{on, text, nil},
				{rep, text, nil},
				{inl, text, nil},

				{ep, text, []LoggerOption{WithExtraText(" | extra")}},
				{ep, text, []LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
				{ep, text, []LoggerOption{WithExtraFieldsV("k", "v")}},
				{ep, text, []LoggerOption{WithExtraText(" | extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
				{ep, text, []LoggerOption{WithExtraText(" | extra"), WithExtraFieldsV("k", "v")}},
			} {
				if custom {
					FormatReceiveFunc = func(p *ReceiveLoggerParam) string {
						return fmt.Sprintf("[Telebot] %4d - %30s - %d %s", p.MessageID, p.FormattedEp, p.ChatID, p.ChatName)
					}
					FieldifyReceiveFunc = func(p *ReceiveLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": "received"}
					}
				}
				if !std {
					LogReceiveToLogrus(l1, tc.giveEndpoint, tc.giveMessage, tc.giveOptions...)
				} else {
					LogReceiveToLogger(l2, tc.giveEndpoint, tc.giveMessage, tc.giveOptions...)
				}
				if custom {
					FormatReceiveFunc = nil
					FieldifyReceiveFunc = nil
				}
			}
		}
	}

	xtesting.NotPanic(t, func() {
		LogReceiveToLogrus(l1, "x", nil)
		LogReceiveToLogger(l2, "x", nil)
		LogReceiveToLogrus(l1, "", text)
		LogReceiveToLogger(l2, "\a", text)
		LogReceiveToLogrus(l1, "\a", text)
		LogReceiveToLogger(l2, 0, text)
		LogReceiveToLogrus(l1, nil, nil)
		LogReceiveToLogger(l2, nil, nil)
		LogReceiveToLogrus(nil, nil, nil)
		LogReceiveToLogger(nil, nil, nil)
	})
}

/*

func TestReplyLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.Default()

	for _, std := range []bool{false, true} {
		for _, custom := range []bool{false, true} {
			for _, tc := range []struct {
				giveReceived *telebot.Message
				giveReplied  *telebot.Message
				giveError    error
				giveOptions  []LoggerOption
			}{
				{text, nil, nil, nil},
				{nil, text, nil, nil},
				{text, nil, telebot.ErrBlockedByUser, nil},

				{text, text2, nil, nil},
				{text, photo, nil, nil},
				{text, sticker, nil, nil},
				{text, video, nil, nil},
				{text, audio, nil, nil},
				{text, voice, nil, nil},
				{text, loc, nil, nil},
				{text, animation, nil, nil},
				{text, dice, nil, nil},
				{text, document, nil, nil},
				{text, invoice, nil, nil},
				{text, poll, nil, nil},
				{text, venue, nil, nil},
				{text, videoNote, nil, nil},

				{text, text2, nil, []LoggerOption{WithExtraText(" | extra")}},
				{text, text2, nil, []LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
				{text, text2, nil, []LoggerOption{WithExtraFieldsV("k", "v")}},
				{text, text2, nil, []LoggerOption{WithExtraText(" | extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
				{text, text2, nil, []LoggerOption{WithExtraText(" | extra"), WithExtraFieldsV("k", "v")}},
				{text, text2, telebot.ErrBlockedByUser, []LoggerOption{WithExtraText(" | extra")}},
				{text, text2, telebot.ErrBlockedByUser, []LoggerOption{WithExtraFieldsV("k", "v")}},
			} {
				if custom {
					FormatReplyFunc = func(p *ReplyLoggerParam) string {
						if p.ErrorMsg != "" {
							return fmt.Sprintf("[Telebot] err: %s", p.ErrorMsg)
						}
						return fmt.Sprintf("[Telebot] %4d - %12s - %8s - %4d - %d %s", p.RepliedID, p.Latency.String(), p.RepliedType, p.ReceivedID, p.ChatID, p.ChatName)

					}
					FieldifyReplyFunc = func(p *ReplyLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": p.Action}
					}
				}
				if !std {
					LogReplyToLogrus(l1, tc.giveReceived, tc.giveReplied, tc.giveError, tc.giveOptions...)
				} else {
					LogReplyToLogger(l2, tc.giveReceived, tc.giveReplied, tc.giveError, tc.giveOptions...)
				}
				if custom {
					FormatReplyFunc = nil
					FieldifyReplyFunc = nil
				}
			}
		}
	}

	xtesting.NotPanic(t, func() {
		LogReplyToLogrus(l1, text, nil, nil)
		LogReplyToLogger(l2, text, nil, nil)
		LogReplyToLogrus(l1, nil, text, nil)
		LogReplyToLogger(l2, nil, text, nil)
		LogReplyToLogrus(l1, nil, nil, nil)
		LogReplyToLogger(l2, nil, nil, nil)
		LogReplyToLogrus(nil, nil, nil, nil)
		LogReplyToLogger(nil, nil, nil, nil)
	})
}

func TestSendLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.Default()

	for _, std := range []bool{false, true} {
		for _, custom := range []bool{false, true} {
			for _, tc := range []struct {
				giveChat    *telebot.Chat
				giveSent    *telebot.Message
				giveError   error
				giveOptions []LoggerOption
			}{
				{nil, text, nil, nil},
				{chat, nil, nil, nil},
				{chat, nil, telebot.ErrBlockedByUser, nil},

				{chat, text, nil, nil},

				{chat, text, nil, []LoggerOption{WithExtraText(" | extra")}},
				{chat, text, nil, []LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
				{chat, text, nil, []LoggerOption{WithExtraFieldsV("k", "v")}},
				{chat, text, nil, []LoggerOption{WithExtraText(" | extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
				{chat, text, nil, []LoggerOption{WithExtraText(" | extra"), WithExtraFieldsV("k", "v")}},
				{chat, text, telebot.ErrBlockedByUser, []LoggerOption{WithExtraText(" | extra")}},
				{chat, text, telebot.ErrBlockedByUser, []LoggerOption{WithExtraFieldsV("k", "v")}},
			} {
				if custom {
					FormatRespondFunc = func(p *RespondLoggerParam) string {
						if p.ErrorMsg != "" {
							return fmt.Sprintf("[Telebot] err: %s", p.ErrorMsg)
						}
						return fmt.Sprintf("[Telebot] %4d - %12s - %8s - %4s - %d %s", p.SentID, "x", p.SentType, "x", p.ChatID, p.ChatName)
					}
					FieldifyRespondFunc = func(p *RespondLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": p.Action}
					}
				}
				if !std {
					LogRespondToLogrus(l1, tc.giveChat, tc.giveSent, tc.giveError, tc.giveOptions...)
				} else {
					LogRespondToLogger(l2, tc.giveChat, tc.giveSent, tc.giveError, tc.giveOptions...)
				}
				if custom {
					FormatRespondFunc = nil
					FieldifyRespondFunc = nil
				}
			}
		}
	}

	xtesting.NotPanic(t, func() {
		LogRespondToLogrus(l1, chat, nil, nil)
		LogRespondToLogger(l2, chat, nil, nil)
		LogRespondToLogrus(l1, nil, text, nil)
		LogRespondToLogger(l2, nil, text, nil)
		LogRespondToLogrus(l1, nil, nil, nil)
		LogRespondToLogger(l2, nil, nil, nil)
		LogRespondToLogrus(nil, nil, nil, nil)
		LogRespondToLogger(nil, nil, nil, nil)
	})
}

*/
