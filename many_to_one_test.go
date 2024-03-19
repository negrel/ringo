package ringo

import (
	"math/rand"
	"sync/atomic"
	"testing"
)

func TestManyToOne(t *testing.T) {
	t.Run("SingleGoRoutine", func(t *testing.T) {
		t.Run("SequentialReadWrite", func(t *testing.T) {
			size := 1000
			buffer := NewManyToOne[int](size)
			for i := 0; i < size; i++ {
				v := rand.Int()
				buffer.Push(v)

				r, ok, dropped := buffer.TryNext()
				if !ok {
					t.Fatal("TryNext() returned false, expecting true")
				}
				if r != v {
					t.Fatal("value read from buffer doesn't match expected")
				}
				if dropped != 0 {
					t.Fatal("buffer reported some dropped value")
				}
			}
		})

		t.Run("FullBufferThenEmptyIt", func(t *testing.T) {
			size := 1000
			buffer := NewManyToOne[int](size)
			pushedData := make([]int, size)

			for i := 0; i < size; i++ {
				pushedData[i] = rand.Int()
				buffer.Push(pushedData[i])
			}

			for i := 0; i < size; i++ {
				r, ok, dropped := buffer.TryNext()
				if !ok {
					t.Fatal("TryNext() returned false, expecting true")
				}
				if dropped != 0 {
					t.Fatal("buffer reported some dropped value:", dropped)
				}
				if r != pushedData[i] {
					t.Fatal("value read from buffer doesn't match expected")
				}
			}
		})
	})

	t.Run("MultipleWriter", func(t *testing.T) {
		writerCounter := 100

		size := 1000
		buffer := NewManyToOne[int](size)

		done := atomic.Int32{}

		for i := 0; i < writerCounter; i++ {
			go func(i int) {
				for j := 0; j < size; j++ {
					buffer.Push(j)
				}

				done.Add(1)
			}(i)
		}

		totalDropped := 0
		totalRead := 0

		// Read concurrently to writer.
		for done.Load() != int32(writerCounter) {
			r, ok, dropped := buffer.TryNext()
			totalDropped += dropped
			if ok {
				totalRead++
				if r < 0 || r > size {
					t.Fatal("value read from buffer doesn't match expected")
				}
			}
		}

		// Read unread values.
		for {
			r, ok, dropped := buffer.TryNext()
			totalDropped += dropped
			if ok {
				totalRead++
				if r < 0 || r > size {
					t.Fatal("value read from buffer doesn't match expected")
				}
			} else {
				break
			}
		}

		if totalDropped+totalRead != writerCounter*size {
			t.Fatal("number of read and dropped value doesn't match expected")
		}
	})

	t.Run("ReadEmptyBuffer", func(t *testing.T) {
		size := 1000
		buffer := NewManyToOne[int](size)

		next, ok, dropped := buffer.TryNext()
		if ok {
			t.Fatal("TryNext() returned true, expecting false")
		}
		if dropped != 0 {
			t.Fatal("buffer reported some dropped value:", dropped)
		}
		if next != 0 {
			t.Fatal("value read from buffer doesn't match expected")
		}
	})

	t.Run("DroppedData", func(t *testing.T) {
		size := 100
		buffer := NewManyToOne[int](size)

		for i := 0; i < 1000; i++ {
			buffer.Push(i)
		}

		next, ok, dropped := buffer.TryNext()
		if !ok {
			t.Fatal("TryNext() returned true, expecting false")
		}
		if dropped != 900 {
			t.Fatal("buffer reported some dropped value:", dropped)
		}
		if next != 900 {
			t.Fatal("value read from buffer doesn't match expected")
		}
	})
}
