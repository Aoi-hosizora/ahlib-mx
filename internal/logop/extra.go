package logop

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// loggerExtra represents some extra options for logger.
type loggerExtra struct {
	Text   string                 // extra text
	Fields map[string]interface{} // extra fields
}

// LoggerOption represents an option for loggerExtra, created by WithXXX functions.
type LoggerOption func(*loggerExtra)

// WithExtraText creates a logop.LoggerOption for logger to log with extra text.
func WithExtraText(text string) LoggerOption {
	return func(extra *loggerExtra) {
		extra.Text = strings.TrimSpace(text)
	}
}

// WithExtraFields creates a logop.LoggerOption for logger to log with extra fields.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return func(extra *loggerExtra) {
		extra.Fields = fields
	}
}

// WithExtraFieldsV creates a logop.LoggerOption for logger to log with extra fields in vararg.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return func(extra *loggerExtra) {
		extra.Fields = sliceToMap(fields)
	}
}

// NewLoggerExtra creates a loggerExtra from given LoggerOption-s.
func NewLoggerExtra(options ...LoggerOption) *loggerExtra {
	extra := &loggerExtra{}
	for _, op := range options {
		if op != nil {
			op(extra)
		}
	}
	return extra
}

// AddToMessage adds extra string to message.
func (l *loggerExtra) AddToMessage(m *string) {
	if l.Text != "" {
		*m += fmt.Sprintf(" | %s", l.Text)
	}
}

// AddToFields adds extra fields to logrus.Fields.
func (l *loggerExtra) AddToFields(f logrus.Fields) {
	for k, v := range l.Fields {
		f[k] = v
	}
}

// sliceToMap returns a string-interface{} map from interface{} slice.
func sliceToMap(args []interface{}) map[string]interface{} {
	l := len(args)
	out := make(map[string]interface{}, l/2)

	for i := 0; i < l; i += 2 {
		ki := i
		vi := i + 1
		if i+1 >= l {
			break // ignore the final arg
		}

		key := "" // string
		keyItf := args[ki]
		value := args[vi] // interface{}
		if keyItf == nil || value == nil {
			continue
		}
		if k, ok := keyItf.(string); ok {
			key = k
		} else {
			key = fmt.Sprintf("%v", keyItf)
		}

		out[key] = value
	}

	return out
}
