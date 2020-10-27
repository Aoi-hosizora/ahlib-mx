package xdto

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"strings"
	"testing"
)

func TestErrorDto(t *testing.T) {
	dto := BuildBasicErrorDto(nil, nil, nil)
	xtesting.Nil(t, dto)

	dto = BuildBasicErrorDto("test error", []string{""}, nil)
	xtesting.Equal(t, dto.Type, "string")
	xtesting.Equal(t, dto.Detail, "test error")
	xtesting.Equal(t, dto.Request, []string{""})
	xtesting.Equal(t, dto.Others, map[string]interface{}{})
	xtesting.Equal(t, dto.Filename, "")
	xtesting.Equal(t, dto.Funcname, "")
	xtesting.Equal(t, dto.LineIndex, 0)
	xtesting.Equal(t, dto.Line, "")
	xtesting.Equal(t, dto.Stacks, []string(nil))

	dto = BuildBasicErrorDto(fmt.Errorf("test error"), nil, map[string]interface{}{"test": "error"})
	xtesting.Equal(t, dto.Type, "*errors.errorString")
	xtesting.Equal(t, dto.Detail, "test error")
	xtesting.Equal(t, dto.Request, []string{})
	xtesting.Equal(t, dto.Others["test"], "error")

	dto = BuildErrorDto("", nil, nil, 0, false)
	xtesting.Equal(t, dto.Detail, "")
	xtesting.Equal(t, dto.Funcname, `xdto.TestErrorDto`)
	xtesting.True(t, strings.HasSuffix(dto.Filename, `xdto_test.go`))
	xtesting.Equal(t, dto.Line, `dto = BuildErrorDto("", nil, nil, 0, false)`)

	_ = BuildErrorDto("", nil, nil, 0, true)
}
