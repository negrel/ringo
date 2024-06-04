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
	readIndex        atomic.Uint64
	collisionHandler CollisionHandler
}

type ManyToOneOption[T any] func(*ManyToOne[T])

// WithManyToOneCollisionHandler sets ManyToOne ring buffer collision handler.
// If this option is not provided ring buffer defaults to global handler.
func WithManyToOneCollisionHandler[T any](ch CollisionHandler) ManyToOneOption[T] {
	return func(mto *ManyToOne[T]) {
		mto.collisionHandler = ch
	}
}

// NewManyToOne return a new ManyToOne ring buffer with the given
// size. The buffer is safe for one reader and multiple writer.
func NewManyToOne[T any](size int, options ...ManyToOneOption[T]) *ManyToOne[T] {
	if size <= 0 {
		panic("ring buffer size can't be negative or zero")
	}

	mto := &ManyToOne[T]{
		buffer:           make([]atomic.Pointer[box[T]], size),
		collisionHandler: *globalCollisionHandler.Load(),
	}

	// First increment will overflow to 0.
	mto.writeIndex.Store(math.MaxUint64)

	for _, opt := range options {
		opt(mto)
	}

	return mto
}

// Size implements Buffer.
func (mto *ManyToOne[T]) Size() int {
	return len(mto.buffer)
}

// Push implements Buffer.
func (mto *ManyToOne[T]) Push(data T) {
	for {
		writeIndex := mto.writeIndex.Add(1)
		index := writeIndex % uint64(mto.Size())

		old := mto.buffer[index].Load()
		if old != nil && old.index > writeIndex {
			mto.collisionHandler.OnCollision(mto)
			continue
		}

		box := box[T]{
			index: writeIndex,
			data:  data,
		}

		if !mto.buffer[index].CompareAndSwap(old, &box) {
			mto.collisionHandler.OnCollision(mto)
			continue
		}

		return
	}
}

// TryNext implements Buffer.
func (mto *ManyToOne[T]) TryNext() (result T, ok bool, dropped int) {
	readIndex := mto.readIndex.Load()
	index := readIndex % uint64(mto.Size())
	// Swap(nil) is tempting to allow garbage collection of box[T] but
	// it breaks collision detection (CompareAndSwap(old, ...) call) of Push().
	box := mto.buffer[index].Load()

	if box == nil {
		return
	}

	// already read
	if box.index < readIndex {
		return
	}

	// cell have been overwritten
	if box.index > readIndex {
		dropped = int(box.index - readIndex)
		// move readIndex to catch up.
		mto.readIndex.Store(box.index)
	}

	mto.readIndex.Add(1)
	data := box.data

	// Replace box.data with zeroed value to allow gc to collect box.data or
	// its content.
	var zeroT T
	box.data = zeroT

	return data, true, dropped
}
