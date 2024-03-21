package ringo

import (
	"fmt"
	"log/slog"
	"sync/atomic"
)

type CollisionHandler interface {
	OnCollision(buffer any)
}

type CollisionHandlerFunc func(buffer any)

// OnCollision implements CollisionHandler.
func (chf CollisionHandlerFunc) OnCollision(buffer any) {
	chf(buffer)
}

var globalCollisionHandler = atomic.Pointer[CollisionHandler]{}

func init() {
	defaultHandler := CollisionHandlerFunc(defaultCollisionHandler)
	SetCollisionHandler(defaultHandler)
}

func defaultCollisionHandler(buffer any) {
	slog.Warn(
		"ringo: ring buffer collision detected, consider increasing size of your ring buffer",
		slog.String("ring_buffer", fmt.Sprintf("%T %p", buffer, buffer)))
}

// SetCollisionHandler sets collision handler called when a collision occurs.
func SetCollisionHandler(handler CollisionHandler) {
	globalCollisionHandler.Store(&handler)
}
