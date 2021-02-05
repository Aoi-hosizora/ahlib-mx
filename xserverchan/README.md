# xserverchan

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/Aoi-hosizora/go-serverchan
+ github.com/sirupsen/logrus

## Documents

### Types

+ None

### Variables

+ None

### Constants

+ None

### Functions

+ `func WithExtraText(text string) logop.LoggerOption`
+ `func WithExtraFields(fields map[string]interface{}) logop.LoggerOption`
+ `func WithExtraFieldsV(fields ...interface{}) logop.LoggerOption`
+ `func LogToLogrus(logger *logrus.Logger, sckey, title, body string, err error, options ...logop.LoggerOption)`
+ `func LogToLogger(logger logrus.StdLogger, sckey, title, body string, err error, options ...logop.LoggerOption)`

### Methods

+ None
