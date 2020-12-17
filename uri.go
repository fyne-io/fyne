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
	// 'file://foo/bar.baz' is 'baz'. May return an empty string if no the
	// referenced resource has none.
	Extension() string

	// Name should return the base name of the item referenced by the URI.
	// For example, the Name() of 'file://foo/bar.baz' is 'bar.baz'.
	Name() string

	// MimeType should return the content type of the resource referenced
	// by the URI. The returned string should be in the format described
	// by section 5 of RFC2045 ("Content-Type Header Field").
	MimeType() string

	// Scheme should return the URI scheme of the URI. For example,
	// the Scheme() of 'file://foo/bar.baz` is 'file'.
	Scheme() string

	// Parent should return a URI referencing the parent resource of the
	// resource referenced by the URI. For example, the Parent() of
	// 'file://foo/bar.baz' is 'file://foo'.
	//
	// NOTE: it is not required that the implementation return a parent URI
	// with the same Scheme(), though this will normally be the case.
	//
	// This can fail in several ways:
	//
	// * If the URI refers to a filesystem root, then the Parent()
	//   implementation must return (nil, URIRootError).
	//
	// * If the URI refers to a resource which does not exist in a
	//   hierarchical context (e.g. the URI references something which
	//   does not have a semantically meaningful "parent"), the Parent()
	//   implementation may return an error.
	//
	// * If determining the parent of the referenced resource requires
	//   interfacing with some external system, failures may propagate
	//   through the Parent() implementation. For example if determining
	//   the parent of a file:// URI requires reading information from
	//   the filesystem, it could fail with a permission error.
	//
	// NOTE: To the extent possible, Parent() should not modify, create,
	// or interact with referenced resources. For example, it should
	// usually be possible to use Parent() to create a reference to
	// a resource which does not exist, though future operations on this
	// resource may fail.
	//
	// Since: 2.0
	Parent() (URI, error)

	// Child should return a URI referencing a resource nested
	// hierarchically below the given URI, identified by a string. For
	// example, the child with the string component 'quux' of
	// 'file://foo/bar' is 'file://foo/bar/quux'.
	//
	// This can fail in several ways:
	//
	// * If the URI refers to a resource which does not exist in a
	//   hierarchical context (e.g. the URI references something which
	//   does not have a semantically meaningful "child"), the Child()
	//   implementation may return an error.
	//
	// * If generating a reference to a child of the referenced resource
	//   requires interfacing with some external system, failures may
	//   propagate through the Child() implementation. It is expected that
	//   this case would occur very rarely if ever.
	//
	// NOTE: To the extent possible, Child() should not modify, create,
	// or interact with referenced resources. For example, it should
	// usually be possible to use Child() to create a reference to
	// a resource which does not exist, though future operations on this
	// resource may fail.
	//
	// Since: 2.0
	Child(URI, string) (URI, error)

	// Exists should determine if the resource referenced by the URI
	// exists.
	//
	// This can fail in several ways:
	//
	// * If checking the existence of a resource requires interfacing
	//   with some external system, then failures may propagate through
	//   Exists(). For example, checking the existence of a resource
	//   requires reading a directory may result in a permissions error.
	//
	// In the event that an error occurs, implementations of Exists() must
	// return false along with the error. It is understood that a non-nil
	// error value signals that the existence or non-existence of the
	// resource cannot be determined and is undefined.
	//
	// Since: 2.0
	Exists(URI) (bool, error)

	// Destroy should destroy, delete, or otherwise remove the resource
	// referenced by the URI.
	//
	// This can fail in several ways:
	//
	// * If removing the resource requires interfacing with some external
	//   system, failures may propagate through Destroy(). For example,
	//   deleting a file may fail with a permissions error.
	//
	// * If the referenced resource does not exist, attempting to destroy
	//   it should throw an error.
	//
	// Since: 2.0
	Destroy(URI) error
}

// ListableURI represents a URI that can have child items, most commonly a
// directory on disk in the native filesystem.
//
// Since: 1.4
type ListableURI interface {
	URI

	// List returns a list of child URIs of this URI.
	List() ([]URI, error)
}
