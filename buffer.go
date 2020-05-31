package ringo

// Buffer define common methods for ring buffers.
type Buffer interface {
	// The capacity of the buffer.
	Cap() uint32

	// Push the given data to the buffer.
	// It may overwrite some unreaded data.
	Push(Generic)

	// Shift data from buffer and return if
	// the data is valid.
	Shift() (Generic, bool)
}

type buffer struct {
	buffer []Generic
	head   uint64
	tail   uint64
	// NOTE Stored as uint64 to avoid conversion
	// but will never overflow uint32
	capacity uint64
}
