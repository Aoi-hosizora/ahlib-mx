package xredis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

// type redis.Conn interface{xxx}
type RedisConn struct {
	conn    redis.Conn
	logger  *logrus.Logger
	LogMode bool
}

func NewRedisConnWithLogger(conn redis.Conn, logger *logrus.Logger) *RedisConn {
	return &RedisConn{conn: conn, logger: logger}
}

func (r *RedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := r.conn.Do(commandName, args...)
	if r.LogMode {
		r.print(reply, err, commandName, args...)
	}
	return reply, err
}

func (r *RedisConn) print(reply interface{}, err error, commandName string, v ...interface{}) {
	cmd := r.render(commandName, v)
	field := r.logger.WithFields(logrus.Fields{
		"module":  "redis",
		"command": cmd,
		"error":   err,
	})

	if err == nil {
		field.Info(fmt.Sprintf("[Redis] return: %8T | %s", reply, cmd))
	} else {
		field.Error(fmt.Sprintf("[Redis] error: %v | %s", err, cmd))
	}
}

func (r *RedisConn) render(cmd string, args []interface{}) string {
	out := cmd
	for _, arg := range args {
		out += " " + fmt.Sprintf("%v", arg)
	}
	return out
}

func (r *RedisConn) Close() error {
	return r.conn.Close()
}

func (r *RedisConn) Err() error {
	return r.conn.Err()
}

func (r *RedisConn) Send(commandName string, args ...interface{}) error {
	return r.conn.Send(commandName, args...)
}

func (r *RedisConn) Flush() error {
	return r.conn.Flush()
}

func (r *RedisConn) Receive() (reply interface{}, err error) {
	return r.conn.Receive()
}
