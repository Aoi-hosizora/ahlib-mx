package xtelebot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func mockTelebotApi() (bot *telebot.Bot, shutdown func()) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(func(c *gin.Context) {
		// log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
	})

	// mocked apis
	meHandler := func(c *gin.Context) {
		// <<< Me
		fakeUser := &telebot.User{ID: 1, FirstName: "FIRSTNAME", LastName: "LASTNAME", Username: "USERNAME", IsBot: true, SupportsInline: true}
		c.JSON(200, gin.H{"ok": true, "Result": fakeUser})
	}
	messageHandler := func(c *gin.Context) {
		// <<< Send & Reply & Edit
		bs, _ := ioutil.ReadAll(c.Request.Body)
		data := make(map[string]interface{})
		_ = json.Unmarshal(bs, &data)
		chatId, _ := strconv.ParseInt(data["chat_id"].(string), 10, 64)
		fakeMessage := &telebot.Message{ID: 111, Text: data["text"].(string), Chat: &telebot.Chat{ID: chatId, Username: "?"}, Unixtime: time.Now().Unix() + 2}
		c.JSON(200, gin.H{"ok": true, "Result": fakeMessage})
	}
	errorHandler := func(c *gin.Context) {
		// <<< Delete & Respond
		c.JSON(200, gin.H{"ok": true})
	}
	pollHandler := func(c *gin.Context) {
		// <<< Poll
		bs, _ := ioutil.ReadAll(c.Request.Body)
		data := make(map[string]interface{})
		_ = json.Unmarshal(bs, &data)
		if data["offset"] != "1" {
			c.JSON(200, gin.H{"ok": true})
			return
		}
		now := time.Now().Unix()
		fakeChat1 := &telebot.Chat{ID: 11111111, Username: "?"}
		fakeChat2 := &telebot.Chat{ID: 22222222, Username: "panic"}
		fakeUpdate1 := &telebot.Update{ID: 1, Message: &telebot.Message{ID: 111, Text: "/command something", Unixtime: now, Chat: fakeChat1}}
		fakeUpdate2 := &telebot.Update{ID: 2, Message: &telebot.Message{ID: 111, Text: "/command panic1", Unixtime: now, Chat: fakeChat2}}
		fakeUpdate3 := &telebot.Update{ID: 3, Message: &telebot.Message{ID: 111, Text: "/command panic2", Unixtime: now, Chat: fakeChat2}}
		fakeUpdate4 := &telebot.Update{ID: 4, Message: &telebot.Message{ID: 111, Text: "reply", Unixtime: now, Chat: fakeChat1}}
		fakeUpdate5 := &telebot.Update{ID: 5, Message: &telebot.Message{ID: 111, Text: "reply", Unixtime: now, Chat: fakeChat2}}
		fakeUpdate6 := &telebot.Update{ID: 6, Callback: &telebot.Callback{Data: "\finline", Message: &telebot.Message{ID: 111, Text: "inline", Unixtime: now, Chat: fakeChat1}}}
		fakeUpdate7 := &telebot.Update{ID: 7, Callback: &telebot.Callback{Data: "\finline", Message: &telebot.Message{ID: 111, Text: "inline", Unixtime: now, Chat: fakeChat2}}}
		c.JSON(200, gin.H{"ok": true, "Result": []*telebot.Update{fakeUpdate1, fakeUpdate2, fakeUpdate3, fakeUpdate4, fakeUpdate5, fakeUpdate6, fakeUpdate7}})
	}
	app.POST("/botxxx:yyy/getMe", meHandler)
	app.POST("/botxxx:yyy/sendMessage", messageHandler)
	app.POST("/botxxx:yyy/editMessageText", messageHandler)
	app.POST("/botxxx:yyy/deleteMessage", errorHandler)
	app.POST("/botxxx:yyy/answerCallbackQuery", errorHandler)
	app.POST("/botxxx:yyy/getUpdates", pollHandler)

	// serve the mocked api
	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()

	mockedBot, err := telebot.NewBot(telebot.Settings{
		URL:     "http://localhost:12345",
		Token:   "xxx:yyy",
		Verbose: false,
		Poller:  LongPoller(0),
		// Poller:  &telebot.LongPoller{Timeout: time.Millisecond * 200},
	})
	if err != nil {
		panic(err)
	}
	return mockedBot, func() {
		server.Shutdown(context.Background())
	}
}

