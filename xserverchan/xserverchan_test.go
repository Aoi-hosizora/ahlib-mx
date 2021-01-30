package xserverchan

import (
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	key := "SCU123456789abcde123456789abcde123456789abcde123456789"
	LogSendToLogger(l, key, "test title", nil)
}
