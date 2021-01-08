package storage

// Repository represents a storage repository, which is a set of methods which
// implement specific functions on a URI. Repositories are registered to handle
// specific URI schemes, and the higher-level functions that operate on URIs
// internally look up an appropriate method from the relevant Repository.
//
// Some repositories may choose not to implement some functions; in such cases,
// they should simply return nil for un-implemented functions.
//
// Repository methods which return (func ..., bool) can use the bool flag to
// indicate if they would like to use a default/generic implementation of the
// relevant method (by returning a true value), or if the operations is not
// supported by the type of information the URI works with (a false value).
//
// Repository methods which only return a func will always use a
// default/generic implementation if the repository does not provide one (e.g.
// returns nil). These operations must be supported by all URIs.
//
// To the extent possible, Repository interfaces should attempt to use these
// methods as a surface area for optimization. For example, a repository for
// accessing ssh:// URIs might perform copy operations on the remote end,
// rather than reading all the content to the local side and then writing
// it back over the wire.
//
// Since: 2.0.0
type Repository interface {

	// ParseURIImpl returns a method that given a text string creates a
	// new URI instance. This may cause an error, for example if the string
	// is an invalid URI or contains information that is invalid in the
	// context of the repository's scheme.
	//
	// This method will be called only after the URI scheme has been
	// verified to match the one the repository was registered to handle,
	// however the URI is not validated further than this. It is required
	// by convention (but not enforced by technical means) that the
	// returned function should throw an error if the string is not a valid
	// IETF RFC 3986 complaint URI. Fail to validate this at your own risk.
	//
	// NOTE: it is highly recommended to use net/url (
	// https://golang.org/pkg/net/url/ ) for parsing URI strings unless you
	// have a good reason not to. It is mature and implements RFC 3986.
	//
	// If no implementation is provided, a generic URI based on RFC3986 is
	// parsed without any special validation logic.
	//
	// Since 2.0.0
	ParseURIImpl() func(string) (URI, error)

	// ExtensionImpl returns a method that given a URI matching the
	// scheme this repository is registered to handle, will return the
	// file extension for the URI.
	//
	// If no implementation is provided, then a generic one will be used,
	// which will string split on instances of the '.' character, and the
	// final such component will be returned.
	//
	// Since 2.0.0
	ExtensionImpl() func(URI) string

	// NameImpl returns a method that given a URI matching the scheme
	// this repository is registered to handle, will return the
	// referenced base name of the URI. For example, the base name of
	// 'file://foo/bar/baz.txt' is 'baz.txt'.
	//
	// If no implementation is provided, a generic one will be used. The
	// generic implementation will string-split the path component of the
	// URI on instances of the '/' character, and return the final such
	// component.
	//
	// Since 2.0.0
	NameImpl() func(URI) string

	// MimeTypeImpl returns a method that given a URI matching the
	// scheme this repository is registered to handle, will return the
	// MIME type of the resource the URI references.
	//
	// If no implementation is provided, a generic one will be used which
	// simply returns `application/octet-stream`, which is appropriate for
	// use when the MIME type is otherwise unknown per IETF RFC2046, pp.
	// 12.
	//
	// Since 2.0.0
	MimeTypeImpl() func(URI) (string, error)

	// ParentImpl returns a method that given a URI 'X' matching the
	// scheme this repository is registered to handle, will return a new
	// URI 'Y' such that the resource referenced by 'Y' is a parent of the
	// resource referenced by 'X'. For example, the parent for
	// 'file://foo/bar' is 'file://foo'.
	//
	// The provided Parent implementation may fail in several ways:
	//
	// * Determining the parent requires a hardware, network, or OS
	//   operation which has failed for some reason.
	//
	// * Determining the parent requires an operation for which different
	//   permissions or credentials are required.
	//
	// * The repository can determine the parent of certain URIs of the
	//   relevant scheme, but not others.
	//
	// * 'X' references the root of the hierarchy of resources this URI
	//   scheme can reference (e.g. the root directory of a UNIX
	//   filesystem). NOTE: in such cases, the implementation MUST return a
	//   URIRootError. This is especially important, as some routines which
	//   traverse a hierarchy of referenced resources use this to determine
	//   when to stop iterating.
	//
	// If the boolean return value is 'false', then any attempt to
	// determine the parent of URIs with the scheme this repository is
	// registered to handle with an error indicating the operation is not
	// supported.
	//
	// If the boolean return value is 'true' and the returned method is
	// nil, then a generic implementation will be used instead, which
	// discards the query and fragment components of the URI, string splits
	// the path on '/' characters to remove the final component, then uses
	// this repository's URI parsing implementation to parse the result
	// into a new URI, which is used as the parent.
	//
	//
	// Since 2.0.0
	ParentImpl() (func(URI) (URI, error), bool)

	// ReaderFromImpl returns a method which given a URI matching the
	// scheme that this repository is registered to handle, will return a
	// URIReadCloser set up to read from the resource that the URI
	// references.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to read the
	//   referenced resource.
	//
	// * This URI scheme could represent some resources that can be read,
	//   but this particular URI references a resources that is not
	//   something that can be read.
	//
	// * Attempting to set up the reader depended on a lower level
	//   operation such as a network or filesystem access that has failed
	//   in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// then any attempt to get a reader for a URI with the scheme this
	// repository is meant to handle will fail with an 'not supported'
	// error.
	//
	// Since 2.0.0
	ReaderFromImpl() func(URI) (URIReadCloser, error)

	// WriterFromImpl returns a method which given a URI matching the
	// scheme that this repository is registered to handle, will return a
	// URIWriteCloser set up to write to the resource that the URI
	// references.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to write to the
	//   referenced resource.
	//
	// * This URI scheme could represent some resources that can be
	//   written, but this particular URI references a resources that is
	//   not something that can be written.
	//
	// * Attempting to set up the writer depended on a lower level
	//   operation such as a network or filesystem access that has failed
	//   in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// then any attempt to get a writer for a URI with the scheme this
	// repository is meant to handle will fail with an 'not supported'
	// error.
	//
	// Since 2.0.0
	WriterToImpl() func(URI) (URIWriteCloser, error)

	// CopyImpl returns a method that given two URIs, 'src', and 'dest'
	// both of the scheme this repository is registered to handle, will
	// copy one to the other.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to perform the
	//   copy operation.
	//
	// * This URI scheme could represent some resources that can be copied,
	//   but either the source, destination, or both are not resources
	//   that support copying.
	//
	// * Performing the copy operation depended on a lower level operation
	//   such as network or filesystem access that has failed in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// and a "false" value is returned by the boolean return value, then
	// any attempt to perform a copy operation for a URI with the scheme
	// this repository is meant to handle will fail with a 'not supported'
	// error. If a "true" value is returned by the boolean return value, a
	// default implementation will be used which creates a new URI for the
	// destination, reads from the source, and writes to the destination.
	//
	// NOTE: a generic implementation which can work for URIs of different
	// schemes is provided via storage.Duplicate().
	//
	// Since 2.0.0
	CopyImpl() (func(URI, URI) error, bool)

	// RenameImpl returns a method that given two URIs, 'src' and 'dest'
	// both of the scheme this repository is registered to handle, will
	// rename src to dest. This means the resource referenced by src will
	// be copied into the resource referenced by dest, and the resource
	// referenced by src will no longer exist after the operation is
	// complete.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to perform the
	//   rename operation.
	//
	// * This URI scheme could represent some resources that can be renamed,
	//   but either the source, destination, or both are not resources
	//   that support renaming.
	//
	// * Performing the rename operation depended on a lower level operation
	//   such as network or filesystem access that has failed in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// and a "false" value is returned by the boolean return value, then
	// any attempt to perform a rename operation for a URI with the scheme
	// this repository is meant to handle will fail with a 'not supported'
	// error. If a "true" value is returned by the boolean return value, a
	// default implementation will be used which uses this repository's
	// copy implementation, and then deletes the source.
	//
	// NOTE: a generic implementation which can work for URIs of different
	// schemes is provided via storage.Move().
	//
	// Since 2.0.0
	RenameImpl() (func(URI, URI) error, bool)

	// DeleteImpl returns a method that given a URI of the scheme this
	// repository is registered to handle. This method should cause the
	// resource referenced by the URI To be destroyed, removed, or
	// otherwise deleted.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to perform the
	//   delete operation.
	//
	// * This URI scheme could represent some resources that can be
	//   deleted, but this specific URI is not one of them.
	//
	// * Performing the delete operation depended on a lower level operation
	//   such as network or filesystem access that has failed in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// then any attempt to perform a delete operation on a URI with the
	// scheme that this repository is registered to handle will
	// fail with a 'not supported' error.
	//
	// Since 2.0.0
	DeleteImpl() func(URI) error

	// ListableImpl returns a method that given a URI of the scheme this
	// repository is registered to handle, will determine if it is listable
	// or not.
	//
	// NOTE: If ListeImpl() is not implemented in this repository, then
	// attempts to check the listability of URIs of the scheme it is
	// registered to handle will fail, even if ListableImpl() is
	// implemented.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to check if the
	//   URI supports listing.
	//
	// * This URI scheme could represent some resources that can be listed,
	//   but this specific URI is not one of them (e.g. a file on a
	//   filesystem, as opposed to a directory).
	//
	// * Checking for listability depended on a lower level operation
	//   such as network or filesystem access that has failed in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// then any attempt to check if the URI is listable will fail with a
	// 'not supported' error. This implies that no URI of the scheme that
	// this repository is registered to handle could ever be listable.
	//
	// Since 2.0.0
	ListableImpl() func(URI) (bool, error)

	// ListImpl returns a method that given a URI of the scheme that this
	// repository is registered to handle, will return a list of URIs that
	// reference resources which are nested below the resource referenced
	// by the argument. For example, listing a directory on a filesystem
	// should return a list of files and directories it contains.
	//
	// NOTE: If ListableImpl() is not implemented in this repository,
	// then attempts to list URIs of the scheme it is registered to handle
	// will fail, even if ListImpl() is implemented.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to obtain a
	//   listing for the given URI.
	//
	// * This URI scheme could represent some resources that can be listed,
	//   but this specific URI is not one of them (e.g. a file on a
	//   filesystem, as opposed to a directory).
	//
	// * Obtaining the listing depended on a lower level operation such as
	//   network or filesystem access that has failed in some way.
	//
	// If this method is not implemented (e.g. a nil function is returned),
	// then any attempt to list a URI of the scheme this repository is
	// registered to handle will fail with a 'not supported' error.
	//
	// Since 2.0.0
	ListImpl() func(URI) ([]URI, error)

	// ChildImpl returns a method that given a URI of the scheme
	// that this repository is registered to handle, create a new URI
	// nested under it. For example, if the method is called on the URI
	// file://foo/bar with the string argument 'baz.txt', then the
	// URI returned would be 'file:.//foo/bar/baz.txt'.
	//
	// The returned method may fail in several ways:
	//
	// * Different permissions or credentials are required to create a
	//   child of the given URI.
	//
	// * This URI scheme could represent some resources that can be have
	//   children, but this specific URI is not one of them (e.g. a file on
	//   a filesystem, as opposed to a directory).
	//
	// * Creating the child depended on a lower level operation such as
	//   network or filesystem access that has failed in some way.
	//
	// NOTE: the fact that a child can be created for a given URI implies
	// that the URI must be listable. If your repository is implemented
	// inconsistently with this, things may break in unexpected ways.
	//
	// If this method is not implemented and a boolean return value of
	// 'true' is given, then a default implementation will be used that
	// calls the given ParseURI implementation with the string appended to
	// the URI path component, separated by a '/' character. If 'false' is
	// given, then attempts to create children for any URI of the scheme
	// this repository is registered to handle will fail with a 'not
	// supported' error.
	//
	// Since 2.0.0
	ChildImpl() (func(URI, string) (URI, error), bool)
}
