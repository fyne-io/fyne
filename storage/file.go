// Package storage provides storage access and management functionality.
package storage

import (
	"fyne.io/fyne"
)

// OpenFileFromURI loads a file read stream from a resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
func OpenFileFromURI(uri fyne.URI) (fyne.URIReadCloser, error) {
	return fyne.CurrentApp().Driver().FileReaderForURI(uri)
}

// SaveFileToURI loads a file write stream to a resource identifier.
// This is mostly provided so that file references can be saved using their URI and written to again later.
func SaveFileToURI(uri fyne.URI) (fyne.URIWriteCloser, error) {
	return fyne.CurrentApp().Driver().FileWriterForURI(uri)
}
