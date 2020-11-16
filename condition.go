package dasel

import (
	"reflect"
)

// Condition defines a Check we can use within dynamic selectors.
type Condition interface {
	Check(other reflect.Value) (bool, error)
}
