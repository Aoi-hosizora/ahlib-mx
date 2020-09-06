package xdto

import (
	"fmt"
	"log"
	"testing"
)

func TestErrorDto(t *testing.T) {
	log.Println(BuildBasicErrorDto(fmt.Errorf("test error"), []string{"test"}, map[string]interface{}{"a": 1, "b": 2}))
	log.Println()
	log.Println(BuildErrorDto(fmt.Errorf("test error"), nil, nil, 1, true))
}
