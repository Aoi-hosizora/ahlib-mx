# xtelebot

### Functions

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
+ `type UserStatus uint64`
+ `type UsersData struct {}`
+ `NewUsersData(noneStatus UserStatus) *UsersData`
+ `(u *UsersData) SetStatus(chatID int64, status UserStatus)`
+ `(u *UsersData) GetStatus(chatID int64) UserStatus`
+ `(u *UsersData) ResetStatus(chatID int64)`
+ `(u *UsersData) SetCache(chatID int64, key string, value interface{})`
+ `(u *UsersData) GetCache(chatID int64, key string) interface{}`
+ `(u *UsersData) DeleteCache(chatID int64, key string)`
