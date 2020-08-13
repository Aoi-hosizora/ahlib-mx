package xredis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"strings"
	"time"
)

// logrus.Logger

type RedisLogrus struct {
	redis.Conn
	logger  *logrus.Logger
	LogMode bool
}

func NewRedisLogrus(conn redis.Conn, logger *logrus.Logger, logMode bool) *RedisLogrus {
	return &RedisLogrus{Conn: conn, logger: logger, LogMode: logMode}
}

func (r *RedisLogrus) Do(commandName string, args ...interface{}) (interface{}, error) {
	s := time.Now()
	reply, err := r.Conn.Do(commandName, args...)
	e := time.Now()
	if r.LogMode {
		r.print(reply, err, commandName, e.Sub(s).String(), args...)
	}
	return reply, err
}

func (r *RedisLogrus) print(reply interface{}, err error, commandName string, du string, v ...interface{}) {
	cmd := renderCommand(commandName, v)

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"module":  "redis",
			"command": cmd,
			"error":   err,
		}).Error(fmt.Sprintf("[Redis] %v | %s", err, cmd))
		return
	}

	cnt, t := renderReply(reply)
	r.logger.WithFields(logrus.Fields{
		"module":   "redis",
		"command":  cmd,
		"count":    cnt,
		"duration": du,
	}).Info(fmt.Sprintf("[Redis] #: %2d | %10s | %12s | %s", cnt, du, t, cmd))
}

// logrus.Logger

type RedisLogger struct {
	redis.Conn
	logger  *log.Logger
	LogMode bool
}

func NewRedisLogger(conn redis.Conn, logger *log.Logger, logMode bool) *RedisLogger {
	return &RedisLogger{Conn: conn, logger: logger, LogMode: logMode}
}

func (r *RedisLogger) Do(commandName string, args ...interface{}) (interface{}, error) {
	s := time.Now()
	reply, err := r.Conn.Do(commandName, args...)
	e := time.Now()
	if r.LogMode {
		r.print(reply, err, commandName, e.Sub(s).String(), args...)
	}
	return reply, err
}

func (r *RedisLogger) print(reply interface{}, err error, commandName string, du string, v ...interface{}) {
	cmd := renderCommand(commandName, v)

	if err != nil {
		r.logger.Printf("[Redis] %v | %s", err, cmd)
		return
	}

	cnt, t := renderReply(reply)
	r.logger.Printf("[Redis] #: %2d | %10s | %12s | %s", cnt, du, t, cmd)
}

// render

func renderCommand(cmd string, args []interface{}) string {
	out := cmd
	for _, arg := range args {
		out += " " + fmt.Sprintf("%v", arg)
	}
	return strings.TrimSpace(out)
}

func renderReply(reply interface{}) (cnt int, t string) {
	if reply == nil {
		cnt = 0
		t = "<nil>"
	} else if val := reflect.ValueOf(reply); val.Kind() == reflect.Slice && val.IsValid() {
		cnt = val.Len()
		t = val.Type().Elem().String()
		if t == "uint8" { // byte
			cnt = 1
			t = "string"
		} else if t == "interface {}" && val.Len() >= 1 {
			t = reflect.TypeOf(val.Index(0).Interface()).String()
			if t == "[]uint8" { // string
				t = "string"
			}
		}
	} else {
		cnt = 1
		t = fmt.Sprintf("%T", reply)
	}
	if reply == "OK" {
		t = "string (OK)"
	}
	return
}
