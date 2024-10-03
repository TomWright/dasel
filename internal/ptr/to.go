package ptr

// To returns a pointer to the value passed in.
func To[T any](v T) *T {
	return &v
}
