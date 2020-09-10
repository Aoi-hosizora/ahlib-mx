package xdto

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"time"
)

// An error response model for fiber and gin.
// Request: Need to dump gin.Request or fiber.Fasthttp.
// Filename... need for runtime stack, need to provide skip.
type ErrorDto struct {
	Time    string                 `json:"time"`    // current time
	Type    string                 `json:"type"`    // error type
	Detail  string                 `json:"detail"`  // error detail message
	Request []string               `json:"request"` // request details
	Others  map[string]interface{} `json:"others"`  // other message

	Filename  string   `json:"filename,omitempty"`   // stack filename
	Funcname  string   `json:"funcname,omitempty"`   // stack function name
	LineIndex int      `json:"line_index,omitempty"` // file line index
	Line      string   `json:"line,omitempty"`       // stack current line
	Stacks    []string `json:"stacks,omitempty"`     // stacks in skip
}

// Build a basic dto (only include time, type, detail, request).
func BuildBasicErrorDto(err interface{}, requests []string, others map[string]interface{}) *ErrorDto {
	skip := -2
	return BuildErrorDto(err, requests, others, skip, false)
}

// Build a complete dto (also include runtime parameters).
func BuildErrorDto(err interface{}, requests []string, others map[string]interface{}, skip int, print bool) *ErrorDto {
	skip++

	now := time.Now().Format(time.RFC3339)
	errType := fmt.Sprintf("%T", err)
	errDetail := fmt.Sprintf("%v", err)
	if e, ok := err.(error); ok {
		errDetail = e.Error()
	}
	if requests == nil {
		requests = []string{}
	}
	if others == nil {
		others = map[string]interface{}{}
	}

	dto := &ErrorDto{
		Time:    now,
		Type:    errType,
		Detail:  errDetail,
		Request: requests,
		Others:  others,
	}

	if skip >= 0 {
		var stacks []*xruntime.Stack
		stacks, dto.Filename, dto.Funcname, dto.LineIndex, dto.Line = xruntime.GetStackWithInfo(skip)
		dto.Stacks = make([]string, len(stacks))
		for idx, stack := range stacks {
			dto.Stacks[idx] = stack.String()
		}
		if print {
			fmt.Println()
			xruntime.PrintStacksRed(stacks)
			fmt.Println()
		}
	}

	return dto
}
