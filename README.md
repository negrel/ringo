<h1 align="center"><img height="250" src="https://raw.githubusercontent.com/negrel/ringo/master/.github/atom.svg"></h1>

<p align="center">
	<a href="https://pkg.go.dev/github.com/negrel/ringo">
		<img src="https://godoc.org/github.com/negrel/ringo?status.svg">
	</a>
	<a href="https://goreportcard.com/badge/github.com/negrel/ringo">
		<img src="https://goreportcard.com/badge/github.com/negrel/ringo">
	</a>
	<a href="https://github.com/negrel/ringo/raw/master/LICENSE">
		<img src="https://img.shields.io/badge/license-MIT-green">
	</a>
</p>

# :atom_symbol: Ringo - Fast, lock free ring buffers.

A thread safe, lock free, efficient ring buffer library.

Ringo is heavily inspired by [go-diodes](https://github.com/cloudfoundry/go-diodes/) 
but aims to provide a more safe (no unsage and type safe) alternative.

## Features

- :zap: [**Efficient**](https://github.com/negrel/ringo#zap-benchmarks)
- **Thread-safe** : manipulated via [atomics](https://pkg.go.dev/sync/atomic) operations.
- **Type-safe** : buffers are implemented using [Go 1.18 generics](https://go.dev/doc/tutorial/generics).

## Getting started

### Installation

Using **go get** :

```bash
go get github.com/negrel/ringo.git
```

Using **go modules** :

```go
package "your_package_name"

import (
	"github.com/negrel/ringo"
)

func main() {
    // Your code here
}
```

then

```bash
go mod tidy
```

### Example: Basic Use

```go
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
```

## :zap: Benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/negrel/ringo
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics
BenchmarkRing
BenchmarkRing-16                200605705                6.669 ns/op          16 B/op          0 allocs/op
BenchmarkManyToOne
BenchmarkManyToOne-16           36673826                28.27 ns/op           16 B/op          1 allocs/op
BenchmarkManyToOneWaiter
BenchmarkManyToOneWaiter-16     32007580                35.13 ns/op           16 B/op          1 allocs/op
BenchmarkManyToOnePoller
BenchmarkManyToOnePoller-16     37656290                34.64 ns/op           16 B/op          1 allocs/op
PASS
ok      github.com/negrel/ringo 6.841s
```

## Known issues

If a diode was to be written to 18446744073709551615+1 times it would overflow
a uint64. This will cause problems if the size of the diode is not a power of
two (2^x). If you write into a diode at the rate of one message every
nanosecond, without restarting your process, it would take you 584.54 years to
encounter this issue.

## :stars: Show your support

Please give a :star: if this project helped you!

## Acknowledgments

<a href="https://iconscout.com/icons/atomic" target="_blank">Atomic Icon</a> by <a href="https://iconscout.com/contributors/oviyan">Vignesh Oviyan</a> on <a href="https://iconscout.com">Iconscout</a>

## :scroll: License

MIT © Alexandre Negrel
