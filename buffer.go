package ringo

// Buffer define common methods for ring buffers.
type Buffer interface {
	// The capacity of the buffer.
	Cap() uint32

	Push(Generic)
	Shift() (Generic, bool)
}

type buffer struct {
	buffer []Generic
	head   uint64
	tail   uint64
	// Stored as uint64 to avoid conversion
	// but will never overflow uint32
	capacity uint64
}
