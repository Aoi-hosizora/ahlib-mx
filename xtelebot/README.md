# xtelebot

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ gopkg.in/tucnak/telebot.v2
+ github.com/sirupsen/logrus

## Documents

### Types

+ `type InlineRow []*telebot.InlineButton`
+ `type ReplyRow []*telebot.ReplyButton`
+ `type ChatState uint64`
+ `type StateHandlerSet struct`
+ `type BotData struct`
+ `type BotWrapper struct`
+ `type MessageHandler func`
+ `type CallbackHandler func`
+ `type RespondEventType string`
+ `type RespondEvent struct`
+ `type LoggerOption func`
+ `type ReceiveLoggerParam struct`
+ `type RespondLoggerParam struct`

### Variables

+ `var FormatReceiveFunc func(p *ReceiveLoggerParam) string`
+ `var FieldifyReceiveFunc func(p *ReceiveLoggerParam) logrus.Fields`
+ `var FormatRespondFunc func(p *RespondLoggerParam) string`
+ `var FieldifyRespondFunc func(p *RespondLoggerParam) logrus.Fields`

### Constants

+ `const RespondSendEvent RespondEventType`
+ `const RespondReplyEvent RespondEventType`
+ `const RespondEditEvent RespondEventType`
+ `const RespondDeleteEvent RespondEventType`
+ `const RespondCallbackEvent RespondEventType`

### Functions

+ `func TextBtn(text string) *telebot.ReplyButton`
+ `func DataBtn(text, unique string, data ...string) *telebot.InlineButton`
+ `func URLBtn(text, url string) *telebot.InlineButton`
+ `func InlineKeyboard(rows ...InlineRow) [][]telebot.InlineButton`
+ `func ReplyKeyboard(rows ...ReplyRow) [][]telebot.ReplyButton`
+ `func RemoveInlineKeyboard() *telebot.ReplyMarkup`
+ `func RemoveReplyKeyboard() *telebot.ReplyMarkup`
+ `func CallbackShowAlert(text string, showAlert bool) *telebot.CallbackResponse`
+ `func NewStateHandlerSet() *StateHandlerSet`
+ `func NewBotData() *BotData`
+ `func NewBotWrapper(bot *telebot.Bot) *BotWrapper`
+ `func WithExtraText(text string) LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) LoggerOption`
+ `func LogReceiveToLogrus(logger *logrus.Logger, endpoint interface{}, received *telebot.Message, options ...LoggerOption)`
+ `func LogReceiveToLogger(logger logrus.StdLogger, endpoint interface{}, received *telebot.Message, options ...LoggerOption)`
+ `func LogRespondToLogrus(logger *logrus.Logger, typ RespondEventType, ev *RespondEvent, options ...LoggerOption)`
+ `func LogRespondToLogger(logger logrus.StdLogger, typ RespondEventType, ev *RespondEvent, options ...LoggerOption)`

### Methods

+ `func (s *StateHandlerSet) IsRegistered(state ChatState) bool`
+ `func (s *StateHandlerSet) GetHandler(state ChatState) MessageHandler`
+ `func (s *StateHandlerSet) Register(state ChatState, handler MessageHandler)`
+ `func (s *StateHandlerSet) Unregister(state ChatState)`
+ `func (b *BotData) GetStateChats() []int64`
+ `func (b *BotData) GetState(chatID int64) (ChatState, bool)`
+ `func (b *BotData) GetStateOr(chatID int64, fallbackState ChatState) ChatState`
+ `func (b *BotData) GetStateOrInit(chatID int64) ChatState`
+ `func (b *BotData) SetInitialState(s ChatState)`
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
+ `func (b *BotWrapper) Bot() *telebot.Bot`
+ `func (b *BotWrapper) Data() *BotData`
+ `func (b *BotWrapper) Shs() *StateHandlerSet`
+ `func (b *BotWrapper) IsHandled(endpoint interface{}) bool`
+ `func (b *BotWrapper) RemoveHandler(endpoint interface{})`
+ `func (b *BotWrapper) HandleCommand(command string, handler MessageHandler)`
+ `func (b *BotWrapper) HandleReplyButton(button *telebot.ReplyButton, handler MessageHandler)`
+ `func (b *BotWrapper) HandleInlineButton(button *telebot.InlineButton, handler CallbackHandler)`
+ `func (b *BotWrapper) RespondSend(source *telebot.Chat, what interface{}, options ...interface{}) (*telebot.Message, error)`
+ `func (b *BotWrapper) RespondReply(source *telebot.Message, explicit bool, what interface{}, options ...interface{}) (*telebot.Message, error)`
+ `func (b *BotWrapper) RespondEdit(source *telebot.Message, what interface{}, options ...interface{}) (*telebot.Message, error)`
+ `func (b *BotWrapper) RespondDelete(source *telebot.Message) error`
+ `func (b *BotWrapper) RespondCallback(source *telebot.Callback, answer *telebot.CallbackResponse) error`
+ `func (b *BotWrapper) SetHandledCallback(f func(endpoint interface{}, formattedEndpoint string, handlerName string))`
+ `func (b *BotWrapper) SetReceivedCallback(cb func(endpoint interface{}, received *telebot.Message))`
+ `func (b *BotWrapper) SetRespondedCallback(cb func(typ RespondEventType, event *RespondEvent))`
+ `func (b *BotWrapper) SetPanicHandler(handler func(endpoint, source, value interface{}))`
