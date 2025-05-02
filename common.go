package chartmetric

// Optional is a pointer to a value of type T.
type Optional[T any] *T

// Opt converts a value of type T to an Optional[T].
func Opt[T any](v T) Optional[T] {
	return &v
}
