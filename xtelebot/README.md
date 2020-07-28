# xtelebot

### Functions

+ `type TelebotLogrus struct {}`
+ `NewTelebotLogrus(logger *logrus.Logger, logMode bool) *TelebotLogrus`
+ `(t *TelebotLogrus) Receive(endpoint interface{}, handle interface{})`
+ `(t *TelebotLogrus) Reply(m *telebot.Message, to *telebot.Message, err error)`
+ `(t *TelebotLogrus) Send(c *telebot.Chat, to *telebot.Message, err error)`
