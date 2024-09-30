package model

import "fmt"

type MapKeyNotFound struct {
	Key string
}

func (e *MapKeyNotFound) Error() string {
	return fmt.Sprintf("map key not found: %q", e.Key)
}
