package xrecovery

import (
	"fmt"
	logrus2 "github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logrus := logrus2.New()

	WithLogrus(logrus, fmt.Errorf("test error"))
	WithLogrus(logrus, fmt.Errorf("test error"), WithExtraString("123"))
	WithLogrus(logrus, fmt.Errorf("test error"), WithExtraFields(map[string]interface{}{"a": "b"}))
	WithLogrus(logrus, fmt.Errorf("test error"), WithExtraString("123"), WithExtraFields(map[string]interface{}{"a": "b"}))

	WithLogger(logger, fmt.Errorf("test error"))
	WithLogger(logger, fmt.Errorf("test error"), WithExtraString("123"))
	WithLogger(logger, fmt.Errorf("test error"), WithExtraFields(map[string]interface{}{"a": "b"}))
}
