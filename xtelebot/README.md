# xtelebot

## Dependencies

+ gopkg.in/tucnak/telebot.v2
+ github.com/sirupsen/logrus

## Documents

### Types

+ `type ChatState uint64`
+ `type BotDataOption func`
+ `type BotData struct`
+ `type LoggerOption func`
+ `type ReceiveLoggerParam struct`
+ `type ReplyLoggerParam struct`
+ `type SendLoggerParam struct`

### Variables

+ `var FormatReceiveFunc func(p *ReceiveLoggerParam) string`
+ `var FieldifyReceiveFunc func(p *ReceiveLoggerParam) logrus.Fields`
+ `var FormatReplyFunc func(p *ReplyLoggerParam) string`
+ `var FieldifyReplyFunc func(p *ReplyLoggerParam) logrus.Fields`
+ `var FormatSendFunc func(p *SendLoggerParam) string`
+ `var FieldifySendFunc func(p *SendLoggerParam) logrus.Fields`

### Constants

+ None

### Functions

+ `func WithInitialState(initialState ChatState) BotDataOption`
+ `func NewBotData(options ...BotDataOption) *BotData`
+ `func WithExtraText(text string) LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) LoggerOption`
+ `func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, received *telebot.Message, options ...LoggerOption)`
+ `func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, received *telebot.Message, options ...LoggerOption)`
+ `func LogReplyToLogrus(logger *logrus.Logger, received, replied *telebot.Message, err error, options ...LoggerOption)`
+ `func LogReplyToLogger(logger logrus.StdLogger, received, replied *telebot.Message, err error, options ...LoggerOption)`
+ `func LogSendToLogrus(logger *logrus.Logger, chat *telebot.Chat, sent *telebot.Message, err error, options ...LoggerOption)`
+ `func LogSendToLogger(logger logrus.StdLogger, chat *telebot.Chat, sent *telebot.Message, err error, options ...LoggerOption)`

### Methods

+ `func (b *BotData) GetStateChats() []int64`
+ `func (b *BotData) GetState(chatID int64) (ChatState, bool)`
+ `func (b *BotData) GetStateOr(chatID int64, fallbackState ChatState) ChatState`
+ `func (b *BotData) GetStateOrInit(chatID int64) ChatState`
+ `func (b *BotData) SetState(chatID int64, state ChatState)`
+ `func (b *BotData) ResetState(chatID int64)`
+ `func (b *BotData) DeleteState(chatID int64)`
+ `func (b *BotData) GetCacheChats() []int64`
+ `func (b *BotData) GetCache(chatID int64, key string) (interface{}, bool)`
+ `func (b *BotData) GetCacheOr(chatID int64, key string, fallbackValue interface{}) interface{}`
+ `func (b *BotData) GetChatCaches(chatID int64) (map[string]interface{}, bool)`
+ `func (b *BotData) SetCache(chatID int64, key string, value interface{})`
+ `func (b *BotData) RemoveCache(chatID int64, key string)`
+ `func (b *BotData) ClearCaches(chatID int64)`
