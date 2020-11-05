package xmap

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"testing"
)

func TestSliceToStringMap(t *testing.T) {
	xtesting.Equal(t, SliceToStringMap(nil), map[string]interface{}{})
	xtesting.Equal(t, SliceToStringMap([]interface{}{"a", "b", "c"}), map[string]interface{}{"a": "b"})
	xtesting.Equal(t, SliceToStringMap([]interface{}{"a", "b", "c", "d"}), map[string]interface{}{"a": "b", "c": "d"})
	xtesting.Equal(t, SliceToStringMap([]interface{}{1, 2}), map[string]interface{}{"1": 2})
	xtesting.Equal(t, SliceToStringMap([]interface{}{nil, "  ", " ", nil, " ", " "}), map[string]interface{}{" ": " "})
}
