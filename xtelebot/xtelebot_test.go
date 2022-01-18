package xtelebot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

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
		bd.SetInitialState(InitState)

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

func mockTelebotApi(t *testing.T) (bot *telebot.Bot, shutdown func()) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(func(c *gin.Context) {
		// log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
	})

	app.POST("/botxxx:yyy/getMe", func(c *gin.Context) {
		fakeUser := &telebot.User{ID: 1, FirstName: "FIRSTNAME", LastName: "LASTNAME", Username: "USERNAME", IsBot: true, SupportsInline: true}
		c.JSON(200, gin.H{"ok": true, "Result": fakeUser})
	})
	app.POST("/botxxx:yyy/sendMessage", func(c *gin.Context) {
		bs, _ := ioutil.ReadAll(c.Request.Body)
		data := make(map[string]interface{})
		_ = json.Unmarshal(bs, &data)
		chatId, _ := strconv.ParseInt(data["chat_id"].(string), 10, 64)
		fakeMessage := &telebot.Message{ID: 111, Text: data["text"].(string), Chat: &telebot.Chat{ID: chatId, Username: "?"}, Unixtime: time.Now().Unix() + 2}
		c.JSON(200, gin.H{"ok": true, "Result": fakeMessage})
	})
	app.POST("/botxxx:yyy/getUpdates", func(c *gin.Context) {
		bs, _ := ioutil.ReadAll(c.Request.Body)
		data := make(map[string]interface{})
		_ = json.Unmarshal(bs, &data)
		if data["offset"] != "1" {
			c.JSON(200, gin.H{"ok": true})
			return
		}
		fakeChat1 := &telebot.Chat{ID: 11111111, Username: "?"}
		fakeChat2 := &telebot.Chat{ID: 22222222, Username: "panic"}
		fakeUpdate1 := &telebot.Update{ID: 1, Message: &telebot.Message{ID: 111, Text: "/panic 1", Unixtime: time.Now().Unix(), Chat: fakeChat1}}
		fakeUpdate2 := &telebot.Update{ID: 2, Message: &telebot.Message{ID: 111, Text: "/panic 2", Unixtime: time.Now().Unix(), Chat: fakeChat1}}
		fakeUpdate3 := &telebot.Update{ID: 3, Message: &telebot.Message{ID: 111, Text: "/command something", Unixtime: time.Now().Unix(), Chat: fakeChat1}}
		fakeUpdate4 := &telebot.Update{ID: 4, Message: &telebot.Message{ID: 111, Text: "reply", Unixtime: time.Now().Unix(), Chat: fakeChat1}}
		fakeUpdate5 := &telebot.Update{ID: 4, Message: &telebot.Message{ID: 111, Text: "reply", Unixtime: time.Now().Unix(), Chat: fakeChat2}}
		fakeUpdate6 := &telebot.Update{ID: 4, Callback: &telebot.Callback{Data: "\finline", Message: &telebot.Message{ID: 111, Text: "inline", Unixtime: time.Now().Unix(), Chat: fakeChat1}}}
		fakeUpdate7 := &telebot.Update{ID: 4, Callback: &telebot.Callback{Data: "\finline", Message: &telebot.Message{ID: 111, Text: "inline", Unixtime: time.Now().Unix(), Chat: fakeChat2}}}
		c.JSON(200, gin.H{"ok": true, "Result": []*telebot.Update{fakeUpdate1, fakeUpdate2, fakeUpdate3, fakeUpdate4, fakeUpdate5, fakeUpdate6, fakeUpdate7}})
	})

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()

	mockBot, err := telebot.NewBot(telebot.Settings{
		URL:     "http://localhost:12345",
		Token:   "xxx:yyy",
		Verbose: false,
		Poller:  &telebot.LongPoller{Timeout: time.Millisecond * 200},
	})
	xtesting.Nil(t, err)
	xtesting.Equal(t, mockBot.Me.Username, "USERNAME")
	return mockBot, func() {
		server.Shutdown(context.Background())
	}
}

