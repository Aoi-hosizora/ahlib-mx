package xtelebot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
)

func mockTelebotApi() (bot *telebot.Bot, shutdown func()) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(func(c *gin.Context) {
		// log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
	})

	app.POST("/botxxx:yyy/getMe", func(c *gin.Context) {
		// <<< Me
		fakeUser := &telebot.User{ID: 1, FirstName: "FIRSTNAME", LastName: "LASTNAME", Username: "USERNAME", IsBot: true, SupportsInline: true}
		c.JSON(200, gin.H{"ok": true, "Result": fakeUser})
	})
	app.POST("/botxxx:yyy/sendMessage", func(c *gin.Context) {
		// <<< Send & Reply
		bs, _ := ioutil.ReadAll(c.Request.Body)
		data := make(map[string]interface{})
		_ = json.Unmarshal(bs, &data)
		chatId, _ := strconv.ParseInt(data["chat_id"].(string), 10, 64)
		fakeMessage := &telebot.Message{ID: 111, Text: data["text"].(string), Chat: &telebot.Chat{ID: chatId, Username: "?"}, Unixtime: time.Now().Unix() + 2}
		c.JSON(200, gin.H{"ok": true, "Result": fakeMessage})
	})
	app.POST("/botxxx:yyy/getUpdates", func(c *gin.Context) {
		// <<< Poll
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

	mockedBot, err := telebot.NewBot(telebot.Settings{
		URL:     "http://localhost:12345",
		Token:   "xxx:yyy",
		Verbose: false,
		Poller:  &telebot.LongPoller{Timeout: time.Millisecond * 200},
	})
	if err != nil {
		panic(err)
	}
	return mockedBot, func() {
		server.Shutdown(context.Background())
	}
}

var (
	timestamp = time.Now().Unix()
	chat      = &telebot.Chat{ID: 12345678, Username: "Aoi-hosizora"}
	text      = &telebot.Message{ID: 3344, Chat: chat, Text: "text", Unixtime: timestamp - 2}
	text2     = &telebot.Message{ID: 3345, Chat: chat, Text: "text", Unixtime: timestamp}
)

func TestBotWrapperBasic(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()

	xtesting.Panic(t, func() { NewBotWrapper(nil) })
	br := NewBotWrapper(mockedBot)

	// bot
	xtesting.Equal(t, br.Bot(), mockedBot)
	xtesting.Equal(t, br.Bot().Me.Username, "USERNAME")
	xtesting.Equal(t, strings.Contains(br.Bot().URL, "localhost"), true)
	xtesting.Equal(t, br.Bot().Token, "xxx:yyy")

	// data
	xtesting.Equal(t, br.Data().InitialState(), ChatState(0))
	xtesting.Equal(t, len(br.Data().GetStateChats()), 0)
	xtesting.Equal(t, len(br.Data().GetCacheChats()), 0)
	br.Data().SetInitialState(1)
	br.Data().ResetState(1)
	br.Data().SetCache(1, "key", "value")
	xtesting.Equal(t, br.Data().InitialState(), ChatState(1))
	xtesting.Equal(t, br.Data().GetStateOr(1, 0), ChatState(1))
	xtesting.Equal(t, br.Data().GetCacheOr(1, "key", ""), "value")

	// shs
	xtesting.Equal(t, len(br.Shs().handlers), 0)
	xtesting.Equal(t, br.Shs().IsRegistered(0), false)
	xtesting.Nil(t, br.Shs().GetHandler(0))
	br.Shs().Register(0, func(bw *BotWrapper, m *telebot.Message) {})
	xtesting.Equal(t, br.Shs().IsRegistered(0), true)
	xtesting.NotNil(t, br.Shs().GetHandler(0))

	// other
	xtesting.NotPanic(t, func() {
		br.SetHandledCallback(nil)
		br.SetReceivedCallback(nil)
		br.SetRespondedCallback(nil)
		br.SetPanicHandler(nil)
	})
}

func TestBotWrapperHandle(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()
	br := NewBotWrapper(mockedBot)

	t.Run("HandleCommand", func(t *testing.T) {
		var defaul = func(w *BotWrapper, m *telebot.Message) {}
		xtesting.Panic(t, func() { br.IsHandled(0) })
		xtesting.Panic(t, func() { br.RemoveHandler(0) })
		xtesting.NotPanic(t, func() { br.RemoveHandler("/not_existed") })
		xtesting.Panic(t, func() { br.HandleCommand("", defaul) })
		xtesting.Panic(t, func() { br.HandleCommand("a", defaul) })
		xtesting.Panic(t, func() { br.HandleCommand("aa", defaul) })
		xtesting.Panic(t, func() { br.HandleCommand("\aa", nil) })
		xtesting.Panic(t, func() { br.HandleCommand("/a", nil) })

		xtesting.Equal(t, br.IsHandled("\aa"), false)
		br.HandleCommand("\aa", defaul)
		xtesting.Equal(t, br.IsHandled("\aa"), true)
		xtesting.Equal(t, br.IsHandled("/a"), false)
		br.HandleCommand("/a", defaul)
		xtesting.Equal(t, br.IsHandled("/a"), true)
		br.RemoveHandler("\aa")
		xtesting.Equal(t, br.IsHandled("\aa"), false)
		br.RemoveHandler("/a")
		xtesting.Equal(t, br.IsHandled("/a"), false)
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

	t.Run("HandleReplyButton", func(t *testing.T) {
		var btn1 = &telebot.ReplyButton{Text: ""}
		var btn2 = &telebot.ReplyButton{Text: "x"}
		var defaul = func(*BotWrapper, *telebot.Message) {}
		xtesting.Panic(t, func() { br.HandleReplyButton(nil, defaul) })
		xtesting.Panic(t, func() { br.HandleReplyButton(btn1, defaul) })
		xtesting.Panic(t, func() { br.HandleReplyButton(btn2, nil) })
		xtesting.NotPanic(t, func() { br.HandleReplyButton(btn2, defaul) })
	})
}

func TestBotWrapperRespond(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()
	br := NewBotWrapper(mockedBot)

	t.Run("RespondReply", func(t *testing.T) {
		defaultMsg := &telebot.Message{Chat: &telebot.Chat{ID: 11111111, Username: "Aoi-hosizora"}}
		_, err := br.RespondReply(nil, false, "abc")
		xtesting.NotNil(t, err)
		_, err = br.RespondReply(defaultMsg, false, nil)
		xtesting.NotNil(t, err)
		_, err = br.RespondReply(defaultMsg, false, 0)
		xtesting.NotNil(t, err)
		replied, err := br.RespondReply(defaultMsg, false, "abc")
		xtesting.Nil(t, err)
		xtesting.Equal(t, replied.Text, "abc")
		xtesting.Equal(t, replied.Chat.ID, int64(11111111))
	})

	t.Run("RespondSend", func(t *testing.T) {
		defaultChat := &telebot.Chat{ID: 11111111, Username: "Aoi-hosizora"}
		_, err := br.RespondSend(nil, "abc")
		xtesting.NotNil(t, err)
		_, err = br.RespondSend(defaultChat, nil)
		xtesting.NotNil(t, err)
		_, err = br.RespondSend(defaultChat, 0)
		xtesting.NotNil(t, err)
		replied, err := br.RespondSend(defaultChat, "abc")
		xtesting.Nil(t, err)
		xtesting.Equal(t, replied.Text, "abc")
		xtesting.Equal(t, replied.Chat.ID, int64(11111111))
	})
}

func TestBotWrapperWithPoll(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()
	br := NewBotWrapper(mockedBot)
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})

	br.SetHandledCallback(func(endpoint interface{}, formattedEndpoint string, handlerName string) {
		l.Debugf("[Telebot] %-12s | %s\n", formattedEndpoint, handlerName)
	})
	br.SetReceivedCallback(func(endpoint interface{}, received *telebot.Message) {
		LogReceiveToLogrus(l, endpoint, received)
	})
	// br.SetRepliedCallback(func(received *telebot.Message, replied *telebot.Message, err error) {
	// 	LogReplyToLogrus(l, received, replied, err)
	// })
	// br.SetSentCallback(func(chat *telebot.Chat, sent *telebot.Message, err error) {
	// 	LogRespondToLogrus(l, chat, sent, err)
	// })

	t.Run("Handle", func(t *testing.T) {
		xtesting.Panic(t, func() { xtesting.True(t, br.IsHandled(0)) })
		xtesting.Panic(t, func() { xtesting.True(t, br.IsHandled(struct{}{})) })
		chs := [7]chan bool{make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool)}
		count := int32(0)
		xtesting.False(t, br.IsHandled("/panic"))
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
				br.SetPanicHandler(func(endpoint, _, v interface{}) {
					r, _ := formatEndpoint(endpoint)
					l.Errorf(">> Panic with `%v` | %s", v, r)
					br.SetPanicHandler(origin)
				})
				panic("test panic 2")
			}
		})
		xtesting.True(t, br.IsHandled("/panic"))
		xtesting.False(t, br.IsHandled("/command"))
		br.HandleCommand("/command", func(w *BotWrapper, m *telebot.Message) {
			<-chs[1]
			defer close(chs[2])
			atomic.AddInt32(&count, 1)
			xtesting.Equal(t, m.Text, "/command something")

			received, err := br.RespondReply(m, false, "abc")
			xtesting.Nil(t, err)
			xtesting.Equal(t, received.Text, "abc")
		})
		xtesting.True(t, br.IsHandled("/command"))
		replyBtn := &telebot.ReplyButton{Text: "reply"}
		xtesting.False(t, br.IsHandled(replyBtn))
		br.HandleReplyButton(replyBtn, func(w *BotWrapper, m *telebot.Message) {
			if m.Chat.Username != "panic" {
				<-chs[2]
				defer close(chs[3])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, m.Text, "reply")

				received, err := br.RespondReply(m, false, "def")
				xtesting.Nil(t, err)
				xtesting.Equal(t, received.Text, "def")
			} else {
				<-chs[3]
				defer close(chs[4])
				atomic.AddInt32(&count, 1)
				panic("test panic reply")
			}
		})
		xtesting.True(t, br.IsHandled(replyBtn))
		inlineBtn := &telebot.InlineButton{Unique: "inline"}
		xtesting.False(t, br.IsHandled(inlineBtn))
		br.HandleInlineButton(inlineBtn, func(w *BotWrapper, c *telebot.Callback) {
			if c.Message.Chat.Username != "panic" {
				<-chs[4]
				defer close(chs[5])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, c.Message.Text, "inline")

				sent, err := br.RespondSend(c.Message.Chat, "abc")
				xtesting.Nil(t, err)
				xtesting.Equal(t, sent.Text, "abc")
			} else {
				<-chs[5]
				defer close(chs[6])
				atomic.AddInt32(&count, 1)
				panic("test panic inline")
			}
		})
		xtesting.True(t, br.IsHandled(inlineBtn))

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

	// t.Run("handledCallback", func(t *testing.T) {
	// 	// hack
	// 	processHandledCallback("/aaa", func() {}, nil)
	// 	processHandledCallback("", func() {}, func(e interface{}, s string, s2 string) {})
	// })
}
