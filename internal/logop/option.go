package logop

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// LoggerOption represents an option for logop.LoggerExtraOptions, created by WithXXX functions.
type LoggerOption func(*LoggerExtraOptions)

// WithExtraText creates an logop.LoggerOption for logging extra string.
func WithExtraText(text string) LoggerOption {
	return func(op *LoggerExtraOptions) {
		op.Text = &text
	}
}

// WithExtraFields creates an logop.LoggerOption for logging extra fields.
func WithExtraFields(fields map[string]interface{}) LoggerOption {
	return func(op *LoggerExtraOptions) {
		op.Fields = &fields
	}
}

// WithExtraFieldsV creates an logop.LoggerOption for logging extra fields in vararg.
func WithExtraFieldsV(field ...interface{}) LoggerOption {
	return func(op *LoggerExtraOptions) {
		m := sliceToStringMap(field)
		op.Fields = &m
	}
}

// LoggerExtraOptions represents some extra options for logger, set by LoggerOption.
type LoggerExtraOptions struct {
	Text   *string                 // extra text
	Fields *map[string]interface{} // extra fields
}

// ApplyOptions applies LoggerOption array to loggerExtra.
func (l *LoggerExtraOptions) ApplyOptions(options []LoggerOption) {
	for _, op := range options {
		if op != nil {
			op(l)
		}
	}
}

// AddToMessage adds extra string to message.
func (l *LoggerExtraOptions) AddToMessage(m *string) {
	if l.Text != nil {
		*m += fmt.Sprintf(" | %s", *l.Text)
	}
}

// AddToFields adds extra fields to logrus.Fields.
func (l *LoggerExtraOptions) AddToFields(f logrus.Fields) {
	if l.Fields != nil {
		for k, v := range *l.Fields {
			f[k] = v
		}
	}
}

func sliceToStringMap(args []interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	l := len(args)
	for i := 0; i < l; i += 2 {
		keyIdx := i
		valueIdx := i + 1
		if i+1 >= l {
			break
		}

		keyItf := args[keyIdx]
		value := args[valueIdx]
		key := ""
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
