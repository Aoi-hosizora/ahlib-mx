package xtelebot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
)

type TelebotLogrus struct {
	logger  *logrus.Logger
	logMode bool
}

func NewTelebotLogrus(logger *logrus.Logger, logMode bool) *TelebotLogrus {
	return &TelebotLogrus{logger: logger, logMode: logMode}
}

// Support endpoint type: `string` & `InlineButton` & `ReplyButton`.
// Support handler type: `Message` & `Callback`.
func (t *TelebotLogrus) Receive(endpoint interface{}, handle interface{}) {
	if !t.logMode {
		return
	}

	ep := ""
	if s, ok := endpoint.(string); ok {
		ep = s
	} else if c, ok := endpoint.(telebot.CallbackEndpoint); ok {
		if b, ok := c.(*telebot.InlineButton); ok {
			ep = fmt.Sprintf("$inl:%s", b.Unique)
		} else if b, ok := c.(*telebot.ReplyButton); ok {
			ep = fmt.Sprintf("$rep:%s", b.Text)
		} else {
			ep = fmt.Sprintf("$cb:%T_%v", c, c)
		}
	} else {
		return // unsupported endpoint
	}

	if len(ep) == 0 || ep == "\a" {
		return // empty endpoint
	} else if len(ep) >= 2 && ep[0] == '\a' {
		ep = "$on_" + ep[1:]
	}

	var msgID int
	var chatID int64
	var chatName string

	if msg, ok := handle.(*telebot.Message); ok {
		msgID = msg.ID
		chatID = msg.Chat.ID
		chatName = msg.Chat.Username
	} else if cb, ok := handle.(*telebot.Callback); ok {
		msgID = cb.Message.ID
		chatID = cb.Message.Chat.ID
		chatName = cb.Message.Chat.Username
	} else {
		return // unsupported handle
	}

	t.logger.WithFields(map[string]interface{}{
		"module":    "telebot",
		"messageID": msgID,
		"endpoint":  ep,
		"chatID":    chatID,
		"chatName":  chatName,
	}).Infof("[Telebot] -> %3d | %17v | (%d %s)", msgID, ep, chatID, chatName)
}

func (t *TelebotLogrus) Reply(m *telebot.Message, to *telebot.Message, err error) {
	if !t.logMode || m == nil {
		return
	}

	if err != nil {
		t.logger.Errorf("[Telebot] failed to reply message to %d %s: %v", m.Chat.ID, m.Chat.Username, err)
	} else if to != nil {
		du := to.Time().Sub(m.Time()).String()
		t.logger.WithFields(map[string]interface{}{
			"module":        "telebot",
			"fromMessageId": m.ID,
			"toMessageId":   to.ID,
			"duration":      du,
			"chatID":        to.Chat.ID,
			"chatName":      to.Chat.Username,
		}).Infof("[Telebot] <- %3d | %10s | %4d | (%d %s)", to.ID, du, m.ID, to.Chat.ID, to.Chat.Username)
	}
}

func (t *TelebotLogrus) Send(c *telebot.Chat, to *telebot.Message, err error) {
	if !t.logMode || c == nil {
		return
	}

	if err != nil {
		t.logger.Errorf("[Telebot] failed to send message to %d %s: %v", c.ID, c.Username, err)
	} else if to != nil {
		t.logger.WithFields(map[string]interface{}{
			"module":      "telebot",
			"toMessageId": to.ID,
			"chatID":      to.Chat.ID,
			"chatName":    to.Chat.Username,
		}).Infof("[Telebot] <- %3d | %10s | %4d | (%d %s)", to.ID, "-1", -1, to.Chat.ID, to.Chat.Username)
	}
}
