package ptr

func To[T any](v T) *T {
	return &v
}
