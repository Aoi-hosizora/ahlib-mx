package internal

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/sirupsen/logrus"
	"strings"
)

// loggerOptions is a type of some logger functions' option, each field can be set by Option function type.
type loggerOptions struct {
	text   string                 // extra text
	fields map[string]interface{} // extra fields
}

// LoggerOption represents an option type for some logger functions' option, can be created by WithXXX functions.
type LoggerOption func(*loggerOptions)

// WithExtraText creates a LoggerOption to specific extra text logging in "... | extra_text" style, notes that if you use this multiple times, only the last one will be retained.
func WithExtraText(text string) LoggerOption {
	return func(extra *loggerOptions) {
		extra.text = strings.TrimSpace(text)
	}
}

// WithExtraFields creates a LoggerOption to specific logging with extra fields, notes that if you use this multiple times, only the last one will be retained.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return func(extra *loggerOptions) {
		extra.fields = fields
	}
}

// WithExtraFieldsV creates a LoggerOption to specific logging with extra fields in variadic, notes that if you use this multiple times, only the last one will be retained.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return func(extra *loggerOptions) {
		extra.fields = xstring.SliceToStringMap(fields)
	}
}

// BuildLoggerOptions creates a loggerOptions with given LoggerOption-s.
func BuildLoggerOptions(options []LoggerOption) *loggerOptions {
	out := &loggerOptions{text: "", fields: make(map[string]interface{})}
	for _, op := range options {
		if op != nil {
			op(out)
		}
	}
	return out
}

// ApplyToMessage adds extra string to given message.
func (l *loggerOptions) ApplyToMessage(m *string) {
	if l.text != "" {
		*m += fmt.Sprintf(" | %s", l.text)
	}
}

// ApplyToFields adds extra fields to given logrus.Fields.
func (l *loggerOptions) ApplyToFields(f logrus.Fields) {
	for k, v := range l.fields {
		f[k] = v
	}
}
