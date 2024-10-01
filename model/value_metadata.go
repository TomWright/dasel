package model

func (v *Value) MetadataValue(key string) (any, bool) {
	if v.Metadata == nil {
		return nil, false
	}
	val, ok := v.Metadata[key]
	return val, ok
}

func (v *Value) SetMetadataValue(key string, val any) {
	if v.Metadata == nil {
		v.Metadata = map[string]any{}
	}
	v.Metadata[key] = val
}

func (v *Value) IsSpread() bool {
	val, ok := v.Metadata["spread"]
	if !ok {
		return false
	}
	spread, ok := val.(bool)
	if !ok {
		return false
	}
	return spread
}

func (v *Value) MarkAsSpread() {
	v.SetMetadataValue("spread", true)
}
