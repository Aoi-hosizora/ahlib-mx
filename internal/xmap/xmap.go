package xmap

import (
	"fmt"
)

func SliceToStringMap(args []interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	l := len(args)
	for i := 0; i < l; i += 2 {
		keyIdx := i
		valueIdx := i + 1
		if i+1 >= l {
			break
		}

		keyItf := args[keyIdx]
		value := args[valueIdx]
		key := ""
		if keyItf == nil || value == nil {
			continue
		}
		if k, ok := keyItf.(string); ok {
			key = k
		} else {
			key = fmt.Sprintf("%v", keyItf)
		}
		out[key] = value
	}

	return out
}
