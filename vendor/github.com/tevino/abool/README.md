# ABool :bulb:

[![Go Report Card](https://goreportcard.com/badge/github.com/tevino/abool)](https://goreportcard.com/report/github.com/tevino/abool)
[![GoDoc](https://godoc.org/github.com/tevino/abool?status.svg)](https://godoc.org/github.com/tevino/abool)

Atomic Boolean package for Go, optimized for performance yet simple to use.

Designed for cleaner code.

## Usage

```go
import "github.com/tevino/abool"

cond := abool.New()     // default to false

cond.Set()              // Sets to true
cond.IsSet()            // Returns true
cond.UnSet()            // Sets to false
cond.IsNotSet()         // Returns true
cond.SetTo(any)         // Sets to whatever you want
cond.SetToIf(new, old)  // Sets to `new` only if the Boolean matches the `old`, returns whether succeeded
cond.Toggle()           // Inverts the boolean then returns the value before inverting


// embedding
type Foo struct {
    cond *abool.AtomicBool  // always use pointer to avoid copy
}
```

## Benchmark

- Go 1.14.3
- Linux 4.19.0

```bash
goos: linux
goarch: amd64

# Read
BenchmarkMutexRead-4          	86662128	        14.2 ns/op
BenchmarkAtomicValueRead-4    	1000000000	         0.755 ns/op
BenchmarkAtomicBoolRead-4     	1000000000	         0.720 ns/op  # <--- This package


# Write
BenchmarkMutexWrite-4         	76237544	        13.6 ns/op
BenchmarkAtomicValueWrite-4   	79471124	        14.9 ns/op
BenchmarkAtomicBoolWrite-4    	178218270	         6.73 ns/op  # <--- This package

# CAS
BenchmarkMutexCAS-4           	29416574	        34.7 ns/op
BenchmarkAtomicBoolCAS-4      	171900002	         7.14 ns/op  # <--- This package

# Toggle
BenchmarkMutexToggle-4        	35212117	        34.5 ns/op
BenchmarkAtomicBoolToggle-4   	169871972	         7.02 ns/op  # <--- This package
```

## Special thanks to contributors

- [@barryz](https://github.com/barryz)
  - Added the `Toggle` method
