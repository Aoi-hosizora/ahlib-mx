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

	dto = BuildBasicErrorDto("test error", []string{""})
	xtesting.Equal(t, dto.Type, "string")
	xtesting.Equal(t, dto.Detail, "test error")
	xtesting.Equal(t, dto.Request, []string{""})
	xtesting.Equal(t, dto.Others, map[string]interface{}{})
	xtesting.Equal(t, dto.Filename, "")
	xtesting.Equal(t, dto.Funcname, "")
	xtesting.Equal(t, dto.LineIndex, 0)
	xtesting.Equal(t, dto.Line, "")
	xtesting.Equal(t, dto.Stacks, []string(nil))

	dto = BuildBasicErrorDto(fmt.Errorf("test error"), nil, "test1", "error", "test2", 0, "xxx")
	xtesting.Equal(t, dto.Type, "*errors.errorString")
	xtesting.Equal(t, dto.Detail, "test error")
	xtesting.Equal(t, dto.Request, []string{})
	xtesting.Equal(t, len(dto.Others), 2)
	xtesting.Equal(t, dto.Others["test1"], "error")
	xtesting.Equal(t, dto.Others["test2"], 0)

	dto = BuildErrorDto("", nil, 0, false, 0)
	xtesting.Equal(t, dto.Detail, "")
	xtesting.Equal(t, dto.Funcname, `xdto.TestErrorDto`)
	xtesting.Equal(t, dto.Others, map[string]interface{}{})
	xtesting.True(t, strings.HasSuffix(dto.Filename, `xdto_test.go`))
	xtesting.Equal(t, dto.Line, `dto = BuildErrorDto("", nil, 0, false, 0)`)

	_ = BuildErrorDto("", nil, 0, true)
}
