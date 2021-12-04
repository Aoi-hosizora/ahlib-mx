package xrecovery

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)

	for _, std := range []bool{false, true} {
		for _, tc := range []struct {
			giveErr     interface{}
			giveStack   xruntime.TraceStack
			giveOptions []LoggerOption
		}{
			{nil, nil, nil},
			{"test string", nil, nil},
			{errors.New("test error"), nil, nil},
			{nil, xruntime.RuntimeTraceStack(0), nil},
			{errors.New("test error"), xruntime.RuntimeTraceStack(0), nil},

			{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraText("extra")}},
			{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
			{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraFieldsV("k", "v")}},
			{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraText("extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
			{errors.New("test error"), xruntime.RuntimeTraceStack(0), []LoggerOption{WithExtraText("extra"), WithExtraFieldsV("k", "v")}},
		} {
			if !std {
				LogToLogrus(l1, tc.giveErr, tc.giveStack, tc.giveOptions...)
			} else {
				LogToLogger(l2, tc.giveErr, tc.giveStack, tc.giveOptions...)
			}
		}
	}
}
