package fyne

import "io"

// FileReadCloser represents a cross platform data stream from a file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type FileReadCloser interface {
	io.ReadCloser
	Name() string
	URI() string
}

// FileWriteCloser represents a cross platform data writer for a file resource.
// This will normally refer to a local file resource.
type FileWriteCloser interface {
	io.WriteCloser
	Name() string
	URI() string
}