func TestBotWrapperBasic(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()

	xtesting.Panic(t, func() { NewBotWrapper(nil) })
	bw := NewBotWrapper(mockedBot)

	// bot
	xtesting.Equal(t, bw.Bot(), mockedBot)
	xtesting.Equal(t, bw.Bot().Me.Username, "USERNAME")
	xtesting.Equal(t, strings.Contains(bw.Bot().URL, "localhost"), true)
	xtesting.Equal(t, bw.Bot().Token, "xxx:yyy")

	// data
	xtesting.Equal(t, bw.Data().InitialState(), ChatState(0))
	xtesting.Equal(t, len(bw.Data().GetStateChats()), 0)
	xtesting.Equal(t, len(bw.Data().GetCacheChats()), 0)
	bw.Data().SetInitialState(1)
	bw.Data().ResetState(1)
	bw.Data().SetCache(1, "key", "value")
	xtesting.Equal(t, bw.Data().InitialState(), ChatState(1))
	xtesting.Equal(t, bw.Data().GetStateOr(1, 0), ChatState(1))
	xtesting.Equal(t, bw.Data().GetCacheOr(1, "key", ""), "value")

	// shs
	xtesting.Equal(t, len(bw.Shs().handlers), 0)
	xtesting.Equal(t, bw.Shs().IsRegistered(0), false)
	xtesting.Nil(t, bw.Shs().GetHandler(0))
	bw.Shs().Register(0, func(bw *BotWrapper, m *telebot.Message) {})
	xtesting.Equal(t, bw.Shs().IsRegistered(0), true)
	xtesting.NotNil(t, bw.Shs().GetHandler(0))

	// formatEndpoint
	for _, tc := range []struct {
		give   interface{}
		want   string
		wantOk bool
	}{
		{"", "", false},
		{"a", "", false},
		{"/a", "/a", true},
		{"\aa", "$on_a", true},
		{&telebot.ReplyButton{}, "", false},
		{&telebot.ReplyButton{Text: "text"}, "$rep_btn:text", true},
		{&telebot.InlineButton{}, "", false},
		{&telebot.InlineButton{Unique: "unique"}, "$inl_btn:unique", true},
		{nil, "", false},
		{0, "", false},
	} {
		s, ok := formatEndpoint(tc.give)
		xtesting.Equal(t, s, tc.want)
		xtesting.Equal(t, ok, tc.wantOk)
	}

	// other
	xtesting.NotPanic(t, func() {
		bw.SetHandledCallback(nil)
		bw.SetReceivedCallback(nil)
		bw.SetRespondedCallback(nil)
		bw.SetPanicHandler(nil)
	})
}

