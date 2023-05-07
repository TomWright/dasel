package dencoding

// NewMap returns a new *Map that has its values initialised.
func NewMap() *Map {
	keys := make([]string, 0)
	data := make(map[string]any)
	return &Map{
		keys: keys,
		data: data,
	}
}

// FromMap creates a *Map from the input map.
// Note that while the contents will be ordered, the ordering is not
// guaranteed since the input map is unordered.
func FromMap(source map[string]any) *Map {
	m := NewMap()
	for k, v := range source {
		m.Set(k, v)
	}
	return m
}

// Map is a map implementation that maintains ordering of keys.
type Map struct {
	// keys contains the keys within the map in the order they were added.
	keys []string
	// data contains the actual map data.
	data map[string]any
}

// Get returns the value found under the given key.
func (m *Map) Get(key string) (any, bool) {
	v, ok := m.data[key]
	return v, ok
}

// Set sets a value under the given key.
func (m *Map) Set(key string, value any) *Map {
	if _, ok := m.data[key]; ok {
		m.data[key] = value
	} else {
		m.keys = append(m.keys, key)
		m.data[key] = value
	}
	return m
}

// Delete deletes the value found under the given key.
func (m *Map) Delete(key string) *Map {
	// Delete the data entry.
	delete(m.data, key)

	// Remove the item from the keys.
	foundIndex := -1
	for i, k := range m.keys {
		if k == key {
			foundIndex = i
			break
		}
	}

	if foundIndex >= 0 {
		m.keys = append((m.keys)[:foundIndex], (m.keys)[foundIndex+1:]...)
	}

	return m
}

// KeyValues returns the KeyValue pairs within the map, in the order that they were added.
func (m *Map) KeyValues() []KeyValue {
	res := make([]KeyValue, 0)
	for _, key := range m.keys {
		res = append(res, KeyValue{
			Key:   key,
			Value: m.data[key],
		})
	}
	return res
}

// Keys returns all the keys from the map.
func (m *Map) Keys() []string {
	return m.keys
}

// UnorderedData returns all the data from the map.
func (m *Map) UnorderedData() map[string]any {
	return m.data
}
