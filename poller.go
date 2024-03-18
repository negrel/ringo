package ringo

import (
	"context"
	"time"
)

// Poller polls a ring buffer until a value is available.
type Poller[T any] struct {
	Buffer[T]
	interval time.Duration
	ctx      context.Context
}

// PollerConfigOption can be used to setup the poller.
type PollerConfigOption[T any] func(*Poller[T])

// WithPollingInterval sets the interval at which the ring buffer is queried
// for new data. The default is 10ms.
func WithPollingInterval[T any](interval time.Duration) PollerConfigOption[T] {
	return PollerConfigOption[T](func(c *Poller[T]) {
		c.interval = interval
	})
}

// WithPollingContext sets the context to cancel any retrieval (Next()). It
// will not change any results for adding data (Set()). Default is
// context.Background().
func WithPollingContext[T any](ctx context.Context) PollerConfigOption[T] {
	return PollerConfigOption[T](func(c *Poller[T]) {
		c.ctx = ctx
	})
}

// NewPoller returns a new Poller that wraps the given buffer.
func NewPoller[T any](buf Buffer[T], opts ...PollerConfigOption[T]) Poller[T] {
	p := Poller[T]{
		Buffer:   buf,
		interval: 10 * time.Millisecond,
		ctx:      context.Background(),
	}

	for _, o := range opts {
		o(&p)
	}

	return p
}

// Next polls the buffer until data is available or until the context is done.
// If the context is done, then default value of T will be returned.
func (p *Poller[T]) Next() (next T, dropped int) {
	var ok bool
	for {
		next, ok, dropped = p.Buffer.TryNext()
		if !ok {
			if p.isDone() {
				return
			}

			time.Sleep(p.interval)
			continue
		}

		return
	}
}

func (p *Poller[T]) isDone() bool {
	select {
	case <-p.ctx.Done():
		return true
	default:
		return false
	}
}
