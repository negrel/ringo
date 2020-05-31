package ringo

type oneToOne struct {
	buffer []box
	head   uint64
	tail   uint64
	// Stored as uint64 to avoid conversion
	// but will never overflow uint32
	capacity uint64
}

// OneToOne return an efficient buffer with the given capacity.
// The buffer is safe for one reader and one writer.
func OneToOne(capacity uint32) Buffer {
	return &oneToOne{
		buffer:   make([]box, capacity),
		capacity: uint64(capacity),
	}
}

func (oto *oneToOne) Cap() uint32 {
	return uint32(oto.capacity)
}

func (oto *oneToOne) Push(data Generic) {
	index := oto.head % uint64(oto.capacity)

	box := box{
		index: oto.head,
		data:  data,
	}

	oto.head++
	oto.buffer[index] = box
}

func (oto *oneToOne) Shift() (Generic, bool) {
	index := oto.tail % oto.capacity
	oto.tail++

	box := oto.buffer[index]

	if box.data == nil {
		return nil, false
	}

	return box.data, box.index > oto.tail
}
