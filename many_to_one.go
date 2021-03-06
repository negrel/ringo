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

func (mto *manyToOne) Push(data Generic) (overwrite bool) {
	head := atomic.AddUint64(&mto.head, 1)
	index := head % mto.capacity

	pBox := box{
		index: head,
		data:  data,
	}

	tail1 := atomic.LoadUint64(&mto.tail)
	old := (*box)(atomic.SwapPointer(&mto.buffer[index], Generic(&pBox)))
	tail2 := atomic.LoadUint64(&mto.tail)

	if old == nil {
		return false
	}

	return old.index > tail1 || old.index > tail2
}

func (mto *manyToOne) Shift() (Generic, bool) {
	index := mto.tail % mto.capacity
	box := (*box)(mto.buffer[index])

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
		return nil, false
	}

	mto.tail++
	return box.data, true
}
