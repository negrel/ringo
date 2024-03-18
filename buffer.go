package ringo

// Box wraps a T to store its index within a ring buffer.
type box[T any] struct {
	index uint64
	data  T
}
