package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"sync/atomic"
	"testing"
	"time"
)

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

	br.SetHandledCallback(func(endpoint interface{}, formattedEndpoint string, handlerName string) {
		l.Debugf("[Telebot] %-12s | %s\n", formattedEndpoint, handlerName)
	})
	br.SetReceivedCallback(func(endpoint interface{}, received *telebot.Message) {
		LogReceiveToLogrus(l, endpoint, received)
	})
	br.SetRepliedCallback(func(received *telebot.Message, replied *telebot.Message, err error) {
		LogReplyToLogrus(l, received, replied, err)
	})
	br.SetSentCallback(func(chat *telebot.Chat, sent *telebot.Message, err error) {
		LogRespondToLogrus(l, chat, sent, err)
	})

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
				br.SetPanicHandler(func(endpoint interface{}, v interface{}) {
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

			received, err := br.ReplyTo(m, "abc")
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
		xtesting.True(t, br.IsHandled(replyBtn))
		inlineBtn := &telebot.InlineButton{Unique: "inline"}
		xtesting.False(t, br.IsHandled(inlineBtn))
		br.HandleInlineButton(inlineBtn, func(w *BotWrapper, c *telebot.Callback) {
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

	t.Run("handledCallback", func(t *testing.T) {
		// hack
		processHandledCallback("/aaa", func() {}, nil)
		processHandledCallback("", func() {}, func(e interface{}, s string, s2 string) {})
	})
}
