package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestBotData(t *testing.T) {
	for _, tc := range []struct {
		giveOptions []BotDataOption
		wantInitSts ChatStatus
	}{
		{[]BotDataOption{}, 0},
		{[]BotDataOption{nil}, 0},
		{[]BotDataOption{nil, nil}, 0},
		{[]BotDataOption{WithInitialStatus(0)}, 0},
		{[]BotDataOption{nil, WithInitialStatus(0)}, 0},
		{[]BotDataOption{WithInitialStatus(0), nil}, 0},
		{[]BotDataOption{WithInitialStatus(1)}, 1},
		{[]BotDataOption{WithInitialStatus(1), WithInitialStatus(2)}, 2},
	} {
		xtesting.Equal(t, NewBotData(tc.giveOptions...).config.initialStatus, tc.wantInitSts)
	}
}

func TestBotDataStatus(t *testing.T) {
	const (
		None ChatStatus = iota
		InitStatus
		Status
	)
	getOk := func(s ChatStatus, ok bool) bool {
		return ok
	}

	t.Run("Methods", func(t *testing.T) {
		bd := NewBotData(WithInitialStatus(InitStatus))

		// initial
		xtesting.Equal(t, bd.GetStatusChats(), []int64{})    // GetStatusChats
		xtesting.False(t, getOk(bd.GetStatus(0)))            // GetStatus
		xtesting.Equal(t, bd.GetStatusOr(0, None), None)     // GetStatusOr
		xtesting.Equal(t, bd.GetStatusOrInit(0), InitStatus) // GetStatusOrInit

		// empty
		bd.DeleteStatus(0) // DeleteStatus
		xtesting.Equal(t, bd.GetStatusChats(), []int64{})
		xtesting.False(t, getOk(bd.GetStatus(0)))
		xtesting.Equal(t, bd.GetStatusOr(0, None), None)
		xtesting.Equal(t, bd.GetStatusOrInit(0), InitStatus)

		// has status (by init)
		xtesting.Equal(t, bd.GetStatusChats(), []int64{0})
		sts, ok := bd.GetStatus(0)
		xtesting.Equal(t, sts, InitStatus)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetStatusOr(0, None), InitStatus)
		xtesting.Equal(t, bd.GetStatusOrInit(0), InitStatus)

		// has status (by set)
		bd.SetStatus(0, Status)
		sts, ok = bd.GetStatus(0)
		xtesting.Equal(t, sts, Status)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetStatusOr(0, None), Status)
		xtesting.Equal(t, bd.GetStatusOrInit(0), Status)

		// has status (by reset)
		bd.ResetStatus(0) // ResetStatus
		xtesting.Equal(t, bd.GetStatusChats(), []int64{0})
		sts, ok = bd.GetStatus(0)
		xtesting.Equal(t, sts, InitStatus)
		xtesting.True(t, ok)
		xtesting.Equal(t, bd.GetStatusOr(0, None), InitStatus)
		xtesting.Equal(t, bd.GetStatusOrInit(0), InitStatus)
	})

	t.Run("Mutex", func(t *testing.T) {
		bd := NewBotData()
		wg := sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				xtesting.NotPanic(t, func() {
					bd.GetStatusChats()
					bd.GetStatus(0)
					bd.GetStatusOr(0, 0)
					bd.GetStatusOrInit(0)
					bd.SetStatus(0, 0)
					bd.ResetStatus(0)
					bd.DeleteStatus(0)
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
		for i := 0; i < 20; i++ {
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
	unix      = time.Now().Unix()
	chat      = &telebot.Chat{ID: 12345678, Username: "Aoi-hosizora"}
	text      = &telebot.Message{ID: 3344, Chat: chat, Text: "text", Unixtime: unix - 2}
	text2     = &telebot.Message{ID: 3345, Chat: chat, Text: "text", Unixtime: unix}
	photo     = &telebot.Message{ID: 3345, Chat: chat, Photo: &telebot.Photo{}, Unixtime: unix}
	sticker   = &telebot.Message{ID: 3345, Chat: chat, Sticker: &telebot.Sticker{}, Unixtime: unix}
	video     = &telebot.Message{ID: 3345, Chat: chat, Video: &telebot.Video{}, Unixtime: unix}
	audio     = &telebot.Message{ID: 3345, Chat: chat, Audio: &telebot.Audio{}, Unixtime: unix}
	voice     = &telebot.Message{ID: 3345, Chat: chat, Voice: &telebot.Voice{}, Unixtime: unix}
	loc       = &telebot.Message{ID: 3345, Chat: chat, Location: &telebot.Location{}, Unixtime: unix}
	animation = &telebot.Message{ID: 3345, Chat: chat, Animation: &telebot.Animation{}, Unixtime: unix}
	dice      = &telebot.Message{ID: 3345, Chat: chat, Dice: &telebot.Dice{}, Unixtime: unix}
	document  = &telebot.Message{ID: 3345, Chat: chat, Document: &telebot.Document{}, Unixtime: unix}
	invoice   = &telebot.Message{ID: 3345, Chat: chat, Invoice: &telebot.Invoice{}, Unixtime: unix}
	poll      = &telebot.Message{ID: 3345, Chat: chat, Poll: &telebot.Poll{}, Unixtime: unix}
	venue     = &telebot.Message{ID: 3345, Chat: chat, Venue: &telebot.Venue{}, Unixtime: unix}
	videoNote = &telebot.Message{ID: 3345, Chat: chat, VideoNote: &telebot.VideoNote{}, Unixtime: unix}
)

func TestReceiveLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)

	ep := "/test-endpoint"
	on := "\atext"
	rep := &telebot.ReplyButton{Text: "reply"}
	inl := &telebot.InlineButton{Unique: "inline"}

	for _, std := range []bool{false, true} {
		for _, tc := range []struct {
			giveEndpoint interface{}
			giveMessage  *telebot.Message
			giveOptions  []logop.LoggerOption
		}{
			{nil, nil, nil},                      // x
			{"", nil, nil},                       // x
			{"", text, nil},                      // x
			{&telebot.InlineButton{}, text, nil}, // x
			{&telebot.ReplyButton{}, text, nil},  // x

			{ep, text, nil},
			{on, text, nil},
			{rep, text, nil},
			{inl, text, nil},

			{ep, text, []logop.LoggerOption{WithExtraText("extra")}},
			{ep, text, []logop.LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
			{ep, text, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
			{ep, text, []logop.LoggerOption{WithExtraText("extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
			{ep, text, []logop.LoggerOption{WithExtraText("extra"), WithExtraFieldsV("k", "v")}},
		} {
			if !std {
				LogReceiveToLogrus(l1, tc.giveEndpoint, tc.giveMessage, tc.giveOptions...)
			} else {
				LogReceiveToLogger(l2, tc.giveEndpoint, tc.giveMessage, tc.giveOptions...)
			}
		}
	}
}

func TestReplyLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)

	for _, std := range []bool{false, true} {
		for _, tc := range []struct {
			giveReceived *telebot.Message
			giveReplied  *telebot.Message
			giveError    error
			giveOptions  []logop.LoggerOption
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

			{text, text2, nil, []logop.LoggerOption{WithExtraText("extra")}},
			{text, text2, nil, []logop.LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
			{text, text2, nil, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
			{text, text2, nil, []logop.LoggerOption{WithExtraText("extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
			{text, text2, nil, []logop.LoggerOption{WithExtraText("extra"), WithExtraFieldsV("k", "v")}},
			{text, text2, telebot.ErrBlockedByUser, []logop.LoggerOption{WithExtraText("extra")}},
			{text, text2, telebot.ErrBlockedByUser, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
		} {
			if !std {
				LogReplyToLogrus(l1, tc.giveReceived, tc.giveReplied, tc.giveError, tc.giveOptions...)
			} else {
				LogReplyToLogger(l2, tc.giveReceived, tc.giveReplied, tc.giveError, tc.giveOptions...)
			}
		}
	}
}

func TestSendLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)

	for _, std := range []bool{false, true} {
		for _, tc := range []struct {
			giveChat    *telebot.Chat
			giveSent    *telebot.Message
			giveError   error
			giveOptions []logop.LoggerOption
		}{
			{nil, text, nil, nil},
			{chat, nil, nil, nil},
			{chat, nil, telebot.ErrBlockedByUser, nil},

			{chat, text, nil, nil},

			{chat, text, nil, []logop.LoggerOption{WithExtraText("extra")}},
			{chat, text, nil, []logop.LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
			{chat, text, nil, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
			{chat, text, nil, []logop.LoggerOption{WithExtraText("extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
			{chat, text, nil, []logop.LoggerOption{WithExtraText("extra"), WithExtraFieldsV("k", "v")}},
			{chat, text, telebot.ErrBlockedByUser, []logop.LoggerOption{WithExtraText("extra")}},
			{chat, text, telebot.ErrBlockedByUser, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
		} {
			if !std {
				LogSendToLogrus(l1, tc.giveChat, tc.giveSent, tc.giveError, tc.giveOptions...)
			} else {
				LogSendToLogger(l2, tc.giveChat, tc.giveSent, tc.giveError, tc.giveOptions...)
			}
		}
	}
}
