package xtelebot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

// LogReceiveToLogrus logs the receive events.
//
// Support endpoint type:
// 	string, InlineButton, ReplyButton
// Support handler type:
// 	Message, Callback
func LogReceiveToLogrus(logger *logrus.Logger, endpoint, handler interface{}) {
	// information
	ep, ok := renderEndpoint(endpoint)
	if !ok {
		return
	}
	var m *telebot.Message
	if msg, ok := handler.(*telebot.Message); ok {
		m = msg
	} else if cb, ok := handler.(*telebot.Callback); ok {
		m = cb.Message
	} else {
		return // unsupported handle
	}

	// fields
	fields := logrus.Fields{
		"module":     "telebot",
		"message_id": m.ID,
		"endpoint":   ep,
		"chat_id":    m.Chat.ID,
		"chat_name":  m.Chat.Username,
	}
	entry := logger.WithFields(fields)

	// logger
	msg := fmt.Sprintf("[Telebot] %4d | -> | %17v | (%d %s)", m.ID, ep, m.Chat.ID, m.Chat.Username)
	entry.Info(msg)
}

// LogReplyToLogrus logs the reply events.
func LogReplyToLogrus(logger *logrus.Logger, m *telebot.Message, to *telebot.Message, err error) {
	if m == nil {
		return
	}

	if err != nil {
		logger.Error(fmt.Sprintf("[Telebot] Reply to %d %s error: %v", m.Chat.ID, m.Chat.Username, err))
	} else if to != nil {
		du := to.Time().Sub(m.Time()).String()
		logger.WithFields(map[string]interface{}{
			"module":        "telebot",
			"fromMessageId": m.ID,
			"toMessageId":   to.ID,
			"duration":      du,
			"chatID":        to.Chat.ID,
			"chatName":      to.Chat.Username,
		}).Info(fmt.Sprintf("[Telebot] %4d | %12s | %4d | (%d %s)", to.ID, du, m.ID, to.Chat.ID, to.Chat.Username))
	}
}

// LogSendToLogrus logs the send events.
func LogSendToLogrus(logger *logrus.Logger, c *telebot.Chat, to *telebot.Message, err error) {
	if c == nil {
		return
	}

	if err != nil {
		logger.Error(fmt.Sprintf("[Telebot] Send to %d %s error: %v", c.ID, c.Username, err))
	} else if to != nil {
		logger.WithFields(map[string]interface{}{
			"module":      "telebot",
			"toMessageId": to.ID,
			"chatID":      to.Chat.ID,
			"chatName":    to.Chat.Username,
		}).Info(fmt.Sprintf("[Telebot] %4d | %12s | %4d | (%d %s)", to.ID, "-1", -1, to.Chat.ID, to.Chat.Username))
	}
}

// TelebotLogger is a standard logger used by telebot.
type TelebotLogger struct {
	logger  *log.Logger
	LogMode bool
}

// NewTelebotLogger creates a TelebotLogger with log.Logger.
func NewTelebotLogger(logger *log.Logger, logMode bool) *TelebotLogger {
	return &TelebotLogger{logger: logger, LogMode: logMode}
}

// Receive logs the receive events.
//
// Support endpoint type:
// 	string, InlineButton, ReplyButton
// Support handler type:
// 	Message, Callback
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

	t.logger.Printf("[Telebot] %4d | -> | %17v | (%d %s)", m.ID, ep, m.Chat.ID, m.Chat.Username)
}

// Reply logs the reply events.
func (t *TelebotLogger) Reply(m *telebot.Message, to *telebot.Message, err error) {
	if !t.LogMode || m == nil {
		return
	}

	if err != nil {
		t.logger.Printf("[Telebot] Reply to %d %s error: %v", m.Chat.ID, m.Chat.Username, err)
	} else if to != nil {
		du := to.Time().Sub(m.Time()).String()
		t.logger.Printf("[Telebot] %4d | %12s | %4d | (%d %s)", to.ID, du, m.ID, to.Chat.ID, to.Chat.Username)
	}
}

// Send logs the send events.
func (t *TelebotLogger) Send(c *telebot.Chat, to *telebot.Message, err error) {
	if !t.LogMode || c == nil {
		return
	}

	if err != nil {
		t.logger.Printf("[Telebot] Send to %d %s error: %v", c.ID, c.Username, err)
	} else if to != nil {
		t.logger.Printf("[Telebot] %4d | %12s | %4d | (%d %s)", to.ID, "-1", -1, to.Chat.ID, to.Chat.Username)
	}
}

// renderEndpoint renders the endpoint value to a string.
//
// Support endpoints:
// 	string, InlineButton, ReplyButton
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
