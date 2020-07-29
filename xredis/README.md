# xredis

### Functions

+ `type RedisLogrus struct {}`
+ `NewRedisLogrus(conn redis.Conn, logger *logrus.Logger) *RedisLogrus`
+ `type RedisLogger struct {}`
+ `NewRedisLogger(conn redis.Conn, logger *log.Logger) *RedisLogger`

### Helper functions

+ `type Helper struct {}`
+ ` WithConn(conn redis.Conn) *Helper`
+ `(h *Helper) DeleteAll(pattern string) (total int, del int, err error)`
