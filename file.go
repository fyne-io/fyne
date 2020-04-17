package fyne

import "io"

// FileReader represents a cross platform data stream from a file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type FileReader interface {
	io.ReadCloser
	Name() string
	URI() string
}

// FileWriter represents a cross platform data writer for a file resource.
// This will normally refer to a local file resource.
type FileWriter interface {
	io.WriteCloser
	Name() string
	URI() string
}

// OpenFileFromURI loads a file read stream from a resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
func OpenFileFromURI(uri string) (FileReader, error) {
	return CurrentApp().Driver().FileReaderForURI(uri)
}

// SaveFileToURI loads a file write stream to a resource identifier.
// This is mostly provided so that file references can be saved using their URI and written to again later.
func SaveFileToURI(uri string) (FileWriter, error) {
	return CurrentApp().Driver().FileWriterForURI(uri)
}