func TestBotWrapper(t *testing.T) {
	xtesting.Panic(t, func() { NewBotWrapper(nil) })

	mockBot, shutdown := mockTelebotApi(t)
	defer shutdown()
	br := NewBotWrapper(mockBot)
	xtesting.Equal(t, br.Bot(), mockBot)
	xtesting.Equal(t, br.Data().initialState, ChatState(0))

	t.Run("HandleCommand", func(t *testing.T) {
		var defaul = func(w *BotWrapper, m *telebot.Message) {}
		xtesting.Panic(t, func() { br.HandleCommand("", defaul) })
		xtesting.Panic(t, func() { br.HandleCommand("a", defaul) })
		xtesting.Panic(t, func() { br.HandleCommand("aa", defaul) })
		xtesting.Panic(t, func() { br.HandleCommand("/a", nil) })
		xtesting.NotPanic(t, func() { br.HandleCommand("\aa", defaul) })
		xtesting.NotPanic(t, func() { br.HandleCommand("/a", defaul) })
	})

	t.Run("HandleInlineButton", func(t *testing.T) {
		var btn1 = &telebot.InlineButton{Unique: ""}
		var btn2 = &telebot.InlineButton{Unique: "x"}
		var defaul = func(*BotWrapper, *telebot.Callback) {}
		xtesting.Panic(t, func() { br.HandleInlineButton(nil, defaul) })
		xtesting.Panic(t, func() { br.HandleInlineButton(btn1, defaul) })
		xtesting.Panic(t, func() { br.HandleInlineButton(btn2, nil) })
		xtesting.NotPanic(t, func() { br.HandleInlineButton(btn2, defaul) })
	})

	t.Run("HandleInlineButton", func(t *testing.T) {
		var btn1 = &telebot.ReplyButton{Text: ""}
		var btn2 = &telebot.ReplyButton{Text: "x"}
		var defaul = func(*BotWrapper, *telebot.Message) {}
		xtesting.Panic(t, func() { br.HandleReplyButton(nil, defaul) })
		xtesting.Panic(t, func() { br.HandleReplyButton(btn1, defaul) })
		xtesting.Panic(t, func() { br.HandleReplyButton(btn2, nil) })
		xtesting.NotPanic(t, func() { br.HandleReplyButton(btn2, defaul) })
	})

	t.Run("ReplyTo", func(t *testing.T) {
		defaultMsg := &telebot.Message{Chat: &telebot.Chat{ID: 11111111, Username: "Aoi-hosizora"}}
		_, err := br.ReplyTo(nil, "abc")
		xtesting.NotNil(t, err)
		_, err = br.ReplyTo(defaultMsg, nil)
		xtesting.NotNil(t, err)
		_, err = br.ReplyTo(defaultMsg, 0)
		xtesting.NotNil(t, err)
		replied, err := br.ReplyTo(defaultMsg, "abc")
		xtesting.Nil(t, err)
		xtesting.Equal(t, replied.Text, "abc")
		xtesting.Equal(t, replied.Chat.ID, int64(11111111))
	})

	t.Run("SendTo", func(t *testing.T) {
		defaultChat := &telebot.Chat{ID: 11111111, Username: "Aoi-hosizora"}
		_, err := br.SendTo(nil, "abc")
		xtesting.NotNil(t, err)
		_, err = br.SendTo(defaultChat, nil)
		xtesting.NotNil(t, err)
		_, err = br.SendTo(defaultChat, 0)
		xtesting.NotNil(t, err)
		replied, err := br.SendTo(defaultChat, "abc")
		xtesting.Nil(t, err)
		xtesting.Equal(t, replied.Text, "abc")
		xtesting.Equal(t, replied.Chat.ID, int64(11111111))
	})
}

func TestBotWrapperWithPoll(t *testing.T) {
	mockBot, shutdown := mockTelebotApi(t)
	defer shutdown()
	br := NewBotWrapper(mockBot)
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})

	br.SetEndpointHandledCallback(func(endpoint string, handlerName string) {
		l.Debugf("[Telebot] %-12s | %s\n", endpoint, handlerName)
	})
	br.SetReceivedCallback(func(endpoint interface{}, received *telebot.Message) {
		LogReceiveToLogrus(l, endpoint, received)
	})
	br.SetAfterRepliedCallback(func(received *telebot.Message, replied *telebot.Message, err error) {
		LogReplyToLogrus(l, received, replied, err)
	})
	br.SetAfterSentCallback(func(chat *telebot.Chat, sent *telebot.Message, err error) {
		LogSendToLogrus(l, chat, sent, err)
	})

	t.Run("Handle", func(t *testing.T) {
		chs := [7]chan bool{make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool)}
		count := int32(0)
		br.HandleCommand("/panic", func(w *BotWrapper, m *telebot.Message) {
			if m.Text == "/panic 1" {
				defer close(chs[0])
				atomic.AddInt32(&count, 1)
				panic("test panic")
			} else {
				<-chs[0]
				defer close(chs[1])
				atomic.AddInt32(&count, 1)
				origin := br.panicHandler
				defer br.SetPanicHandler(origin)
				br.SetPanicHandler(func(v interface{}) { l.Errorf("Panic with `%v`", v) })
				panic("test panic 2")
			}
		})
		br.HandleCommand("/command", func(w *BotWrapper, m *telebot.Message) {
			<-chs[1]
			defer close(chs[2])
			atomic.AddInt32(&count, 1)
			xtesting.Equal(t, m.Text, "/command something")

			received, err := br.ReplyTo(m, "abc")
			xtesting.Nil(t, err)
			xtesting.Equal(t, received.Text, "abc")
		})
		br.HandleReplyButton(&telebot.ReplyButton{Text: "reply"}, func(w *BotWrapper, m *telebot.Message) {
			if m.Chat.Username != "panic" {
				<-chs[2]
				defer close(chs[3])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, m.Text, "reply")

				received, err := br.ReplyTo(m, "def")
				xtesting.Nil(t, err)
				xtesting.Equal(t, received.Text, "def")
			} else {
				<-chs[3]
				defer close(chs[4])
				atomic.AddInt32(&count, 1)
				panic("test panic reply")
			}
		})
		br.HandleInlineButton(&telebot.InlineButton{Unique: "inline"}, func(w *BotWrapper, c *telebot.Callback) {
			if c.Message.Chat.Username != "panic" {
				<-chs[4]
				defer close(chs[5])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, c.Message.Text, "inline")

				sent, err := br.SendTo(c.Message.Chat, "abc")
				xtesting.Nil(t, err)
				xtesting.Equal(t, sent.Text, "abc")
			} else {
				<-chs[5]
				defer close(chs[6])
				atomic.AddInt32(&count, 1)
				panic("test panic inline")
			}
		})

		terminated := make(chan struct{})
		go func() {
			br.bot.Start()
			close(terminated)
		}()
		<-chs[6]
		br.bot.Stop()
		<-terminated
		xtesting.Equal(t, int(atomic.LoadInt32(&count)), 7)
	})

	t.Run("endpointHandledCallback", func(t *testing.T) {
		// hack
		handledEndpointCallback(nil, "/aaa", func() {})
		handledEndpointCallback(func(s string, s2 string) {}, "", func() {})
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
