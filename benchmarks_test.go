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

// func BenchmarkOneToOne(b *testing.B) {
// 	buffer := NewOneToOne(uint32(b.N))
//
// 	var wg sync.WaitGroup
// 	wg.Add(1)
//
// 	b.ResetTimer()
//
// 	for i := 0; i < b.N; i++ {
// 		y := i
// 		buffer.Push(Generic(&y))
// 	}
//
// 	go func() {
// 		defer wg.Done()
// 		for i := 0; i < b.N; i++ {
// 			data, _ := buffer.Shift()
// 			result = data
// 		}
// 	}()
//
// 	wg.Wait()
// }

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

// func BenchmarkManyToMany(b *testing.B) {
// 	buffer := ManyToMany(uint32(b.N))
// 	// Number of concurrent writer
// 	rwCount := 100
//
// 	var wg sync.WaitGroup
// 	wg.Add(rwCount)
//
// 	// Avoid b.N / rwCount = 0
// 	loop := float64(b.N)
// 	if loop == 1 {
// 		loop = float64(rwCount)
// 	}
// 	loop /= float64(rwCount)
// 	loopPerGoroutine := int(math.Ceil(loop))
//
// 	writer := func() {
// 		defer wg.Done()
// 		for i := 0; i < loopPerGoroutine; i++ {
// 			j := i
// 			buffer.Push(Generic(&j))
// 		}
// 	}
//
// 	reader := func() {
// 		defer wg.Done()
// 		for i := 0; i < loopPerGoroutine; i++ {
// 			value, _ := buffer.Shift()
// 			_ = value
// 		}
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < rwCount; i++ {
// 		go writer()
// 	}
//
// 	wg.Wait()
// 	wg.Add(rwCount)
//
// 	for i := 0; i < rwCount; i++ {
// 		go reader()
// 	}
// 	wg.Wait()
// }
