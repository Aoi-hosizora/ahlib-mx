# xredis

### Functions

+ `type RedisLogrus struct {}`
+ `NewRedisLogrus(conn redis.Conn, logger *logrus.Logger, logMode bool) *RedisLogrus`
+ `type RedisLogger struct {}`
+ `NewRedisLogger(conn redis.Conn, logger *log.Logger, logMode bool) *RedisLogger`

### Helper functions

+ `type Helper struct {}`
+ ` WithConn(conn redis.Conn) *Helper`
+ `(h *Helper) DeleteAll(pattern string) (total int, del int, err error)`
+ `(h *Helper) SetAll(keys []string, values []string) (total int, add int, err error)`
+ `(h *Helper) SetExAll(keys []string, values []string, exs []int64) (total int, add int, err error)`
