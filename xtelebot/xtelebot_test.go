package xtelebot

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	None ChatStatus = iota
	Status1
	Status2
)

func TestUsersData(t *testing.T) {
	bd := NewBotData(WithInitialStatus(None))

	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := int64(0); i < 20; i++ {
		go func(bd *BotData, i int64) {
			xtesting.Equal(t, bd.GetStatusOrInit(i), None)
			bd.SetStatus(i, Status1)
			xtesting.Equal(t, bd.GetStatusOrInit(i), Status1)
			bd.SetStatus(i, Status2)
			xtesting.Equal(t, bd.GetStatusOrInit(i), Status2)
			bd.ResetStatus(i)
			xtesting.Equal(t, bd.GetStatusOrInit(i), None)

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

func TestXXX(t *testing.T) {
	chat := &telebot.Chat{ID: 12345678, Title: "XXXBot", Username: "Aoi-hosizora"}

	receiveParam := getReceiveLoggerParam("$on_text", &telebot.Message{ID: 3344, Chat: chat})
	log.Println(formatReceiveLogger(receiveParam))
	replyParam := getReplyLoggerParam(&telebot.Message{ID: 3344, Chat: chat, Unixtime: time.Now().Unix() - 2}, &telebot.Message{ID: 3345, Chat: chat, Unixtime: time.Now().Unix(), Text: "hello world"})
	log.Println(formatReplyLogger(replyParam))
	sendParam := getSendLoggerParam(&telebot.Message{ID: 3346, Chat: chat, Unixtime: time.Now().Unix(), Text: "hello world"})
	log.Println(formatSendLogger(sendParam))
}
