//go:generate go run gen.go

// Package async provides unbounded channel data structures that are
// designed for caching unlimited number of a concrete type. For better
// performance, a given type should have a smaller or euqal size of 16-bit.
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
