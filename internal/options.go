package internal

import (
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/sirupsen/logrus"
)

// loggerOptions is a type of some logger functions' option, each field can be set by LoggerOption function type.
type loggerOptions struct {
	text   string
	fields map[string]interface{}
}

// LoggerOption represents an option type for some logger functions' option, can be created by WithXXX functions.
type LoggerOption func(*loggerOptions)

// WithExtraText creates a LoggerOption to specify extra text logging in "...extra_text" style. Note that if you use this multiple times, only the last one will be retained.
func WithExtraText(text string) LoggerOption {
	return func(extra *loggerOptions) {
		extra.text = text // no trim
	}
}

// WithExtraFields creates a LoggerOption to specify logging with extra fields. Note that if you use this multiple times, only the last one will be retained.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return func(extra *loggerOptions) {
		extra.fields = fields
	}
}

// WithExtraFieldsV creates a LoggerOption to specify logging with extra fields in variadic. Note that if you use this multiple times, only the last one will be retained.
func WithExtraFieldsV(fields ...interface{}) LoggerOption {
	return func(extra *loggerOptions) {
		extra.fields = xstring.SliceToStringMap(fields)
	}
}

// BuildLoggerOptions creates a loggerOptions with given LoggerOption-s.
func BuildLoggerOptions(options []LoggerOption) *loggerOptions {
	opt := &loggerOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.fields == nil {
		opt.fields = make(map[string]interface{})
	}
	return opt
}

// ApplyToMessage adds extra string to given message.
func (l *loggerOptions) ApplyToMessage(m *string) {
	if l.text != "" {
		*m += l.text
	}
}

// ApplyToFields adds extra fields to given logrus.Fields.
func (l *loggerOptions) ApplyToFields(f logrus.Fields) {
	for k, v := range l.fields {
		f[k] = v
	}
}