func TestBotWrapperHandle(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi() // serving is no need actually
	defer shutdown()
	bw := NewBotWrapper(mockedBot)

	t.Run("HandleCommand", func(t *testing.T) {
		var empty = func(w *BotWrapper, m *telebot.Message) {}
		xtesting.Panic(t, func() { bw.IsHandled(0) })
		xtesting.Panic(t, func() { bw.IsHandled(struct{}{}) })
		xtesting.Panic(t, func() { bw.RemoveHandler(0) })
		xtesting.NotPanic(t, func() { bw.RemoveHandler("/not_existed") })
		xtesting.Panic(t, func() { bw.HandleCommand("", empty) })
		xtesting.Panic(t, func() { bw.HandleCommand("a", empty) })
		xtesting.Panic(t, func() { bw.HandleCommand("aa", empty) })
		xtesting.Panic(t, func() { bw.HandleCommand("\aa", nil) })
		xtesting.Panic(t, func() { bw.HandleCommand("/a", nil) })

		xtesting.Equal(t, bw.IsHandled("\aa"), false)
		bw.HandleCommand("\aa", empty)
		xtesting.Equal(t, bw.IsHandled("\aa"), true)
		xtesting.Equal(t, bw.IsHandled("/a"), false)
		bw.SetHandledCallback(DefaultColorizedHandledCallback)
		bw.HandleCommand("/a", empty)
		bw.SetHandledCallback(DefaultHandledCallback)
		xtesting.Equal(t, bw.IsHandled("/a"), true)
		bw.RemoveHandler("\aa")
		xtesting.Equal(t, bw.IsHandled("\aa"), false)
		bw.RemoveHandler("/a")
		xtesting.Equal(t, bw.IsHandled("/a"), false)
	})

	t.Run("HandleInlineButton", func(t *testing.T) {
		var btn1 = &telebot.InlineButton{Unique: ""}
		var btn2 = &telebot.InlineButton{Unique: "x"}
		var empty = func(*BotWrapper, *telebot.Callback) {}
		xtesting.Panic(t, func() { bw.HandleInlineButton(nil, empty) })
		xtesting.Panic(t, func() { bw.HandleInlineButton(btn1, empty) })
		xtesting.Panic(t, func() { bw.HandleInlineButton(btn2, nil) })

		xtesting.Equal(t, bw.IsHandled(btn2), false)
		bw.HandleInlineButton(btn2, empty)
		xtesting.Equal(t, bw.IsHandled(btn2), true)
		bw.RemoveHandler(btn2)
		xtesting.Equal(t, bw.IsHandled(btn2), false)
	})

	t.Run("HandleReplyButton", func(t *testing.T) {
		var btn1 = &telebot.ReplyButton{Text: ""}
		var btn2 = &telebot.ReplyButton{Text: "x"}
		var empty = func(*BotWrapper, *telebot.Message) {}
		xtesting.Panic(t, func() { bw.HandleReplyButton(nil, empty) })
		xtesting.Panic(t, func() { bw.HandleReplyButton(btn1, empty) })
		xtesting.Panic(t, func() { bw.HandleReplyButton(btn2, nil) })

		xtesting.Equal(t, bw.IsHandled(btn2), false)
		bw.HandleReplyButton(btn2, empty)
		xtesting.Equal(t, bw.IsHandled(btn2), true)
		bw.RemoveHandler(btn2)
		xtesting.Equal(t, bw.IsHandled(btn2), false)
	})
}

var (
	defaultChat     = &telebot.Chat{ID: 11111111, Username: "Aoi-hosizora"}
	defaultMessage  = &telebot.Message{Text: "...", Chat: defaultChat}
	defaultCallback = &telebot.Callback{Message: defaultMessage}
	defaultAnswer   = &telebot.CallbackResponse{Text: "alert", ShowAlert: true}
)

