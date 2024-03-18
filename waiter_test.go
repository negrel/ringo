package ringo

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestWaiter(t *testing.T) {
	t.Run("WaitsForPush", func(t *testing.T) {
		buf := NewManyToOne[int](10)
		waiter := NewWaiter(buf)

		expected := rand.Int()

		go func() {
			time.Sleep(500 * time.Millisecond)
			_ = waiter.Push(expected)
		}()

		start := time.Now()
		next, dropped := waiter.Next()
		end := time.Now()

		if next != expected {
			t.Fatalf("waited value doesn't match expected")
		}
		if dropped != 0 {
			t.Fatal("buffer reported some dropped value")
		}
		if end.Sub(start) < 500*time.Millisecond || end.Sub(start) > 550*time.Millisecond {
			t.Fatal("waiter is slow")
		}
	})

	t.Run("CancelWaiterContext", func(t *testing.T) {
		buf := NewManyToOne[int](10)
		ctx, cancel := context.WithCancel(context.Background())

		waiter := NewWaiter(buf, WithWaiterContext[int](ctx))

		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		start := time.Now()
		next, dropped := waiter.Next()
		end := time.Now()

		if end.Sub(start) < 500*time.Millisecond {
			t.Fatal("waiter didn't used given context")
		}
		if next != 0 {
			t.Fatalf("waited value doesn't match expected")
		}
		if dropped != 0 {
			t.Fatal("buffer reported some dropped value")
		}
	})
}
