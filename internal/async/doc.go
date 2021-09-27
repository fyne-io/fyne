//go:generate go run gen.go

// Package async provides unbounded channel data structures that are
// designed for caching unlimited number of a concrete type. For better
// performance, a given type should be less or euqal than 16 bytes.
//
// Delicate dance: One must aware that an unbounded channel may lead to
// OOM when the consuming speed of the buffer is lower than the producing
// speed constantly. However, such a channel may be fairly used for event
// delivering if the consumer of the channel consumes the incoming
// forever, such as even processing.
//
// One must close such a channel via Close() method, closing the input
// channel via close() built-in method can leads to memory leak.
//
// To support a new type, one may add the required data in the gen.go,
// for instances:
//
// 	types := map[string]data{
// 		"fyne_canvasobject.go": data{
// 			Type: "fyne.CanvasObject",
// 			Name: "CanvasObject",
// 			Imports: `import "fyne.io/fyne/v2"`,
// 		},
// 		"func.go": data{
// 			Type:    "func()",
// 			Name:    "Func",
// 			Imports: "",
// 		},
// 	}
//
// then run: `go generate ./...` to generate more desired unbounded channels.
package async
