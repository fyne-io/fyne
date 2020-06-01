package fyne

import (
	"fmt"
	"io"
)

// FileReadCloser represents a cross platform data stream from a file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type FileReadCloser interface {
	io.ReadCloser
	Name() string
	URI() URI
}

// FileWriteCloser represents a cross platform data writer for a file resource.
// This will normally refer to a local file resource.
type FileWriteCloser interface {
	io.WriteCloser
	Name() string
	URI() URI
}

// URI represents the identifier of a resource on a target system.
// This resource may be a file or another data source such as an app or file sharing system.
type URI interface {
	fmt.Stringer
	Extension() string
	MimeType() string
	Scheme() string
}
