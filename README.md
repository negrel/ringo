<h1 align="center"><img height="250" src="https://raw.githubusercontent.com/negrel/ringo/master/.github/atom.svg"></h1>

<p align="center">
	<a href="https://goreportcard.com/badge/github.com/negrel/ringo">
		<img src="https://goreportcard.com/badge/github.com/negrel/ringo">
	</a>
	<a href="https://github.com/negrel/ringo/raw/master/LICENSE">
		<img src="https://img.shields.io/badge/license-MIT-green">
	</a>
</p>

# :atom_symbol: Ringo - Efficient ring buffers
*Atomic buffers are thread safe, lock free, efficient ring buffers.*  
Ringo is inspired by [go-diodes](https://github.com/cloudfoundry/go-diodes/) but is faster.

## Features

- **Thread-safe** : manipulated via [atomics](https://pkg.go.dev/sync/atomic) operations.
- :zap: [**Efficient**](https://github.com/negrel/ringo#zap-benchmarks)
- **Easy to use** : Check the [examples](https://github.com/negrel/ringo/tree/master/example)
- **Untyped** : step around type safety thanks to the standard [unsafe](https://pkg.go.dev/unsafe) package.

## Installation

Using **go get** :

```bash
go get -u github.com/negrel/ringo.git
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
go mod init
```

## Getting started

Take a look at the [examples](https://github.com/negrel/ringo/tree/master/example) and especially the *"wrapper"* one, it is recommended to use a wrapper for type safety.

## :zap: Benchmarks
i5-8250U @ 8x 3.4GHz:

```
goos: linux
goarch: amd64
pkg: github.com/negrel/ringo
BenchmarkOneToOne-8             21199524                57.6 ns/op             8 B/op          1 allocs/op
BenchmarkManyToOne-8            11959047                84.3 ns/op            24 B/op          2 allocs/op
BenchmarkManyToMany-8            2868415               414 ns/op               8 B/op          1 allocs/op
PASS
ok      github.com/negrel/ringo 6.002s
```

## :stars: Show your support

Please give a :star: if this project helped you!

## Acknowledgments

<a href="https://iconscout.com/icons/atomic" target="_blank">Atomic Icon</a> by <a href="https://iconscout.com/contributors/oviyan">Vignesh Oviyan</a> on <a href="https://iconscout.com">Iconscout</a>

## :scroll: License

MIT © Alexandre Negrel
