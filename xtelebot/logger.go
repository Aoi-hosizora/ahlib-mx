package xtelebot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

type TelebotLogrus struct {
	logger  *logrus.Logger
	LogMode bool
}

// noinspection GoUnusedExportedFunction
func NewTelebotLogrus(logger *logrus.Logger, logMode bool) *TelebotLogrus {
	return &TelebotLogrus{logger: logger, LogMode: logMode}
}

// Support endpoint type: `string` & `InlineButton` & `ReplyButton`.
// Support handler type: `Message` & `Callback`.
func (t *TelebotLogrus) Receive(endpoint interface{}, handle interface{}) {
	if !t.LogMode {
		return
	}

	ep, ok := renderEndpoint(endpoint)
	if !ok {
		return
	}

	var m *telebot.Message
	if msg, ok := handle.(*telebot.Message); ok {
		m = msg
	} else if cb, ok := handle.(*telebot.Callback); ok {
		m = cb.Message
	} else {
		return // unsupported handle
	}
	t.logger.WithFields(map[string]interface{}{
		"module":    "telebot",
		"messageID": m.ID,
		"endpoint":  ep,
		"chatID":    m.Chat.ID,
		"chatName":  m.Chat.Username,
	}).Infof("[Telebot] -> %3d | %17v | (%d %s)", m.ID, ep, m.Chat.ID, m.Chat.Username)
}

func (t *TelebotLogrus) Reply(m *telebot.Message, to *telebot.Message, err error) {
	if !t.LogMode || m == nil {
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
	if !t.LogMode || c == nil {
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

type TelebotLogger struct {
	logger  *log.Logger
	LogMode bool
}

// noinspection GoUnusedExportedFunction
func NewTelebotLogger(logger *log.Logger, logMode bool) *TelebotLogger {
	return &TelebotLogger{logger: logger, LogMode: logMode}
}

func (t *TelebotLogger) Receive(endpoint interface{}, handle interface{}) {
	if !t.LogMode {
		return
	}

	ep, ok := renderEndpoint(endpoint)
	if !ok {
		return
	}

	var m *telebot.Message
	if msg, ok := handle.(*telebot.Message); ok {
		m = msg
	} else if cb, ok := handle.(*telebot.Callback); ok {
		m = cb.Message
	} else {
		return // unsupported handle
	}
	t.logger.Printf("[Telebot] -> %3d | %17v | (%d %s)", m.ID, ep, m.Chat.ID, m.Chat.Username)
}

func (t *TelebotLogger) Reply(m *telebot.Message, to *telebot.Message, err error) {
	if !t.LogMode || m == nil {
		return
	}

	if err != nil {
		t.logger.Printf("[Telebot] failed to reply message to %d %s: %v", m.Chat.ID, m.Chat.Username, err)
	} else if to != nil {
		du := to.Time().Sub(m.Time()).String()
		t.logger.Printf("[Telebot] <- %3d | %10s | %4d | (%d %s)", to.ID, du, m.ID, to.Chat.ID, to.Chat.Username)
	}
}

func (t *TelebotLogger) Send(c *telebot.Chat, to *telebot.Message, err error) {
	if !t.LogMode || c == nil {
		return
	}

	if err != nil {
		t.logger.Printf("[Telebot] failed to send message to %d %s: %v", c.ID, c.Username, err)
	} else if to != nil {
		t.logger.Printf("[Telebot] <- %3d | %10s | %4d | (%d %s)", to.ID, "-1", -1, to.Chat.ID, to.Chat.Username)
	}
}

func renderEndpoint(endpoint interface{}) (string, bool) {
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
		return "", false // unsupported endpoint
	}

	if len(ep) == 0 || ep == "\a" {
		return "", false // empty endpoint
	}
	if len(ep) >= 2 && ep[0] == '\a' {
		ep = "$on_" + ep[1:]
	}

	return ep, true
}
