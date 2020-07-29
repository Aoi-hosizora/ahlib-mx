package xredis

import (
	"github.com/gomodule/redigo/redis"
)

type Helper struct {
	conn redis.Conn
}

// noinspection GoUnusedExportedFunction
func WithConn(conn redis.Conn) *Helper {
	return &Helper{conn: conn}
}

func (h *Helper) DeleteAll(pattern string) (total int, del int, err error) {
	keys, err := redis.Strings(h.conn.Do("KEYS", pattern))
	if err != nil {
		return 0, 0, err
	}

	cnt := 0
	var someErr error
	for _, key := range keys {
		result, err := redis.Int(h.conn.Do("DEL", key))
		if err == nil {
			cnt += result
		} else if someErr == nil {
			someErr = err
		}
	}
	return len(keys), cnt, someErr
}
