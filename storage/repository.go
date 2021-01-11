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

	// ReaderFrom will be used to implement calls to storage.ReaderFrom()
	// for the registered scheme of this repository.
	//
	// Since 2.0.0
	ReaderFrom(u fyne.URI) (fyne.URIReadCloser, error)

	// CanRead will be used to implement calls to storage.CanRead() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	CanRead(u fyne.URI) (bool, error)
}

// WriteableRepository is an extension of the Repository interface which also
// supports obtaining a writer for URIs of the scheme it is registered to.
//
// If this interface is not implemented for a repository, then attempts to use
// storage.Writer, storage.Delete, and storage.CanWrite will fail with
// URIOperationNotSupportedError when called on URIs of the scheme this is
// registered to.
//
// Since: 2.0.0
type WriteableRepository interface {
	Repository

	// Writer will be used to implement calls to storage.WriterTo() for
	// the registered scheme of this repository.
	//
	// Since 2.0.0
	Writer(u fyne.URI) (fyne.URIWriteCloser, error)

	// CanWrite will be used to implement calls to storage.CanWrite() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	CanWrite(u fyne.URI) (bool, error)

	// Delete will be used to implement calls to storage.Delete() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	Delete(u fyne.URI) error
}

// ListableRepository is an extension of the Repository interface which also
// supports obtaining directory listings (generally analogous to a directory
// listing) for URIs of the scheme it is registered to.
//
// If this interface is not implemented for a repository, then attempts to use
// storage.CanList and storage.List will fail with
// URIOperationNotSupportedError when called for URIs of the scheme it is
// registered to.
//
// Since: 2.0.0
type ListableRepository interface {
	Repository

	// CanList will be used to implement calls to storage.Listable() for
	// the registered scheme of this repository.
	//
	// Since 2.0.0
	CanList(u fyne.URI) (bool, error)

	// List will be used to implement calls to storage.List() for the
	// registered scheme of this repository.
	//
	// Since 2.0.0
	List(u fyne.URI) ([]URI, error)
}

// gets fallbacks
type HierarchicalRepository interface {
	Repository

	Parent(fyne.URI) (fyne.URI, error)

	Child(Fyne.URI) (fyne.URI, error)
}

// gets fallbacks
type CopyableRepository interface {
	Repository

	Copy(fyne.URI, fyne.URI) error
}

// gets fallbacks
type MovableRepository interface {
	Repository

	Move(fyne.URI, fyne.URI) error
}

// Register registers a storage repository so that operations on URIs of the
// registered scheme will use methods implemented by the relevant repository
// implementation.
//
// Since 2.0.0
func Register(scheme string, repository Repository) {
}
