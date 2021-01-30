# xserverchan

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/Aoi-hosizora/go-serverchan
+ github.com/sirupsen/logrus

## Documents

### Types

+ `type SendLoggerParam struct`

### Variables

+ `var FormatSendLoggerFunc func(param *SendLoggerParam) string`

### Constants

+ None

### Functions

+ `func WithExtraText(text string) logop.LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) logop.LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption`
+ `func LogSendToLogrus(logger *logrus.Logger, sckey, title string, err error, options ...logop.LoggerOption)`
+ `func LogSendToLogger(logger logrus.StdLogger, sckey, title string, err error, options ...logop.LoggerOption)`

### Methods

+ None
