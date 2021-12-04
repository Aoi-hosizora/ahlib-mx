# xtelebot

## Dependencies

+ gopkg.in/tucnak/telebot.v2
+ github.com/sirupsen/logrus

## Documents

### Types

+ `type ChatStatus uint64`
+ `type BotDataOption func`
+ `type BotData struct`
+ `type LoggerOption func`

### Variables

+ None

### Constants

+ None

### Functions

+ `func WithInitialStatus(initialStatus ChatStatus) BotDataOption`
+ `func NewBotData(options ...BotDataOption) *BotData`
+ `func WithExtraText(text string) logop.LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) logop.LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption`
+ `func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, message *telebot.Message, options ...logop.LoggerOption)`
+ `func LogReplyToLogrus(logger *logrus.Logger, received, replied *telebot.Message, err error, options ...logop.LoggerOption)`
+ `func LogSendToLogrus(logger *logrus.Logger, chat *telebot.Chat, sent *telebot.Message, err error, options ...logop.LoggerOption)`
+ `func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, message *telebot.Message, options ...logop.LoggerOption)`
+ `func LogReplyToLogger(logger logrus.StdLogger, received, replied *telebot.Message, err error, options ...logop.LoggerOption)`
+ `func LogSendToLogger(logger logrus.StdLogger, chat *telebot.Chat, sent *telebot.Message, err error, options ...logop.LoggerOption)`

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
+ `func (b *BotData) ClearCaches(chatID int64)`
