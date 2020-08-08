package fyne

import (
	"fmt"
	"io"
)

// URIReadCloser represents a cross platform data stream from a file or provider of data.
// It may refer to an item on a filesystem or data in another application that we have access to.
type URIReadCloser interface {
	io.ReadCloser
	// Deprecated, use URI().Name() instead
	Name() string
	URI() URI
}

// URIWriteCloser represents a cross platform data writer for a file resource.
// This will normally refer to a local file resource.
type URIWriteCloser interface {
	io.WriteCloser
	// Deprecated, use URI().Name() instead
	Name() string
	URI() URI
}

// URI represents the identifier of a resource on a target system.
// This resource may be a file or another data source such as an app or file sharing system.
type URI interface {
	fmt.Stringer
	Extension() string
	Name() string
	MimeType() string
	Scheme() string
	Parent() URI
}

// ListableURI represents a URI that can have child items, most commonly a
// directory on disk in the native filesystem.
type ListableURI interface {
	URI

	// List returns a list of child URIs of this URI.
	List() []URI
}

// URIHandler implements the needed methods for a URI to interact with the
// outside world.
//
// The theory is that a URI refers to some resource that presumably exists
// somewhere in the ether. A URIHandler makes this relationship by concrete
// by defining operations that can be done on the intersection of a particular
// URI in conjunction with a specific I/O system.
//
// In general, all of the functions in this interface will error if the URI
// scheme is not one it knows how to handle.
type URIHandler interface {

	// Validate checks that the given URI is a valid instance of a URI
	// which this handler understands. For example, and FTP handler
	// might return an error if the URI provided uses the HTTP Scheme.
	//
	// This function should never cause any network or disk activity, or
	// otherwise interact with the resource in any way.
	Validate(URI) error

	// Create attempts to ensure that the resource to which the URI refers
	// to exists.
	//
	// If the resource already exists, this is a No-Op (think the `touch`
	// command).  If it doesn't exist, it will be created.
	Create(URI) error

	// WriterTo attempts to open the URI for writing, and will return an
	// error if this is not possible.
	WriterTo(URI) (URIWriteCloser, error)

	// ReaderFrom attempts to open a URI for reading, and will return an
	// error if this is not possible.
	ReaderFrom(URI) (URIReadCloser, error)

	// ListerOf attempts to upgrade a URI to a ListableURI, and will return
	// an error if this is not possible.
	ListerOf(URI) (ListableURI, error)
}