func TestBotWrapperRespond(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()
	bw := NewBotWrapper(mockedBot)
	getErr := func(_ *telebot.Message, err error) error { return err }

	typ := RespondEventType("")
	event := (*RespondEvent)(nil)
	bw.SetRespondedCallback(func(t RespondEventType, ev *RespondEvent) { typ = t; event = ev })

	t.Run("RespondSend", func(t *testing.T) {
		xtesting.NotNil(t, getErr(bw.RespondSend(nil, "text")))
		xtesting.NotNil(t, getErr(bw.RespondSend(defaultChat, nil)))
		xtesting.NotNil(t, getErr(bw.RespondSend(defaultChat, 0))) // ErrUnsupportedWhat

		msg, err := bw.RespondSend(defaultChat, "abc", telebot.NoPreview)
		xtesting.Nil(t, err)
		xtesting.Equal(t, msg.Text, "abc")
		xtesting.Equal(t, msg.Chat.ID, defaultChat.ID)
		xtesting.Equal(t, typ, RespondSendEvent)
		xtesting.Equal(t, event.SendSource, defaultChat)
		xtesting.Equal(t, event.SendWhat, "abc")
		xtesting.Equal(t, event.SendOptions, []interface{}{telebot.NoPreview})
		xtesting.Equal(t, event.SendResult, msg)

		msg, err = bw.RespondSend(defaultChat, "123") // &telebot.Photo{}
		xtesting.Nil(t, err)
		xtesting.Equal(t, msg.Text, "123")
		xtesting.Equal(t, event.SendWhat, "123")
		xtesting.Equal(t, event.SendOptions, []interface{}(nil))
		xtesting.Equal(t, event.SendResult, msg)
	})

	t.Run("RespondReply", func(t *testing.T) {
		xtesting.NotNil(t, getErr(bw.RespondReply(nil, false, "text")))
		xtesting.NotNil(t, getErr(bw.RespondReply(defaultMessage, true, nil)))
		xtesting.NotNil(t, getErr(bw.RespondReply(defaultMessage, false, 0))) // ErrUnsupportedWhat

		msg, err := bw.RespondReply(defaultMessage, false, "abc", telebot.NoPreview)
		xtesting.Nil(t, err)
		xtesting.Equal(t, msg.Text, "abc")
		xtesting.Equal(t, msg.Chat.ID, defaultChat.ID)
		xtesting.Equal(t, typ, RespondReplyEvent)
		xtesting.Equal(t, event.ReplySource, defaultMessage)
		xtesting.Equal(t, event.ReplyExplicit, false)
		xtesting.Equal(t, event.ReplyWhat, "abc")
		xtesting.Equal(t, event.ReplyOptions, []interface{}{telebot.NoPreview})
		xtesting.Equal(t, event.ReplyResult, msg)

		msg, err = bw.RespondReply(defaultMessage, true, "123")
		xtesting.Nil(t, err)
		xtesting.Equal(t, msg.Text, "123")
		xtesting.Equal(t, event.ReplyExplicit, true)
		xtesting.Equal(t, event.ReplyWhat, "123")
		xtesting.Equal(t, event.ReplyOptions, []interface{}(nil))
		xtesting.Equal(t, event.ReplyResult, msg)
	})

	t.Run("RespondEdit", func(t *testing.T) {
		xtesting.NotNil(t, getErr(bw.RespondEdit(nil, "text")))
		xtesting.NotNil(t, getErr(bw.RespondEdit(defaultMessage, nil)))
		xtesting.NotNil(t, getErr(bw.RespondEdit(defaultMessage, 0))) // ErrUnsupportedWhat

		msg, err := bw.RespondEdit(defaultMessage, "abc", telebot.NoPreview)
		xtesting.Nil(t, err)
		xtesting.Equal(t, msg.Text, "abc")
		xtesting.Equal(t, msg.Chat.ID, defaultChat.ID)
		xtesting.Equal(t, typ, RespondEditEvent)
		xtesting.Equal(t, event.EditSource, defaultMessage)
		xtesting.Equal(t, event.EditWhat, "abc")
		xtesting.Equal(t, event.EditOptions, []interface{}{telebot.NoPreview})
		xtesting.Equal(t, event.EditResult, msg)

		msg, err = bw.RespondEdit(defaultMessage, "123")
		xtesting.Nil(t, err)
		xtesting.Equal(t, msg.Text, "123")
		xtesting.Equal(t, event.EditWhat, "123")
		xtesting.Equal(t, event.EditOptions, []interface{}(nil))
		xtesting.Equal(t, event.EditResult, msg)
	})

	t.Run("RespondDelete", func(t *testing.T) {
		xtesting.NotNil(t, bw.RespondDelete(nil))

		err := bw.RespondDelete(defaultMessage)
		xtesting.Nil(t, err)
		xtesting.Equal(t, typ, RespondDeleteEvent)
		xtesting.Equal(t, event.DeleteSource, defaultMessage)
		xtesting.Equal(t, event.DeleteResult, defaultMessage)
	})

	t.Run("RespondCallback", func(t *testing.T) {
		xtesting.NotNil(t, bw.RespondCallback(nil, &telebot.CallbackResponse{}))
		xtesting.Nil(t, bw.RespondCallback(defaultCallback, nil)) // answer is nillable
		xtesting.Nil(t, bw.RespondCallback(defaultCallback, &telebot.CallbackResponse{}))

		err := bw.RespondCallback(defaultCallback, defaultAnswer)
		xtesting.Nil(t, err)
		xtesting.Equal(t, typ, RespondCallbackEvent)
		xtesting.Equal(t, event.CallbackSource, defaultCallback)
		xtesting.Equal(t, event.CallbackAnswer, defaultAnswer)
		xtesting.Equal(t, event.CallbackResult, defaultAnswer)

		err = bw.RespondCallback(defaultCallback, nil)
		xtesting.Nil(t, err)
		xtesting.Equal(t, event.CallbackAnswer, (*telebot.CallbackResponse)(nil))
		xtesting.Equal(t, event.CallbackResult, (*telebot.CallbackResponse)(nil))
	})
}

