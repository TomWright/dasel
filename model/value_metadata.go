package model

// MetadataValue returns a metadata value.
func (v *Value) MetadataValue(key string) (any, bool) {
	if v.Metadata == nil {
		return nil, false
	}
	val, ok := v.Metadata[key]
	return val, ok
}

// SetMetadataValue sets a metadata value.
func (v *Value) SetMetadataValue(key string, val any) {
	if v.Metadata == nil {
		v.Metadata = map[string]any{}
	}
	v.Metadata[key] = val
}

// IsSpread returns true if the value is a spread value.
// Spread values are used to represent the spread operator.
func (v *Value) IsSpread() bool {
	if v == nil {
		return false
	}
	val, ok := v.MetadataValue("spread")
	if !ok {
		return false
	}
	spread, ok := val.(bool)
	return ok && spread
}

// MarkAsSpread marks the value as a spread value.
// Spread values are used to represent the spread operator.
func (v *Value) MarkAsSpread() {
	v.SetMetadataValue("spread", true)
}

// IsBranch returns true if the value is a branched value.
func (v *Value) IsBranch() bool {
	if v == nil {
		return false
	}
	val, ok := v.MetadataValue("branch")
	if !ok {
		return false
	}
	branch, ok := val.(bool)
	return ok && branch
}

// MarkAsBranch marks the value as a branch value.
func (v *Value) MarkAsBranch() {
	v.SetMetadataValue("branch", true)
}
