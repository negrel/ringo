package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/negrel/ringo"
)

func main() {
	buffer := ringo.NewManyToOne[int](1024)

	go func() {
		for i := 0; i < math.MaxInt; i++ {
			buffer.Push(i)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Minute)
		cancel()
	}()

	poller := ringo.NewPoller(buffer, ringo.WithPollingContext[int](ctx))
	for {
		next, done, dropped := poller.Next()
		// Writer is faster than reader, some data was overwritten.
		if dropped > 0 {
			log.Printf("lost %v int", dropped)
		}
		// Context canceled.
		if done {
			break
		}
		log.Print(next)
	}
}
