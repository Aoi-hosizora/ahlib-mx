package xgorm

import (
	"database/sql/driver"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"regexp"
	"time"
)

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// logrus.Logger

type GormLogrus struct {
	logger *logrus.Logger
}

// noinspection GoUnusedExportedFunction
func NewGormLogrus(logger *logrus.Logger) *GormLogrus {
	return &GormLogrus{logger: logger}
}

func (g *GormLogrus) Print(v ...interface{}) {
	if len(v) == 0 || len(v) == 1 {
		g.logger.WithFields(logrus.Fields{
			"module": "gorm",
		}).Error(fmt.Sprintf("[Gorm] Unknown message: %v", v))
		return
	}

	if level := v[0]; level == "info" {
		info := v[1]
		g.logger.WithFields(logrus.Fields{
			"module": "gorm",
			"type":   "info",
			"info":   info,
		}).Info(fmt.Sprintf("[Gorm] info: %s", info))
	} else if level == "sql" {
		source := v[1]
		duration := v[2]
		sql := render(v[3].(string), v[4])
		rows := v[5]
		g.logger.WithFields(logrus.Fields{
			"module":   "gorm",
			"type":     "sql",
			"source":   source,
			"duration": duration,
			"aql":      sql,
			"rows":     rows,
		}).Info(fmt.Sprintf("[Gorm] rows: %3d | %10s | %s | %s", rows, duration, sql, source))
	} else {
		g.logger.WithFields(logrus.Fields{
			"module": "gorm",
			"type":   level,
		}).Info(fmt.Sprintf("[Gorm] unknown level %s: %v", level, v))
	}
}

// log.Logger

type GormLogger struct {
	logger *log.Logger
}

// noinspection GoUnusedExportedFunction
func NewGormLogger(logger *log.Logger) *GormLogger {
	return &GormLogger{logger: logger}
}

func (g *GormLogger) Print(v ...interface{}) {
	if len(v) == 0 || len(v) == 1 {
		g.logger.Printf("[Gorm] Unknown message: %v", v)
		return
	}

	if level := v[0]; level == "info" {
		info := v[1]
		g.logger.Printf("[Gorm] info: %s", info)
	} else if level == "sql" {
		source := v[1]
		duration := v[2]
		sql := render(v[3].(string), v[4])
		rows := v[5]
		g.logger.Printf("[Gorm] rows: %3d | %10s | %s | %s", rows, duration, sql, source)
	} else {
		g.logger.Printf("[Gorm] unknown level %s: %v", level, v)
	}
}

// render

func render(sql string, param interface{}) string {
	values := make([]interface{}, 0)
	for _, value := range param.([]interface{}) {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() { // valid
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok { // time
				values = append(values, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
			} else if b, ok := value.([]byte); ok { // bytes
				values = append(values, fmt.Sprintf("'%v'", string(b)))
			} else if r, ok := value.(driver.Valuer); ok { // driver
				if value, err := r.Value(); err == nil && value != nil {
					values = append(values, fmt.Sprintf("'%v'", value))
				} else {
					values = append(values, "NULL")
				}
			} else { // other value
				values = append(values, fmt.Sprintf("'%v'", value))
			}
		} else { // invalid
			values = append(values, fmt.Sprintf("'%v'", value))
		}
	}

	return fmt.Sprintf(sqlRegexp.ReplaceAllString(sql, "%v"), values...)
}
