package xrecovery

import (
	"fmt"
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

	LogToLogrus(l1, fmt.Errorf("test error"))
	LogToLogrus(l1, fmt.Errorf("test error"), WithExtraText("123"))
	LogToLogrus(l1, fmt.Errorf("test error"), WithExtraFields(map[string]interface{}{"a": "b"}))
	LogToLogrus(l1, fmt.Errorf("test error"), WithExtraFieldsV("a", "b"))
	LogToLogrus(l1, fmt.Errorf("test error"), WithExtraText("123"), WithExtraFields(map[string]interface{}{"a": "b"}))

	LogToLogger(l2, fmt.Errorf("test error"))
	LogToLogger(l2, fmt.Errorf("test error"), WithExtraText("123"))
	LogToLogger(l2, fmt.Errorf("test error"), WithExtraFields(map[string]interface{}{"a": "b"}))
}
