package xserverchan

import (
	"github.com/Aoi-hosizora/ahlib-web/internal/logop"
	"github.com/Aoi-hosizora/go-serverchan"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)

	key := "SCU0abcdefg123hijklmn456opqrst789uvwxyz123abcdefg456ab"
	title := strings.Repeat("test", 256/4) // 256B
	body := strings.Repeat("t", 1<<16)     // 64KB

	for _, std := range []bool{false, true} {
		for _, tc := range []struct {
			giveKey     string
			giveTitle   string
			giveBody    string
			giveErr     error
			giveOptions []logop.LoggerOption
		}{
			{"", "", "", nil, nil},
			{"KEY", "title", "body", nil, nil},
			{key, "title", "body", nil, nil},
			{key, title, body, nil, nil},

			{key, title, body, nil, []logop.LoggerOption{WithExtraText("extra")}},
			{key, title, body, nil, []logop.LoggerOption{WithExtraFields(map[string]interface{}{"k": "v"})}},
			{key, title, body, nil, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
			{key, title, body, nil, []logop.LoggerOption{WithExtraText("extra"), WithExtraFields(map[string]interface{}{"k": "v"})}},
			{key, title, body, nil, []logop.LoggerOption{WithExtraText("extra"), WithExtraFieldsV("k", "v")}},

			{"", "", "", serverchan.ErrEmptyTitle, nil},
			{"KEY", "title", "body", serverchan.ErrBadPushToken, nil},
			{key, "title", "body", serverchan.ErrNotSuccess, nil},
			{key, title, body, serverchan.ErrDuplicateMessage, nil},
			{key, title, body, serverchan.ErrDuplicateMessage, []logop.LoggerOption{WithExtraText("extra")}},
			{key, title, body, serverchan.ErrDuplicateMessage, []logop.LoggerOption{WithExtraFieldsV("k", "v")}},
		} {
			if !std {
				LogToLogrus(l1, tc.giveKey, tc.giveTitle, tc.giveBody, tc.giveErr, tc.giveOptions...)
			} else {
				LogToLogger(l2, tc.giveKey, tc.giveTitle, tc.giveBody, tc.giveErr, tc.giveOptions...)
			}
		}
	}
}
