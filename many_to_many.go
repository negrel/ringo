package ringo

import "sync"

type manyToMany struct {
	buffer *oneToOne

	mutex sync.Mutex
}

// ManyToMany return an efficient, thread-safe ring buffer
// with the given capacity. ManyToMany buffers are safe for
// concurrent writers and concurrent readers. This buffer is slower
// because it use the Mutex of the sync package, for faster ring buffer
// take a look at the other concurrent buffers.
func ManyToMany(capacity uint32) Buffer {
	return &manyToMany{
		buffer: &oneToOne{
			head:     ^uint64(0),
			buffer:   make([]Generic, capacity),
			capacity: uint64(capacity),
		},
		mutex: sync.Mutex{},
	}
}

func (mtm *manyToMany) Cap() uint32 {
	return uint32(mtm.buffer.capacity)
}

func (mtm *manyToMany) Push(data Generic) {
	mtm.mutex.Lock()
	mtm.buffer.Push(data)
	mtm.mutex.Unlock()
}

func (mtm *manyToMany) Shift() (shifted Generic, ok bool) {
	mtm.mutex.Lock()
	shifted, ok = mtm.buffer.Shift()
	mtm.mutex.Unlock()

	return shifted, ok
}
