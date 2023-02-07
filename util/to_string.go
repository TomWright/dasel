package util

import (
	"fmt"
)

// ToString converts the given value to a string.
func ToString(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		return v
	case []byte:
		return string(v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprint(v)
	default:
		return fmt.Sprint(v)
	}
}
