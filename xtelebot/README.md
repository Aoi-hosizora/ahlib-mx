# xtelebot

## Dependencies

+ github.com/sirupsen/logrus

## Documents

### Types

+ `type ChatStatus uint64`
+ `type BotDataOption func`
+ `type BotData struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func WithInitialChatStatus(initialStatus ChatStatus) BotDataOption`
+ `func NewBotData(options ...BotDataOption) *BotData`

### Methods

+ `func (b *BotData) GetStatusChats() []int64`
+ `func (b *BotData) GetStatus(chatID int64) (ChatStatus, bool)`
+ `func (b *BotData) GetStatusOr(chatID int64, fallbackStatus ChatStatus) ChatStatus`
+ `func (b *BotData) GetStatusOrInit(chatID int64) ChatStatus`
+ `func (b *BotData) SetStatus(chatID int64, status ChatStatus)`
+ `func (b *BotData) ResetStatus(chatID int64)`
+ `func (b *BotData) DeleteStatus(chatID int64)`
+ `func (b *BotData) GetCacheChats() []int64`
+ `func (b *BotData) GetCache(chatID int64, key string) (interface{}, bool)`
+ `func (b *BotData) GetCacheOr(chatID int64, key string, fallbackValue interface{}) interface{}`
+ `func (b *BotData) GetChatCaches(chatID int64) (map[string]interface{}, bool)`
+ `func (b *BotData) SetCache(chatID int64, key string, value interface{})`
+ `func (b *BotData) RemoveCache(chatID int64, key string)`
+ `func (b *BotData) DeleteChatCaches(chatID int64)`

---

### ...

+ `type TelebotLogrus struct {}`
+ `NewTelebotLogrus(logger *logrus.Logger, logMode bool) *TelebotLogrus`
+ `(t *TelebotLogrus) Receive(endpoint interface{}, handle interface{})`
+ `(t *TelebotLogrus) Reply(m *telebot.Message, to *telebot.Message, err error)`
+ `(t *TelebotLogrus) Send(c *telebot.Chat, to *telebot.Message, err error)`
+ `type TelebotLogger struct {}`
+ `NewTelebotLogger(logger *log.Logger, logMode bool) *TelebotLogrus`
+ `(t *TelebotLogger) Receive(endpoint interface{}, handle interface{})`
+ `(t *TelebotLogger) Reply(m *telebot.Message, to *telebot.Message, err error)`
+ `(t *TelebotLogger) Send(c *telebot.Chat, to *telebot.Message, err error)`
