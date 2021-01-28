package xrecovery

import (
	"fmt"
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

	LogToLogrus(l1, fmt.Errorf("test error"), nil)
	LogToLogrus(l1, fmt.Errorf("test error"), xruntime.RuntimeTraceStack(0))
	LogToLogrus(l1, fmt.Errorf("test error"), nil, WithExtraText("123"))
	LogToLogrus(l1, fmt.Errorf("test error"), xruntime.RuntimeTraceStack(0), WithExtraText("123"))
	LogToLogrus(l1, fmt.Errorf("test error"), nil, WithExtraFields(map[string]interface{}{"a": "b"}))
	LogToLogrus(l1, fmt.Errorf("test error"), nil, WithExtraFieldsV("a", "b"))
	LogToLogrus(l1, fmt.Errorf("test error"), nil, WithExtraText("123"), WithExtraFields(map[string]interface{}{"a": "b"}))

	LogToLogger(l2, fmt.Errorf("test error"), nil)
	LogToLogger(l2, fmt.Errorf("test error"), xruntime.RuntimeTraceStack(0))
	LogToLogger(l2, fmt.Errorf("test error"), nil, WithExtraText("123"))
	LogToLogger(l2, fmt.Errorf("test error"), xruntime.RuntimeTraceStack(0), WithExtraText("123"))
	LogToLogger(l2, fmt.Errorf("test error"), nil, WithExtraFields(map[string]interface{}{"a": "b"}))
}
