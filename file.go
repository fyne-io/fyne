package fyne

import "io"

// File represents a cross platform file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type File interface {
	Open() (io.ReadCloser, error)
	Save() (io.WriteCloser, error)
	ReadOnly() bool
	Name() string
	URI() string
}

// FileFromURI loads a file descriptor from a saved resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
func FileFromURI(uri string) File {
	return CurrentApp().Driver().FileFromURI(uri)
}