func TestBotWrapperWithPoll(t *testing.T) {
	mockedBot, shutdown := mockTelebotApi()
	defer shutdown()
	bw := NewBotWrapper(mockedBot)
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})

	bw.SetHandledCallback(func(endpoint interface{}, formattedEndpoint string, handlerName string) {
		l.Debugf("[Telebot] %-12s | %s\n", formattedEndpoint, handlerName)
	})
	bw.SetReceivedCallback(func(endpoint interface{}, received *telebot.Message) {
		LogReceiveToLogrus(l, endpoint, received)
	})
	bw.SetRespondedCallback(func(typ RespondEventType, event *RespondEvent) {
		LogRespondToLogrus(l, typ, event)
	})

	t.Run("Handle", func(t *testing.T) {
		chs := [7]chan bool{make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool), make(chan bool)}
		count := int32(0)

		// command
		xtesting.False(t, bw.IsHandled("/command"))
		bw.HandleCommand("/command", func(w *BotWrapper, m *telebot.Message) {
			if m.Payload == "something" {
				defer close(chs[0])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, m.Text, "/command something")
				bw.RespondSend(m.Chat, m.Text) // <<< test send
			} else if m.Payload == "panic1" {
				<-chs[0]
				defer close(chs[1])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, m.Text, "/command panic1")
				origin := bw.panicHandler
				bw.SetPanicHandler(func(endpoint, _, v interface{}) {
					r, _ := formatEndpoint(endpoint)
					l.Errorf(">> Panic with `%v` | %s", v, r)
					bw.SetPanicHandler(origin)
				})
				panic("test panic 1") // new handler
			} else if m.Payload == "panic2" {
				<-chs[1]
				defer close(chs[2])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, m.Text, "/command panic2")
				bw.RespondReply(m, true, m.Text) // <<< test reply
				panic("test panic 2")            // default handler
			}
		})
		xtesting.True(t, bw.IsHandled("/command"))

		// reply
		replyBtn := &telebot.ReplyButton{Text: "reply"}
		xtesting.False(t, bw.IsHandled(replyBtn))
		bw.HandleReplyButton(replyBtn, func(w *BotWrapper, m *telebot.Message) {
			if m.Chat.Username != "panic" {
				<-chs[2]
				defer close(chs[3])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, m.Text, "reply")
				bw.RespondEdit(m, m.Text) // <<< test edit
			} else {
				<-chs[3]
				defer close(chs[4])
				atomic.AddInt32(&count, 1)
				bw.RespondDelete(m) // <<< test delete
				panic("test panic reply")
			}
		})
		xtesting.True(t, bw.IsHandled(replyBtn))

		// inline
		inlineBtn := &telebot.InlineButton{Unique: "inline"}
		xtesting.False(t, bw.IsHandled(inlineBtn))
		bw.HandleInlineButton(inlineBtn, func(w *BotWrapper, c *telebot.Callback) {
			if c.Message.Chat.Username != "panic" {
				<-chs[4]
				defer close(chs[5])
				atomic.AddInt32(&count, 1)
				xtesting.Equal(t, c.Message.Text, "inline")
				bw.RespondCallback(c, &telebot.CallbackResponse{Text: "alert", ShowAlert: false}) // test respond
			} else {
				<-chs[5]
				defer close(chs[6])
				atomic.AddInt32(&count, 1)
				panic("test panic inline")
			}
		})
		xtesting.True(t, bw.IsHandled(inlineBtn))

		// test count
		terminated := make(chan struct{})
		go func() {
			bw.bot.Start()
			close(terminated)
		}()
		<-chs[6]
		bw.bot.Stop()
		<-terminated
		xtesting.Equal(t, int(atomic.LoadInt32(&count)), 7)
	})
}

var (
	timestamp = time.Now().Unix()
	chat      = &telebot.Chat{ID: 12345678, Username: "Aoi-hosizora"}
	text      = &telebot.Message{ID: 3344, Chat: chat, Text: "text", Unixtime: timestamp - 1}
	text2     = &telebot.Message{ID: 3345, Chat: chat, Text: "text", Unixtime: timestamp + 1}
	callback  = &telebot.Callback{Message: text}
	cbresp    = &telebot.CallbackResponse{Text: "x"}
)

