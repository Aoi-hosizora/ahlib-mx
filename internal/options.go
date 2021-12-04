package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// loggerOptions represents some logger options, such as for extra data, set by LoggerOption.
type loggerOptions struct {
	Text   string                 // extra text
	Fields map[string]interface{} // extra fields
}

// LoggerOption represents an option for loggerOptions, created by WithXXX functions.
type LoggerOption func(*loggerOptions)

// WithExtraText creates a logopt.LoggerOption to log with extra text.
func WithExtraText(text string) LoggerOption {
	return func(extra *loggerOptions) {
		extra.Text = strings.TrimSpace(text)
	}
}

// WithExtraFields creates a logopt.LoggerOption to log with extra fields.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return func(extra *loggerOptions) {
		extra.Fields = fields
	}
}

// WithExtraFieldsV creates a logopt.LoggerOption to log with extra fields in vararg.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return func(extra *loggerOptions) {
		extra.Fields = sliceToMap(fields)
	}
}

// NewLoggerOptions creates a loggerOptions from given LoggerOption-s.
func NewLoggerOptions(options []LoggerOption) *loggerOptions {
	out := &loggerOptions{
		Text:   "",
		Fields: make(map[string]interface{}),
	}
	for _, op := range options {
		if op != nil {
			op(out)
		}
	}
	return out
}

// AddToMessage adds extra string to message.
func (l *loggerOptions) AddToMessage(m *string) {
	if l.Text != "" {
		*m += fmt.Sprintf(" | %s", l.Text)
	}
}

// AddToFields adds extra fields to logrus.Fields.
func (l *loggerOptions) AddToFields(f logrus.Fields) {
	for k, v := range l.Fields {
		f[k] = v
	}
}

// sliceToMap returns a string-interface{} map from interface{} slice.
func sliceToMap(args []interface{}) map[string]interface{} {
	l := len(args)
	out := make(map[string]interface{}, l/2)

	for i := 0; i < l; i += 2 {
		if i+1 >= l {
			break // ignore the final arg
		}
		key := ""
		keyItf, value := args[i], args[i+1]
		if keyItf == nil {
			i--
			continue
		}
		if k, ok := keyItf.(string); ok {
			key = k
		} else {
			key = fmt.Sprintf("%v", keyItf) // %v
		}
		out[key] = value
	}

	return out
}
