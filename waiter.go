package ringo

import (
	"context"
)

// Waiter will use a channel signal to alert the reader to when data is
// available.
type Waiter[T any] struct {
	Buffer[T]
	c   chan struct{}
	ctx context.Context
}

// WaiterConfigOption can be used to setup the waiter.
type WaiterConfigOption[T any] func(*Waiter[T])

// WithWaiterContext sets the context to cancel any retrieval (Next()). It
// will not change any results for adding data (Set()). Default is
// context.Background().
func WithWaiterContext[T any](ctx context.Context) WaiterConfigOption[T] {
	return WaiterConfigOption[T](func(c *Waiter[T]) {
		c.ctx = ctx
	})
}

// NewWaiter returns a new Waiter that wraps the given ring buffer.
func NewWaiter[T any](buffer Buffer[T], opts ...WaiterConfigOption[T]) Waiter[T] {
	w := Waiter[T]{
		Buffer: buffer,
		c:      make(chan struct{}, 1),
		ctx:    context.Background(),
	}
	w.Buffer = buffer
	w.c = make(chan struct{}, 1)

	for _, opt := range opts {
		opt(&w)
	}

	return w
}

// Push invokes the wrapped Buffer's Push with the given data and uses broadcast
// to wake up any readers.
func (w *Waiter[T]) Push(data T) {
	w.Buffer.Push(data)
	w.broadcast()
}

// broadcast sends to the channel if it can.
func (w *Waiter[T]) broadcast() {
	select {
	case w.c <- struct{}{}:
	default:
	}
}

// Next returns the next data point on the wrapped ring buffer. If there is no new
// data, it will wait for Set to be called or the context to be done. If the
// context is done, then default value of T will be returned.
func (w *Waiter[T]) Next() (next T, done bool, dropped int) {
	var ok bool
	for {
		next, ok, dropped = w.Buffer.TryNext()
		if ok {
			return
		}
		select {
		case <-w.ctx.Done():
			done = true
			return
		case <-w.c:
		}
	}
}
