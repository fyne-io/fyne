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

	// Extension should return the file extension of the resource
	// referenced by the URI. For example, the Extension() of
	// 'file://foo/bar.baz' is 'baz'. May return an empty string if the
	// referenced resource has none.
	Extension() string

	// Name should return the base name of the item referenced by the URI.
	// For example, the Name() of 'file://foo/bar.baz' is 'bar.baz'.
	Name() string

	// MimeType should return the content type of the resource referenced
	// by the URI. The returned string should be in the format described
	// by Section 5 of RFC2045 ("Content-Type Header Field").
	MimeType() string

	// Scheme should return the URI scheme of the URI. For example,
	// the Scheme() of 'file://foo/bar.baz` is 'file'.
	Scheme() string
}

// ListableURI represents a URI that can have child items, most commonly a
// directory on disk in the native filesystem.
//
// Deprecated: use the IsListable() and List() methods that operate on URI in
// the storage package instead.
//
// Since: 1.4
type ListableURI interface {
	URI

	// List returns a list of child URIs of this URI.
	List() ([]URI, error)
}
