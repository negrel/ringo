package ringo

import (
	"math"
	"sync/atomic"
)

var _ Buffer[any] = &ManyToOne[any]{}

// ManyToOne define a ring buffer safe for use by concurrent writers and a
// single reader.
type ManyToOne[T any] struct {
	buffer     []atomic.Pointer[box[T]]
	writeIndex atomic.Uint64
	// Also atomic as we read it on Push().
	readIndex atomic.Uint64
}

// NewManyToOne return a new ManyToOne ring buffer with the given
// capacity. The buffer is safe for one reader and multiple writer.
func NewManyToOne[T any](capacity int) *ManyToOne[T] {
	mto := &ManyToOne[T]{
		buffer: make([]atomic.Pointer[box[T]], capacity),
	}

	// First increment will overflow to 0.
	mto.writeIndex.Store(math.MaxUint64)

	return mto
}

// Size implements Buffer.
func (mto *ManyToOne[T]) Size() int {
	return len(mto.buffer)
}

// Push implements Buffer.
func (mto *ManyToOne[T]) Push(data T) (overwrite bool) {
	writeIndex := mto.writeIndex.Add(1)
	index := writeIndex % uint64(mto.Size())

	box := box[T]{
		index: writeIndex,
		data:  data,
	}

	tail1 := mto.readIndex.Load()
	old := mto.buffer[index].Swap(&box)

	if old == nil {
		return false
	}

	return old.index > tail1
}

// Shift implements Buffer.
func (mto *ManyToOne[T]) TryNext() (result T, ok bool, dropped int) {
	readIndex := mto.readIndex.Load()
	index := readIndex % uint64(mto.Size())
	box := mto.buffer[index].Load()

	// never written before.
	if box == nil {
		return
	}

	// already read
	if box.index < mto.readIndex.Load() {
		return
	}

	// cell have been overwritten
	if box.index > readIndex {
		dropped = int(box.index - readIndex)
		// set the tail so next shift box.data will be valid
		mto.readIndex.Store(box.index)
	}

	mto.readIndex.Add(1)
	return box.data, true, dropped
}
