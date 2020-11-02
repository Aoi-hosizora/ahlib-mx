package xwlogger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// LoggerOption represents some options for WithLogrus and WithLogger.
type LoggerOption func(*LoggerExtra)

// LoggerExtra represents an extra option for logger, modified by LoggerOption.
type LoggerExtra struct {
	ExtraString *string
	ExtraFields *map[string]interface{}
}

// ApplyOptions applies LoggerOption array to loggerExtra.
func (l *LoggerExtra) ApplyOptions(options []LoggerOption) {
	if len(options) > 0 {
		for _, op := range options {
			if op != nil {
				op(l)
			}
		}
	}
}

// AddToString adds extraString to message.
func (l *LoggerExtra) AddToString(m *string) {
	if l.ExtraString != nil {
		*m += fmt.Sprintf(" | %s", *l.ExtraString)
	}
}

// AddToFields adds extraFields to fields.
func (l *LoggerExtra) AddToFields(f logrus.Fields) {
	if l.ExtraFields != nil {
		for k, v := range *l.ExtraFields {
			f[k] = v
		}
	}
}
