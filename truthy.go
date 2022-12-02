package dasel

import (
	"reflect"
	"strings"
)

func IsTruthy(value interface{}) bool {
	switch v := value.(type) {
	case Value:
		return IsTruthy(v.Unpack().Interface())

	case reflect.Value:
		return IsTruthy(unpackReflectValue(v).Interface())

	case bool:
		return v

	case string:
		v = strings.ToLower(strings.TrimSpace(v))
		switch v {
		case "false", "no", "0":
			return false
		default:
			return v != ""
		}

	case []byte:
		return IsTruthy(string(v))

	case int:
		return v > 0
	case int8:
		return v > 0
	case int16:
		return v > 0
	case int32:
		return v > 0
	case int64:
		return v > 0

	case uint:
		return v > 0
	case uint8:
		return v > 0
	case uint16:
		return v > 0
	case uint32:
		return v > 0
	case uint64:
		return v > 0

	case float32:
		return v >= 1
	case float64:
		return v >= 1

	default:
		return false
	}
}
