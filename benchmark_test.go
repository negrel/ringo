package ringo

import (
	"math"
	"sync"
	"testing"

	"code.cloudfoundry.org/go-diodes"
)

var result Generic

func BenchmarkDiodeOneToOne(b *testing.B) {
	buffer := diodes.NewOneToOne(b.N, nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		y := i
		buffer.Set(diodes.GenericDataType(&y))
	}

	for i := 0; i < b.N; i++ {
		data, _ := buffer.TryNext()
		_ = data
	}
}

func BenchmarkOneToOne(b *testing.B) {
	buffer := OneToOne(uint32(b.N))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		y := i
		buffer.Push(Generic(&y))
	}

	for i := 0; i < b.N; i++ {
		data, _ := buffer.Shift()
		result = data
	}

}

func BenchmarkConcurrentDiodeOneToOne(b *testing.B) {
	buffer := diodes.NewOneToOne(b.N, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		y := i
		buffer.Set(diodes.GenericDataType(&y))
	}

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			data, _ := buffer.TryNext()
			_ = data
		}
	}()

	wg.Wait()
}

func BenchmarkConcurrentOneToOne(b *testing.B) {
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

func BenchmarkConcurrentManyToOne(b *testing.B) {
	buffer := ManyToOne(uint32(b.N))
	grc := 4

	var wg sync.WaitGroup
	wg.Add(grc)

	// 2 concurrent writers.
	// Avoid b.N / grc = 0
	loop := float64(b.N)
	if loop == 1 {
		loop = float64(grc)
	}
	loop /= float64(grc)
	loopPerGoroutine := int(math.Ceil(loop))

	writer := func() {
		defer wg.Done()
		for i := 0; i < loopPerGoroutine; i++ {
			j := i
			buffer.Push(Generic(&j))
		}
	}

	b.ResetTimer()
	for i := 0; i < grc; i++ {
		go writer()
	}

	wg.Wait()

	for i := 0; i < b.N; i++ {
		data, _ := buffer.Shift()
		result = data
	}
}
