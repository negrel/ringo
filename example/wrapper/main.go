package main

import (
	"fmt"
	"sync"

	"github.com/negrel/ringo"
)

func main() {
	// string wrapper
	buffer := &stRingBuffer{
		buffer: ringo.ManyToOne(100),
		// buffer: ringo.OneToOne(100),
		// buffer: ringo.ManyToMany(100),
	}

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		// Concurrent push to the buffer
		go func(j int) {
			buffer.Push("My string")
			wg.Done()
		}(i)
	}

	// Waiting all push are finish
	wg.Wait()

	for i := 0; i < 100; i++ {
		str := buffer.Shift()
		// Shifting & printing shifted data
		fmt.Printf("Shifting the %T \"%v\"\n", str, str)
	}
}
