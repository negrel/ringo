package ringo

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestPoller(t *testing.T) {
	t.Run("PollAvailableData", func(t *testing.T) {
		buf := NewRing[int](10)
		poller := NewPoller(buf, WithPollingInterval[int](time.Second))

		expected := rand.Int()

		poller.Push(expected)
		next, done, dropped := poller.Next()
		if next != expected {
			t.Fatalf("polled value doesn't match expected")
		}
		if dropped != 0 {
			t.Fatal("buffer reported some dropped value")
		}
		if done {
			t.Fatal("poller done without context cancellation")
		}
	})

	t.Run("PollUntilDataAvailable", func(t *testing.T) {
		buf := NewManyToOne[int](10)
		poller := NewPoller(buf, WithPollingInterval[int](time.Second))

		expected := rand.Int()

		go func() {
			time.Sleep(10 * time.Millisecond)
			poller.Push(expected)
		}()

		start := time.Now()
		next, done, dropped := poller.Next()
		end := time.Now()
		if next != expected {
			t.Fatalf("polled value doesn't match expected")
		}
		if dropped != 0 {
			t.Fatal("buffer reported some dropped value")
		}
		if end.Sub(start) < time.Second {
			t.Fatal("poller didn't used given polling interval")
		}
		if done {
			t.Fatal("poller done without context cancellation")
		}
	})

	t.Run("CancelPollerContext", func(t *testing.T) {
		buf := NewManyToOne[int](10)
		ctx, cancel := context.WithCancel(context.Background())
		poller := NewPoller(
			buf,
			WithPollingInterval[int](time.Second),
			WithPollingContext[int](ctx),
		)

		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		start := time.Now()
		next, done, dropped := poller.Next()
		end := time.Now()

		if end.Sub(start) < 500*time.Millisecond {
			t.Fatal("poller didn't used given context")
		}
		if next != 0 {
			t.Fatalf("polled value doesn't match expected")
		}
		if dropped != 0 {
			t.Fatal("buffer reported some dropped value")
		}
		if !done {
			t.Fatal("poller not done after context cancellation")
		}
	})
}
