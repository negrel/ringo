package main

import "github.com/negrel/ringo"

// stRingBuffer will be our buffer wrapper
// it will only accept string data
type stRingBuffer struct {
	buffer ringo.Buffer
}

// Push will be used to write our string to the
// buffer.
func (sb *stRingBuffer) Push(str string) {
	// We convert the string to a Generic (= unsafe.Pointer) object
	sb.buffer.Push(ringo.Generic(&str))
}

// Shift the strings of the buffer.
func (sb *stRingBuffer) Shift() string {
	val, ok := sb.buffer.Shift()

	if !ok {
		return ""
	}

	// Reconvert it to string
	return *(*string)(val)
}
