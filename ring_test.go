package ringo

import (
	"math/rand"
	"testing"
)

func TestRing(t *testing.T) {
	t.Run("SequentialReadWrite", func(t *testing.T) {
		size := 1000
		buffer := NewRing[int](size)
		for i := 0; i < size; i++ {
			v := rand.Int()
			overwrite := buffer.Push(v)
			if overwrite {
				t.Fatal("unexpected overwrite")
			}

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
		buffer := NewRing[int](size)
		pushedData := make([]int, size)

		for i := 0; i < size; i++ {
			pushedData[i] = rand.Int()
			overwrite := buffer.Push(pushedData[i])
			if overwrite {
				t.Fatal("unexpected overwrite")
			}
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

	t.Run("ReadEmptyBuffer", func(t *testing.T) {
		size := 1000
		buffer := NewRing[int](size)

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
}
