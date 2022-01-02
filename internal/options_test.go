package internal

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLoggerOptions(t *testing.T) {
	for _, tc := range []struct {
		give       []LoggerOption
		wantMsg    string
		wantFields logrus.Fields
	}{
		{[]LoggerOption{}, "", logrus.Fields{}},
		{[]LoggerOption{nil}, "", logrus.Fields{}},
		{[]LoggerOption{nil, nil, nil}, "", logrus.Fields{}},

		{[]LoggerOption{WithExtraText("")}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("   ")}, "   ", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("  x x  ")}, "  x x  ", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("test")}, "test", logrus.Fields{}},
		{[]LoggerOption{WithExtraText("test1"), WithExtraText(" | test2")}, " | test2", logrus.Fields{}},

		{[]LoggerOption{WithExtraFields(map[string]interface{}{})}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"true": 2, "3": 4.4})}, "", logrus.Fields{"true": 2, "3": 4.4}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"true": 2, "3": 4.4}), WithExtraFields(map[string]interface{}{"k": "v"})}, "", logrus.Fields{"k": "v"}},

		{[]LoggerOption{WithExtraFieldsV()}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil)}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil, "a", nil)}, "", logrus.Fields{"a": nil}},
		{[]LoggerOption{WithExtraFieldsV(nil, nil, "a")}, "", logrus.Fields{}},
		{[]LoggerOption{WithExtraFieldsV(nil, nil, 1, nil)}, "", logrus.Fields{"1": nil}},
		{[]LoggerOption{WithExtraFieldsV(true, 2, 3.3, 4, 5)}, "", logrus.Fields{"true": 2, "3.3": 4}},
		{[]LoggerOption{WithExtraFieldsV(true, 2, 3.3, 4, 5), WithExtraFieldsV("k", "v")}, "", logrus.Fields{"k": "v"}},

		{[]LoggerOption{WithExtraText("test"), WithExtraFields(map[string]interface{}{"1": 2})}, "test", logrus.Fields{"1": 2}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"1": 2}), WithExtraText(" | test")}, " | test", logrus.Fields{"1": 2}},
		{[]LoggerOption{WithExtraText("test"), WithExtraFieldsV(3, 4)}, "test", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFieldsV(3, 4), WithExtraText(" | test")}, " | test", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFields(map[string]interface{}{"1": 2}), WithExtraFieldsV(3, 4)}, "", logrus.Fields{"3": 4}},
		{[]LoggerOption{WithExtraFieldsV(3, 4), WithExtraFields(map[string]interface{}{"1": 2})}, "", logrus.Fields{"1": 2}},
	} {
		ops := BuildLoggerOptions(tc.give)
		msg := ""
		fields := logrus.Fields{}
		ops.ApplyToMessage(&msg)
		ops.ApplyToFields(fields)
		xtesting.Equal(t, msg, tc.wantMsg)
		xtesting.Equal(t, fields, tc.wantFields)
	}
}
