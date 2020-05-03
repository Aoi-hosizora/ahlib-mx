package xredis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type RedisLogger struct {
	conn   redis.Conn
	logger *logrus.Logger
}

func NewRedisLogger(conn redis.Conn, logger *logrus.Logger) *RedisLogger {
	return &RedisLogger{conn: conn, logger: logger}
}

func (r *RedisLogger) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := r.conn.Do(commandName, args...)
	r.print(reply, err, commandName, args...)
	return reply, err
}

func (r *RedisLogger) print(reply interface{}, err error, commandName string, v ...interface{}) {
	cmd := r.render(commandName, v)
	if err == nil {
		r.logger.WithFields(logrus.Fields{
			"Module":  "redis",
			"Command": cmd,
			"Error":   err,
		}).Info(fmt.Sprintf("[Redis] return: %8T | %s", reply, cmd))
	} else {
		r.logger.WithFields(logrus.Fields{
			"Module":  "redis",
			"Command": cmd,
			"Error":   err,
		}).Error(fmt.Sprintf("[Redis] error: %v | %s", err, cmd))
	}
}

func (r *RedisLogger) render(cmd string, args []interface{}) string {
	out := cmd
	for _, arg := range args {
		out += " " + fmt.Sprintf("%v", arg)
	}
	return out
}

func (r *RedisLogger) Close() error {
	return r.conn.Close()
}

func (r *RedisLogger) Err() error {
	return r.conn.Err()
}

func (r *RedisLogger) Send(commandName string, args ...interface{}) error {
	return r.conn.Send(commandName, args...)
}

func (r *RedisLogger) Flush() error {
	return r.conn.Flush()
}

func (r *RedisLogger) Receive() (reply interface{}, err error) {
	return r.conn.Receive()
}
