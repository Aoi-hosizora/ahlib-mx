# xrecovery

### Functions

+ `WithExtraString(s string) LoggerOption`
+ `WithExtraFields(m map[string]interface{}) LoggerOption`
+ `WithExtraFieldsV(m ...interface{}) LoggerOption`
+ `WithLogrus(logger *logrus.Logger, err interface{}, options ...LoggerOption)`
+ `WithLogger(logger *log.Logger, err interface{}, options ...LoggerOption)`
