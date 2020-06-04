package fyne

import (
	"fmt"
	"io"
)

// URIReadCloser represents a cross platform data stream from a file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type URIReadCloser interface {
	io.ReadCloser
	Name() string
	URI() URI
}

// URIWriteCloser represents a cross platform data writer for a file resource.
// This will normally refer to a local file resource.
type URIWriteCloser interface {
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
