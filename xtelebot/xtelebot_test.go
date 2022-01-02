package xtelebot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewBotData(t *testing.T) {
	for _, tc := range []struct {
		giveOptions []BotDataOption
		wantInitSts ChatState
	}{
		{[]BotDataOption{}, 0},
		{[]BotDataOption{nil}, 0},
		{[]BotDataOption{nil, nil}, 0},
		{[]BotDataOption{WithInitialState(0)}, 0},
		{[]BotDataOption{nil, WithInitialState(0)}, 0},
		{[]BotDataOption{WithInitialState(0), nil}, 0},
		{[]BotDataOption{WithInitialState(1)}, 1},
		{[]BotDataOption{WithInitialState(1), WithInitialState(2)}, 2},
	} {
		xtesting.Equal(t, NewBotData(tc.giveOptions...).option.initialState, tc.wantInitSts)
	}
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
		bd := NewBotData(WithInitialState(InitState))

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

	t.Run("Mutex", func(t *testing.T) {
		bd := NewBotData()
		wg := sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				xtesting.NotPanic(t, func() {
					bd.GetStateChats()
					bd.GetState(0)
					bd.GetStateOr(0, 0)
					bd.GetStateOrInit(0)
					bd.SetState(0, 0)
					bd.ResetState(0)
					bd.DeleteState(0)
				})
			}()
		}
		wg.Wait()
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

	t.Run("Mutex", func(t *testing.T) {
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
	})
}

var (
	timestamp = time.Now().Unix()
	chat      = &telebot.Chat{ID: 12345678, Username: "Aoi-hosizora"}
	text      = &telebot.Message{ID: 3344, Chat: chat, Text: "text", Unixtime: timestamp - 2}
	text2     = &telebot.Message{ID: 3345, Chat: chat, Text: "text", Unixtime: timestamp}
	photo     = &telebot.Message{ID: 3345, Chat: chat, Photo: &telebot.Photo{}, Unixtime: timestamp}
	sticker   = &telebot.Message{ID: 3345, Chat: chat, Sticker: &telebot.Sticker{}, Unixtime: timestamp}
	video     = &telebot.Message{ID: 3345, Chat: chat, Video: &telebot.Video{}, Unixtime: timestamp}
	audio     = &telebot.Message{ID: 3345, Chat: chat, Audio: &telebot.Audio{}, Unixtime: timestamp}
	voice     = &telebot.Message{ID: 3345, Chat: chat, Voice: &telebot.Voice{}, Unixtime: timestamp}
	loc       = &telebot.Message{ID: 3345, Chat: chat, Location: &telebot.Location{}, Unixtime: timestamp}
	animation = &telebot.Message{ID: 3345, Chat: chat, Animation: &telebot.Animation{}, Unixtime: timestamp}
	dice      = &telebot.Message{ID: 3345, Chat: chat, Dice: &telebot.Dice{}, Unixtime: timestamp}
	document  = &telebot.Message{ID: 3345, Chat: chat, Document: &telebot.Document{}, Unixtime: timestamp}
	invoice   = &telebot.Message{ID: 3345, Chat: chat, Invoice: &telebot.Invoice{}, Unixtime: timestamp}
	poll      = &telebot.Message{ID: 3345, Chat: chat, Poll: &telebot.Poll{}, Unixtime: timestamp}
	venue     = &telebot.Message{ID: 3345, Chat: chat, Venue: &telebot.Venue{}, Unixtime: timestamp}
	videoNote = &telebot.Message{ID: 3345, Chat: chat, VideoNote: &telebot.VideoNote{}, Unixtime: timestamp}
)

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
						return fmt.Sprintf("[Telebot] %4d - %30s - %d %s", p.ReceivedID, p.Endpoint, p.ChatID, p.ChatName)
					}
					FieldifyReceiveFunc = func(p *ReceiveLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": p.Action}
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
					FormatSendFunc = func(p *SendLoggerParam) string {
						if p.ErrorMsg != "" {
							return fmt.Sprintf("[Telebot] err: %s", p.ErrorMsg)
						}
						return fmt.Sprintf("[Telebot] %4d - %12s - %8s - %4s - %d %s", p.SentID, "x", p.SentType, "x", p.ChatID, p.ChatName)
					}
					FieldifySendFunc = func(p *SendLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": p.Action}
					}
				}
				if !std {
					LogSendToLogrus(l1, tc.giveChat, tc.giveSent, tc.giveError, tc.giveOptions...)
				} else {
					LogSendToLogger(l2, tc.giveChat, tc.giveSent, tc.giveError, tc.giveOptions...)
				}
				if custom {
					FormatSendFunc = nil
					FieldifySendFunc = nil
				}
			}
		}
	}

	xtesting.NotPanic(t, func() {
		LogSendToLogrus(l1, chat, nil, nil)
		LogSendToLogger(l2, chat, nil, nil)
		LogSendToLogrus(l1, nil, text, nil)
		LogSendToLogger(l2, nil, text, nil)
		LogSendToLogrus(l1, nil, nil, nil)
		LogSendToLogger(l2, nil, nil, nil)
		LogSendToLogrus(nil, nil, nil, nil)
		LogSendToLogger(nil, nil, nil, nil)
	})
}
