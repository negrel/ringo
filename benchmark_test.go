package ringo

import (
	"math"
	"sync"
	"testing"
)

var result Generic

func BenchmarkOneToOne(b *testing.B) {
	buffer := OneToOne(uint32(b.N))

	var wg sync.WaitGroup
	wg.Add(1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		y := i
		buffer.Push(Generic(&y))
	}

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			data, _ := buffer.Shift()
			result = data
		}
	}()

	wg.Wait()
}

func BenchmarkManyToOne(b *testing.B) {
	buffer := ManyToOne(uint32(b.N))
	// Number of concurrent writer
	wCount := 100

	var wg sync.WaitGroup
	wg.Add(wCount)

	// Avoid b.N / rwCount = 0
	loop := float64(b.N)
	if loop == 1 {
		loop = float64(wCount)
	}
	loop /= float64(wCount)
	loopPerGoroutine := int(math.Ceil(loop))

	writer := func() {
		defer wg.Done()
		for i := 0; i < loopPerGoroutine; i++ {
			j := i
			buffer.Push(Generic(&j))
		}
	}

	b.ResetTimer()
	for i := 0; i < wCount; i++ {
		go writer()
	}

	wg.Wait()

	for i := 0; i < b.N; i++ {
		data, _ := buffer.Shift()
		result = data
	}
}

func BenchmarkManyToMany(b *testing.B) {
	buffer := ManyToMany(uint32(b.N))
	// Number of concurrent writer
	rwCount := 100

	var wg sync.WaitGroup
	wg.Add(rwCount)

	// Avoid b.N / rwCount = 0
	loop := float64(b.N)
	if loop == 1 {
		loop = float64(rwCount)
	}
	loop /= float64(rwCount)
	loopPerGoroutine := int(math.Ceil(loop))

	writer := func() {
		defer wg.Done()
		for i := 0; i < loopPerGoroutine; i++ {
			j := i
			buffer.Push(Generic(&j))
		}
	}

	reader := func() {
		defer wg.Done()
		for i := 0; i < loopPerGoroutine; i++ {
			value, _ := buffer.Shift()
			_ = value
		}
	}

	b.ResetTimer()
	for i := 0; i < rwCount; i++ {
		go writer()
	}

	wg.Wait()
	wg.Add(rwCount)

	for i := 0; i < rwCount; i++ {
		go reader()
	}
	wg.Wait()
}
