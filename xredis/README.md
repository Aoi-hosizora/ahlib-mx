# xredis

### Functions

+ `DeleteAll(conn redis.Conn, pattern string) (total int, del int, err error)`

### Logger Functions

+ `type RedisLogrus struct {}`
+ `NewRedisLogrus(conn redis.Conn, logger *logrus.Logger) *RedisLogrus`
+ `type RedisLogger struct {}`
+ `NewRedisLogger(conn redis.Conn, logger *log.Logger) *RedisLogger`
