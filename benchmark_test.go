package ringo

import (
	"log"
	"sync"
	"testing"

	"code.cloudfoundry.org/go-diodes"
)

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
		_ = data
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
			_ = data
		}
	}()

	wg.Wait()
}

func BenchmarkConcurrentManyToOne(b *testing.B) {
	buffer := ManyToOne(uint32(b.N))

	var wg sync.WaitGroup
	wg.Add(b.N)

	b.ResetTimer()

	// 2 concurrent writers.
	writer := func() {
		log.Println("Goroutines for ", b.N)
		for i := 0; i < b.N/2; i++ {
			j := i
			buffer.Push(Generic(&j))
			wg.Done()
		}
	}

	go writer()
	go writer()

	wg.Wait()

	for i := 0; i < b.N; i++ {
		data, _ := buffer.Shift()
		_ = data
	}

}
