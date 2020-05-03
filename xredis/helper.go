package xredis

import (
	"github.com/gomodule/redigo/redis"
)

func DeleteAll(conn redis.Conn, pattern string) (total int, del int, err error) {
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return 0, 0, err
	}

	cnt := 0
	var someErr error
	for _, key := range keys {
		result, err := redis.Int(conn.Do("DEL", key))
		if err == nil {
			cnt += result
		} else {
			someErr = err
		}
	}
	return len(keys), cnt, someErr
}
