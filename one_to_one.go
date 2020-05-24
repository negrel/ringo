package ringo

type oneToOne buffer

// OneToOne return an efficient buffer with the given capacity.
// The buffer is safe for one reader and one writer.
func OneToOne(capacity uint32) Buffer {
	return &oneToOne{
		buffer:   make([]Generic, capacity),
		capacity: uint64(capacity),
	}
}

func (oto *oneToOne) Cap() uint32 {
	return uint32(oto.capacity)
}

// Push the given data to the buffer.
func (oto *oneToOne) Push(data Generic) {
	index := oto.head % uint64(oto.capacity)

	box := box{
		index: oto.head,
		data:  data,
	}

	oto.head++

	oto.buffer[index] = Generic(&box)
}

// Push the given data to the buffer and return if
// the data is valid.
func (oto *oneToOne) Shift() (Generic, bool) {
	i := oto.tail % oto.capacity
	oto.tail++

	box := (*box)(oto.buffer[i])

	if box == nil {
		return nil, false
	}

	return box.data, box.index > oto.tail
}