// ATTENTION: loggerOptions related code and unit tests in xgin package and xtelebot package should keep the same as each other.
func TestLoggerOptions(t *testing.T) {
	for _, tc := range []struct {
		give       []LoggerOption
		wantMsg    string
		wantFields logrus.Fields
	}{
		{[]LoggerOption{}, "", logrus.Fields{}},
		{[]LoggerOption{nil}, "", logrus.Fields{}},
		{[]LoggerOption{nil, nil, nil}, "", logrus.Fields{}},

		{[]LoggerOption{WithExtraText("")}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("   ")}, "   ", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("  x x  ")}, "  x x  ", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("test")}, "test", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("test1"), WithExtraText(" | test2")}, " | test2", logrus.Fields{}},

		{[]LoggerOption{WithExtraFields(map[string]interface{}{})}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"true": 2, "3": 4.4})}, "", logrus.Fields{"true": 2, "3": 4.4}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"true": 2, "3": 4.4}), WithExtraFields(map[string]interface{}{"k": "v"})}, "", logrus.Fields{"k": "v"}},

		{[]LoggerOption{WithExtraFieldsV()}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil)}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil, "a", nil)}, "", logrus.Fields{"a": nil}},
		{[]LoggerOption{WithExtraFieldsV(nil, nil, "a")}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil, nil, 1, nil)}, "", logrus.Fields{"1": nil}},
		{[]LoggerOption{WithExtraFieldsV(true, 2, 3.3, 4, 5)}, "", logrus.Fields{"true": 2, "3.3": 4}},
		{[]LoggerOption{WithExtraFieldsV(true, 2, 3.3, 4, 5), WithExtraFieldsV("k", "v")}, "", logrus.Fields{"k": "v"}},

		{[]LoggerOption{WithExtraText("test"), WithExtraFields(map[string]interface{}{"1": 2})}, "test", logrus.Fields{"1": 2}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"1": 2}), WithExtraText(" | test")}, " | test", logrus.Fields{"1": 2}},
		{[]LoggerOption{WithExtraText("test"), WithExtraFieldsV(3, 4)}, "test", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFieldsV(3, 4), WithExtraText(" | test")}, " | test", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"1": 2}), WithExtraFieldsV(3, 4)}, "", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFieldsV(3, 4), WithExtraFields(map[string]interface{}{"1": 2})}, "", logrus.Fields{"1": 2}},
	} {
		ops := buildLoggerOptions(tc.give)
		msg := ""
		fields := logrus.Fields{}
		ops.ApplyToMessage(&msg)
		ops.ApplyToFields(fields)
		xtesting.Equal(t, msg, tc.wantMsg)
		xtesting.Equal(t, fields, tc.wantFields)
	}
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
				{nil, text, nil},               // x
				{"x", nil, nil},                // x
				{"x", &telebot.Message{}, nil}, // x

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
						return fmt.Sprintf("[Telebot] recv - msg# %7s - %28s - %d %s", xnumber.Itoa(p.MessageID), p.FormattedEp, p.ChatID, p.ChatName)
					}
					FieldifyReceiveFunc = func(p *ReceiveLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": "receive"}
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
		LogReceiveToLogger(l2, "", text)
		LogReceiveToLogrus(l1, "\a", text)
		LogReceiveToLogger(l2, "\a", text)
		LogReceiveToLogrus(l1, 0, text)
		LogReceiveToLogger(l2, 0, text)
		LogReceiveToLogrus(l1, "\a", &telebot.Message{})
		LogReceiveToLogger(l2, "\a", &telebot.Message{})
		LogReceiveToLogrus(l1, nil, nil)
		LogReceiveToLogger(l2, nil, nil)
		LogReceiveToLogrus(nil, nil, nil)
		LogReceiveToLogger(nil, nil, nil)
	})

	t.Run("colorizeEventType", func(t *testing.T) {
		xtesting.Equal(t, colorizeEventType(""), xcolor.Blue.Sprintf("recv")) // <<<
		xtesting.Equal(t, colorizeEventType(RespondSendEvent), xcolor.Green.Sprint("send"))
		xtesting.Equal(t, colorizeEventType(RespondReplyEvent), xcolor.Green.Sprint(" rep"))
		xtesting.Equal(t, colorizeEventType(RespondEditEvent), xcolor.Yellow.Sprint("edit"))
		xtesting.Equal(t, colorizeEventType(RespondDeleteEvent), xcolor.Red.Sprint(" del"))
		xtesting.Equal(t, colorizeEventType(RespondCallbackEvent), xcolor.Cyan.Sprint("call"))
		xtesting.Equal(t, colorizeEventType("x"), " ???")
	})
}

func TestRespondLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.Default()

	for _, std := range []bool{false, true} {
		for _, custom := range []bool{false, true} {
			for _, tc := range []struct {
				giveType    RespondEventType
				giveEvent   *RespondEvent
				giveOptions []LoggerOption
			}{
				{"", nil, nil},     // x
				{"send", nil, nil}, // x

				{RespondSendEvent, &RespondEvent{SendSource: chat, SendResult: text}, nil},
				{RespondReplyEvent, &RespondEvent{ReplySource: text, ReplyResult: text2}, nil},
				{RespondEditEvent, &RespondEvent{EditSource: text, EditResult: text2}, nil},
				{RespondDeleteEvent, &RespondEvent{DeleteSource: text, DeleteResult: text2}, nil},
				{RespondCallbackEvent, &RespondEvent{CallbackSource: callback, CallbackAnswer: cbresp, CallbackResult: cbresp}, nil},
				{RespondSendEvent, &RespondEvent{SendSource: chat, ReturnedError: telebot.ErrBlockedByUser}, nil},
				{RespondCallbackEvent, &RespondEvent{CallbackSource: callback, CallbackResult: &telebot.CallbackResponse{Text: "x", ShowAlert: true}, ReturnedError: telebot.ErrBlockedByUser}, nil},

				{RespondSendEvent, &RespondEvent{SendSource: chat, SendResult: text}, []LoggerOption{WithExtraText(" | extra")}},
				{RespondReplyEvent, &RespondEvent{ReplySource: text, ReplyResult: text2}, []LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
				{RespondEditEvent, &RespondEvent{EditSource: text, EditResult: text2}, []LoggerOption{WithExtraFieldsV("k", "v")}},
				{RespondDeleteEvent, &RespondEvent{DeleteSource: text, DeleteResult: text2}, []LoggerOption{WithExtraText(" | extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
				{RespondCallbackEvent, &RespondEvent{CallbackSource: callback, CallbackAnswer: nil, CallbackResult: &telebot.CallbackResponse{}}, []LoggerOption{WithExtraText(" | extra"), WithExtraFieldsV("k", "v")}},
			} {
				if custom {
					FormatRespondFunc = func(p *RespondLoggerParam) string {
						if p.ReturnedErrorMsg != "" {
							return fmt.Sprintf("[Telebot] %s - msg# %7s - %d %s - err: %s", p.EventType, xnumber.Itoa(p.ResultMessageID), p.SourceChatID, p.SourceChatName, p.ReturnedErrorMsg)
						}
						return fmt.Sprintf("[Telebot] %s - msg# %7s - %d %s", p.EventType, xnumber.Itoa(p.ResultMessageID), p.SourceChatID, p.SourceChatName)
					}
					FieldifyRespondFunc = func(p *RespondLoggerParam) logrus.Fields {
						return logrus.Fields{"module": "telebot", "action": string(p.EventType)}
					}
				}
				if !std {
					LogRespondToLogrus(l1, tc.giveType, tc.giveEvent, tc.giveOptions...)
				} else {
					LogRespondToLogger(l2, tc.giveType, tc.giveEvent, tc.giveOptions...)
				}
				if custom {
					FormatRespondFunc = nil
					FieldifyRespondFunc = nil
				}
			}
		}
	}

	xtesting.NotPanic(t, func() {
		xtesting.Equal(t, formatRespondLoggerParam(&RespondLoggerParam{EventType: "x"}), "")
		LogRespondToLogrus(l1, RespondSendEvent, nil)
		LogRespondToLogger(l2, RespondSendEvent, nil)
		LogRespondToLogrus(l1, "", &RespondEvent{})
		LogRespondToLogger(l2, "", &RespondEvent{})
		LogRespondToLogrus(l1, "send", &RespondEvent{})
		LogRespondToLogger(l2, "send", &RespondEvent{})
		LogRespondToLogrus(l1, " rep", &RespondEvent{ReplySource: text})
		LogRespondToLogger(l2, " rep", &RespondEvent{ReplySource: text})
		LogRespondToLogrus(l1, "edit", &RespondEvent{EditSource: text, ReturnedError: telebot.ErrBlockedByUser})
		LogRespondToLogger(l2, "edit", &RespondEvent{EditSource: text, ReturnedError: telebot.ErrBlockedByUser})
	})
}
