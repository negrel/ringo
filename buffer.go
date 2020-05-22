package ringo

import "sync/atomic"

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

// Push the given data to the buffer, block if the
func (mto *manyToOne) Push(data Generic) {
	index := mto.head % uint64(mto.capacity)

	box := box{
		index: mto.head,
		data:  data,
	}

	atomic.AddUint64(&mto.head, 1)

	mto.buffer[index] = Generic(&box)
}

// Push the given data to the buffer and return if
// the data is valid.
func (mto *manyToOne) Shift() (Generic, bool) {
	i := mto.tail % mto.capacity
	mto.tail++

	box := (*box)(mto.buffer[i])

	if box == nil {
		return nil, false
	}

	return box.data, box.index > mto.tail
}
