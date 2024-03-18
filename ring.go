package ringo

// Buffer define common methods of ring buffers.
type Buffer[T any] interface {
	// Size returns size of internal buffer.
	Size() int
	// Push data to buffer.
	Push(data T) bool
	// Read value from buffer and returns it.
	// Returned boolean is true if a value was successfully read.
	// Int correspond to the number of dropped value since last read.
	TryNext() (T, bool, int)
}

var _ Buffer[any] = &Ring[any]{}

type Ring[T any] struct {
	buffer     []box[T]
	writeIndex uint64
	readIndex  uint64
}

func NewRing[T any](size int) *Ring[T] {
	return &Ring[T]{
		buffer:     make([]box[T], size),
		writeIndex: 0,
		readIndex:  1, // Makes first TryNext() return false if no write before.
	}
}

// Size implements Buffer.
func (r *Ring[T]) Size() int {
	return len(r.buffer)
}

// Push implements Buffer.
func (r *Ring[T]) Push(data T) (overwrite bool) {
	r.writeIndex++
	index := r.writeIndex % uint64(r.Size())

	old := r.buffer[index]
	r.buffer[index] = box[T]{r.writeIndex, data}

	return old.index > r.readIndex
}

// TryNext implements Buffer.
func (r *Ring[T]) TryNext() (result T, ok bool, dropped int) {
	index := r.readIndex % uint64(r.Size())
	box := r.buffer[index]

	// read index is ahead of write index.
	if box.index < r.readIndex {
		return
	}

	// writer is faster that reader and have overwritten data.
	if box.index > r.readIndex {
		dropped = int(box.index - r.readIndex)
		r.readIndex = box.index
	}

	r.readIndex++

	return box.data, true, dropped
}
