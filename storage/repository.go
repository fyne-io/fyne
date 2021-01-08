package storage

import (
	"fyne.io/fyne"
)

// Repository represents a storage repository, which is a set of methods which
// implement specific functions on a URI. Repositories are registered to handle
// specific URI schemes, and the higher-level functions that operate on URIs
// internally look up an appropriate method from the relevant Repository.
//
// The repository interface includes only methods which must be implemented at
// a minimum. Without implementing all of the methods in this interface, a URI
// would not be usable in a useful way. Some additional methods which can offer
// more optimization opportunities can also be implemented by calling the
// appropriate Register...() function inside of the Repository.Init() function.
// For example, storage.Copy() will by default use a reader and a writer to
// perform a copy operation, but a repository might use RegisterCopy() to cause
// a more optimal version to be used for URIs of the relevant scheme (e.g. by
// performing the copy on the remote side of an ssh connection).
//
// In some cases, you may wish to omit an implementation of one of these
// functions anyway. In such cases, you should do this by returning a
// URIOperationNotSupportedError. For example, this may be appropriate when
// building a storage repository which represents only resources that are
// read-only - in such cases functions like WriterTo() may not make sense.
// Please use care when deciding to do this, and consider how it will impact
// users of your storage repository.
//
// NOTE: most developers who use Fyne should *not* generally attempt to
// call repository methods directly. You should use the methods in the storage
// package, which will automatically detect the scheme of a URI and call into
// the appropriate repository.
//
// NOTE: functions in a particular repository can assume they will only ever be
// used on URIs which match the scheme that they have been registered to
// handle, however they should not assume that they will only be called on URIs
// of the same implementation as the ParseURI() for this repository has
// returned.
//
// Since: 2.0.0
type Repository interface {

	// ParseURI returns a method that given a text string creates a
	// new URI instance. This may cause an error, for example if the string
	// is an invalid URI or contains information that is invalid in the
	// context of the repository's scheme.
	//
	// This method will be called only after the URI scheme has been
	// verified to match the one the repository was registered to handle,
	// however the URI is not validated further than this. It is required
	// by convention (but not enforced by technical means) that this
	// function should throw an error if the string is not a valid IETF RFC
	// 3986 complaint URI. Fail to validate this at your own risk.
	//
	// NOTE: it is highly recommended to use net/url (
	// https://golang.org/pkg/net/url/ ) for parsing URI strings unless you
	// have a good reason not to. It is mature and implements RFC 3986.
	//
	// If no implementation is provided, a generic URI based on RFC3986 is
	// parsed without any special validation logic.
	//
	// Since 2.0.0
	ParseURI() func(string) (fyne.URI, error)

	// Init will be called while the repository is being registered.  This
	// is the appropriate place to register any optional handlers relevant
	// to a particular URI scheme.
	//
	// Since 2.0.0
	Init() func() error

	// Exists will be used to implement calls to storage.Exists() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	Exists(u fyne.URI) (bool, error)

	// Delete will be used to implement calls to storage.Delete() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	Delete(u fyne.URI) error

	// ReaderFrom will be used to implement calls to storage.ReaderFrom()
	// for the registered scheme of this repository.
	//
	// Since 2.0.0
	ReaderFrom(u fyne.URI) (fyne.URIReadCloser, error)

	// WriterTo will be used to implement calls to storage.WriterTo() for
	// the registered scheme of this repository.
	//
	// Since 2.0.0
	WriterTo(u fyne.URI) (fyne.URIWriteCloser, error)

	// Listable will be used to implement calls to storage.Listable() for
	// the registered scheme of this repository.
	//
	// Since 2.0.0
	Listable(u fyne.URI) (bool, error)

	// List will be used to implement calls to storage.List() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	List(u fyne.URI) ([]URI, error)
}

// RegisterCopy registers an implementation of a Copy() function for a specific
// URI scheme. This function should only be called in Repository.Init(). If
// the scheme does not have a registered repository, then this function will
// fail with an error. If a Copy implementation is already registered for this
// scheme, it will be replaced with this one silently.
//
// If no Copy implementation is registered for a particular scheme, then copies
// will be implemented as if the URIs were of different types, that is the
// source URI will be read using ReaderFrom(), and the destination written
// using WriterTo(). This function is an opportunity to implement a more
// optimal approach, such as leveraging remote operations for a network-backed
// repository.
//
// Since 2.0.0
func RegisterCopy(scheme string, copyImplementation func(fyne.URI, fyne.URI) error) error {
	// TODO
}

// RegisterRename registers an implementation of a Rename() function for a
// specific URI scheme. This function should only be called in
// Repository.Init(). If the scheme does not have a registered repository, then
// this function will fail with an error. If a Rename implementation is already
// registered for this scheme, it will be replaced with this one silently.
//
// If no Rename implementation is registered for a particular scheme, then
// renames will be implemented as if the URIs were of different types, that is
// the source URI will be read using ReaderFrom(), and the destination written
// using WriterTo(), then the source will be deleted using Delete(). This
// function is an opportunity to implement a more optimal approach, such as
// leveraging remote operations for a network-backed repository.
//
// Since 2.0.0
func RegisterRename(scheme string, renameImplementation func(fyne.URI, fyne.URI) error) error {
}

// RegisterRepository registers a storage repository so that operations on URIs
// of the registered scheme will use methods implemented by the relevant
// repository implementation. This method will call repository.Init(), and may
// error if it does.
//
// Since 2.0.0
func RegisterRepository(scheme string, repository Repository) error {
}
