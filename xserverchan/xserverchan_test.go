package xserverchan

import (
	"github.com/Aoi-hosizora/go-serverchan"
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

	key := "SCU0abcdefg123hijklmn456opqrst789uvwxyz123abcdefg456"
	title := "Test Title"
	LogToLogrus(l1, key, title, nil)
	LogToLogrus(l1, key, title, serverchan.ErrDuplicateMessage)
	LogToLogrus(l1, key, title, serverchan.ErrBadPushToken)
	LogToLogger(l2, key, title, nil)
	LogToLogger(l2, key, title, serverchan.ErrDuplicateMessage)
	LogToLogger(l2, key, title, serverchan.ErrBadPushToken)
}
