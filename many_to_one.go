package ringo

import (
	"sync/atomic"
)

type manyToOne buffer

// ManyToOne return an efficient ring buffer with the given
// capacity. The buffer is safe for one reader and multiple writer.
// The ManyToOne buffer will panic if you use the -race flag
// because you must use only one reader, no runtime check is performed for
// better performance.
func ManyToOne(capacity uint32) Buffer {
	return &manyToOne{
		head:     ^uint64(0),
		buffer:   make([]Generic, capacity),
		capacity: uint64(capacity),
	}
}

func (mto *manyToOne) Cap() uint32 {
	return uint32(mto.capacity)
}

// Push the given data to the buffer.
func (mto *manyToOne) Push(data Generic) {
	head := atomic.AddUint64(&mto.head, 1)
	index := head % mto.capacity

	box := box{
		index: head,
		data:  data,
	}

	atomic.SwapPointer(&mto.buffer[index], Generic(&box))
}

// Push the given data to the buffer and return if
// the data is valid.
func (mto *manyToOne) Shift() (Generic, bool) {
	i := mto.tail % mto.capacity

	box := (*box)(mto.buffer[i])

	// never written before.
	if box == nil {
		return nil, false
	}

	// already readed
	if box.index < mto.tail {
		return nil, false
	}

	// cell have been overwritten
	if box.index > mto.tail {
		// set the tail so next shift box.data will be valid
		mto.tail = box.index
		return nil, true
	}

	mto.tail++
	return box.data, true
}
