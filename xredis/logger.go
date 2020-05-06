package xredis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"log"
)

// logrus.Logger

type RedisLogrus struct {
	redis.Conn
	logger  *logrus.Logger
	LogMode bool
}

func NewRedisLogrus(conn redis.Conn, logger *logrus.Logger) *RedisLogrus {
	return &RedisLogrus{Conn: conn, logger: logger}
}

func (r *RedisLogrus) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := r.Conn.Do(commandName, args...)
	if r.LogMode {
		r.print(reply, err, commandName, args...)
	}
	return reply, err
}

func (r *RedisLogrus) print(reply interface{}, err error, commandName string, v ...interface{}) {
	cmd := render(commandName, v)
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

// logrus.Logger

type RedisLogger struct {
	redis.Conn
	logger  *log.Logger
	LogMode bool
}

func NewRedisLogger(conn redis.Conn, logger *log.Logger) *RedisLogger {
	return &RedisLogger{Conn: conn, logger: logger}
}

func (r *RedisLogger) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := r.Conn.Do(commandName, args...)
	if r.LogMode {
		r.print(reply, err, commandName, args...)
	}
	return reply, err
}

func (r *RedisLogger) print(reply interface{}, err error, commandName string, v ...interface{}) {
	cmd := render(commandName, v)
	if err == nil {
		r.logger.Printf("[Redis] return: %8T | %s", reply, cmd)
	} else {
		r.logger.Printf("[Redis] error: %v | %s", err, cmd)
	}
}

// render

func render(cmd string, args []interface{}) string {
	out := cmd
	for _, arg := range args {
		out += " " + fmt.Sprintf("%v", arg)
	}
	return out
}
