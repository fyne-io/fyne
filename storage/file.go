package storage

import (
	"fyne.io/fyne"
)

// OpenFileFromURI loads a file read stream from a resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
func OpenFileFromURI(uri fyne.URI) (fyne.FileReadCloser, error) {
	return fyne.CurrentApp().Driver().FileReaderForURI(uri)
}

// OpenFileFromURIString loads a file read stream from a simple resource identifier.
// This is a helper for apps that have persisted URI.String() and want to re-use.
func OpenFileFromURIString(uri string) (fyne.FileReadCloser, error) {
	return OpenFileFromURI(NewURI(uri))
}

// SaveFileToURI loads a file write stream to a resource identifier.
// This is mostly provided so that file references can be saved using their URI and written to again later.
func SaveFileToURI(uri fyne.URI) (fyne.FileWriteCloser, error) {
	return fyne.CurrentApp().Driver().FileWriterForURI(uri)
}

// SaveFileToURIString loads a file write stream to a simple resource identifier.
// This is a helper for apps that have persisted URI.String() and want to re-use.
func SaveFileToURIString(uri string) (fyne.FileWriteCloser, error) {
	return SaveFileToURI(NewURI(uri))
}
