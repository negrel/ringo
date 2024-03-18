package ringo

import (
	"sync"
	"testing"
)

func BenchmarkRing(b *testing.B) {
	buffer := NewRing[int](b.N)

	for i := 0; i < b.N; i++ {
		buffer.Push(i)
	}
	for i := 0; i < b.N; i++ {
		value, ok, dropped := buffer.TryNext()
		if !ok || dropped != 0 || value != i {
			b.FailNow()
		}
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func BenchmarkManyToOne(b *testing.B) {
	// Number of concurrent writer
	wCount := 100

	var wg sync.WaitGroup
	wg.Add(wCount)

	loop := max(b.N, wCount)
	loopPerGoroutine := loop / wCount

	buffer := NewManyToOne[int](loop)
	writer := func() {
		defer wg.Done()
		for i := 0; i < loopPerGoroutine; i++ {
			buffer.Push(i)
		}
	}

	b.ResetTimer()
	for i := 0; i < wCount; i++ {
		go writer()
	}

	wg.Wait()

	for i := 0; i < b.N; i++ {
		value, ok, dropped := buffer.TryNext()

		if i >= wCount*loopPerGoroutine {
			if ok {
				b.FailNow()
			}
		} else {
			if !ok || dropped != 0 || value < 0 || value > b.N {
				b.Fatalf("i: %v, ok: %v, dropped: %v, value: %v", i, ok, dropped, value)
			}
		}
	}
}

func BenchmarkManyToOneWaiter(b *testing.B) {
	buffer := NewWaiter(NewManyToOne[int](b.N))

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			buffer.Push(i)
		}
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buffer.Next()
	}
}

func BenchmarkManyToOnePoller(b *testing.B) {
	buffer := NewPoller(NewManyToOne[int](b.N))

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			buffer.Push(i)
		}
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buffer.Next()
	}
}
