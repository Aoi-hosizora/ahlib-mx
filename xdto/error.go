package xdto

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/internal/xmap"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"log"
	"os"
	"reflect"
	"time"
)

// ErrorDto is a general error response model.
type ErrorDto struct {
	Time    string   `json:"time"`    // current time
	Type    string   `json:"type"`    // error type
	Detail  string   `json:"detail"`  // error detail message
	Request []string `json:"request"` // request details

	Others map[string]interface{} `json:"others,omitempty"` // other message

	Filename  string   `json:"filename,omitempty"`   // stack filename
	Funcname  string   `json:"funcname,omitempty"`   // stack function name
	LineIndex int      `json:"line_index,omitempty"` // file line index
	Line      string   `json:"line,omitempty"`       // stack current line
	Stacks    []string `json:"stacks,omitempty"`     // stacks in skip
}

// BuildBasicErrorDto builds a basic dto (only include time, type, detail, request).
func BuildBasicErrorDto(err interface{}, requests []string, otherKvs ...interface{}) *ErrorDto {
	skip := -2
	return BuildErrorDto(err, requests, skip, false, otherKvs...)
}

// BuildErrorDto builds a complete dto (also include runtime parameters).
func BuildErrorDto(err interface{}, requests []string, skip int, doPrint bool, otherKvs ...interface{}) *ErrorDto {
	if err == nil {
		return nil
	}

	// basic
	now := time.Now().Format(time.RFC3339)
	errType := reflect.TypeOf(err).String()
	errDetail := ""
	if e, ok := err.(error); ok {
		errDetail = e.Error()
	} else {
		errDetail = fmt.Sprintf("%v", err)
	}
	if requests == nil {
		requests = []string{}
	}
	dto := &ErrorDto{Time: now, Type: errType, Detail: errDetail, Request: requests}

	// other
	dto.Others = xmap.SliceToStringMap(otherKvs)

	// runtime
	if skip >= 0 {
		skip++
		var stacks xruntime.TraceStack
		stacks, dto.Filename, dto.Funcname, dto.LineIndex, dto.Line = xruntime.RuntimeTraceStackWithInfo(skip)
		dto.Stacks = make([]string, len(stacks))
		for idx, stack := range stacks {
			dto.Stacks[idx] = stack.String()
		}
		if doPrint {
			l := log.New(os.Stderr, "", 0)
			l.Println()
			l.Println(xcolor.Red.Sprint(stacks.String()))
			l.Println()
		}
	}

	return dto
}
