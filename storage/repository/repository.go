package repository

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"

	"fyne.io/fyne"
)

// repositoryTable stores the mapping of schemes to Repository implementations.
// It should only ever be used by ForURI() and Register().
var repositoryTable map[string]Repository = map[string]Repository{}

// Repository represents a storage repository, which is a set of methods which
// implement specific functions on a URI. Repositories are registered to handle
// specific URI schemes, and the higher-level functions that operate on URIs
// internally look up an appropriate method from the relevant Repository.
//
// The repository interface includes only methods which must be implemented at
// a minimum. Without implementing all of the methods in this interface, a URI
// would not be usable in a useful way. Additional functionality can be exposed
// by using interfaces which extend Repository.
//
// Repositories are registered to handle a specific URI scheme (or schemes)
// using the Register() method. When a higher-level URI function such as
// storage.Copy() is called, the storage package will internally look up
// the repository associated with the scheme of the URI, then it will use
// a type assertion to check if the repository implements CopyableRepository.
// If so, the Copy() function will be run from the repository, otherwise
// storage.Copy() will return NotSupportedError. This works similarly for
// all other methods in repository-related interfaces.
//
// Note that a repository can be registered for multiple URI schemes. In such
// cases, the repository must internally select and implement the correct
// behavior for each URI scheme.
//
// NOTE: most developers who use Fyne should *not* generally attempt to
// call repository methods directly. You should use the methods in the storage
// package, which will automatically detect the scheme of a URI and call into
// the appropriate repository.
//
// Since: 2.0.0
type Repository interface {

	// Exists will be used to implement calls to storage.Exists() for the
	// registered scheme of this repository.
	//
	// Since: 2.0.0
	Exists(u fyne.URI) (bool, error)

	// Reader will be used to implement calls to storage.Reader()
	// for the registered scheme of this repository.
	//
	// Since: 2.0.0
	Reader(u fyne.URI) (fyne.URIReadCloser, error)

	// CanRead will be used to implement calls to storage.CanRead() for the
	// registered scheme of this repository.
	//
	// Since: 2.0.0
	CanRead(u fyne.URI) (bool, error)

	// Destroy is called when the repository is un-registered from a given
	// URI scheme.
	//
	// The string parameter will be the URI scheme that the repository was
	// registered for. This may be useful for repositories that need to
	// handle more than one URI scheme internally.
	//
	// Since: 2.0.0
	Destroy(string)
}

// CanonicalRepository is an extension of the repository interface which
// allows the behavior of storage.ParseURI to be overridden. This is only
// needed if you wish to generate custom URI types, rather than using Fyne's
// URI implementation and net/url based parsing.
//
// NOTE: even for URIs with non-RFC3986-compliant encoding, the URI MUST begin
// with 'scheme:', or storage.ParseURI() will not be able to determine which
// storage repository to delegate to for parsing.
//
// Since: 2.0.0
type CanonicalRepository interface {
	Repository

	// ParseURI will be used to implement calls to storage.ParseURI()
	// for the registered scheme of this repository.
	ParseURI(string) (fyne.URI, error)
}

// WriteableRepository is an extension of the Repository interface which also
// supports obtaining a writer for URIs of the scheme it is registered to.
//
// Since: 2.0.0
type WriteableRepository interface {
	Repository

	// Writer will be used to implement calls to storage.WriterTo() for
	// the registered scheme of this repository.
	//
	// Since: 2.0.0
	Writer(u fyne.URI) (fyne.URIWriteCloser, error)

	// CanWrite will be used to implement calls to storage.CanWrite() for
	// the registered scheme of this repository.
	//
	// Since: 2.0.0
	CanWrite(u fyne.URI) (bool, error)

	// Delete will be used to implement calls to storage.Delete() for the
	// registered scheme of this repository.
	//
	// Since: 2.0.0
	Delete(u fyne.URI) error
}

// ListableRepository is an extension of the Repository interface which also
// supports obtaining directory listings (generally analogous to a directory
// listing) for URIs of the scheme it is registered to.
//
// Since: 2.0.0
type ListableRepository interface {
	Repository

	// CanList will be used to implement calls to storage.Listable() for
	// the registered scheme of this repository.
	//
	// Since: 2.0.0
	CanList(u fyne.URI) (bool, error)

	// List will be used to implement calls to storage.List() for the
	// registered scheme of this repository.
	//
	// Since: 2.0.0
	List(u fyne.URI) ([]fyne.URI, error)
}

// HierarchicalRepository is an extension of the Repository interface which
// also supports determining the parent and child items of a URI.
//
// Since: 2.0.0
type HierarchicalRepository interface {
	Repository

	// Parent will be used to implement calls to storage.Parent() for the
	// registered scheme of this repository.
	//
	// A generic implementation is provided in GenericParent(), which
	// is based on the RFC3986 definition of a URI parent.
	//
	// Since: 2.0.0
	Parent(fyne.URI) (fyne.URI, error)

	// Child will be used to implement calls to storage.Child() for
	// the registered scheme of this repository.
	//
	// A generic implementation is provided in GenericParent(), which
	// is based on RFC3986.
	//
	// Since: 2.0.0
	Child(fyne.URI, string) (fyne.URI, error)
}

// CopyableRepository is an extension of the Repository interface which also
// supports copying referenced resources from one URI to another.
//
// Since: 2.0.0
type CopyableRepository interface {
	Repository

	// Copy will be used to implement calls to storage.Copy() for the
	// registered scheme of this repository.
	//
	// A generic implementation is provided by GenericCopy().
	//
	// NOTE: the first parameter is the source, the second is the
	// destination.
	//
	// NOTE: if storage.Copy() is given two URIs of different schemes, it
	// is possible that only the source URI will be of the type this
	// repository is registered to handle. In such cases, implementations
	// are suggested to fail-over to GenericCopy().
	//
	// Since: 2.0.0
	Copy(fyne.URI, fyne.URI) error
}

// MovableRepository is an extension of the Repository interface which also
// supports moving referenced resources from one URI to another.
//
// Note: both Moveable and Movable are correct spellings, but Movable is newer
// and more accepted. Source: https://grammarist.com/spelling/movable-moveable/
//
// Since: 2.0.0
type MovableRepository interface {
	Repository

	// Move will be used to implement calls to storage.Move() for the
	// registered scheme of this repository.
	//
	// A generic implementation is provided by GenericMove().
	//
	// NOTE: the first parameter is the source, the second is the
	// destination.
	//
	// NOTE: if storage.Move() is given two URIs of different schemes, it
	// is possible that only the source URI will be of the type this
	// repository is registered to handle. In such cases, implementations
	// are suggested to fail-over to GenericMove().
	//
	// Since: 2.0.0
	Move(fyne.URI, fyne.URI) error
}

// Register registers a storage repository so that operations on URIs of the
// registered scheme will use methods implemented by the relevant repository
// implementation.
//
// Since: 2.0.0
func Register(scheme string, repository Repository) {
	prev, ok := repositoryTable[scheme]

	if ok {
		prev.Destroy(scheme)
	}

	repositoryTable[scheme] = repository
}

// ForURI returns the Repository instance which is registered to handle URIs of
// the given scheme.
//
// NOTE: this function is intended to be used specifically by the storage
// package. It generally should not be used outside of the fyne package -
// instead you should use the methods in the storage package.
//
// Since: 2.0.0
func ForURI(u fyne.URI) (Repository, error) {
	scheme := u.Scheme()
	repo, ok := repositoryTable[scheme]

	if !ok {
		return nil, fmt.Errorf("no repository registered for scheme '%s'", scheme)
	}

	return repo, nil
}

// ParseURI implements the back-end logic for storage.ParseURI, which you
// should use instead. This is only here because other functions in repository
// need to call it, and it prevents a circular import.
//
// Since: 2.0.0
func ParseURI(s string) (fyne.URI, error) {
	// Extract the scheme.
	scheme := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			break
		}
		scheme += string(s[i])
	}

	if scheme == "file" {
		// Does this really deserve to be special? In principle, the
		// purpose of this check is to pass it to NewFileURI, which
		// allows platform path seps in the URI (against the RFC, but
		// easier for people building URIs naively on Windows). Maybe
		// we should punt this to whoever generated the URI in the
		// first place?

		path := s[5:] // everything after file:
		if len(path) > 2 && path[:2] == "//" {
			path = path[2:]
		}

		// this looks weird, but it makes sure that we still pass
		// url.Parse()
		s = NewFileURI(path).String()
	}

	repo, err := ForURI(&uri{scheme: scheme})
	crepo, ok := repo.(CanonicalRepository)
	if !ok || err != nil {
		// Either the repository registered for this scheme does not
		// implement a parser, or there is no registered repository.
		// Either way, we'll do our best-effort default.
		l, err := url.Parse(s)
		if err != nil {
			return nil, err
		}

		return &uri{
			scheme:    l.Scheme,
			authority: l.User.String() + l.Host,
			// workaround for net/url, see type uri struct comments
			haveAuthority: true,
			path:          l.Path,
			query:         l.RawQuery,
			fragment:      l.Fragment,
		}, nil
	}

	return crepo.ParseURI(s)

}

// NewFileURI implements the back-end logic to storage.NewFileURI, which you
// should use instead. This is only here because other functions in repository
// need to call it, and it prevents a circular import.
//
// Since: 2.0.0
func NewFileURI(path string) fyne.URI {
	// URIs are supposed to use forward slashes. On Windows, it
	// should be OK to use the platform native filepath with UNIX
	// or NT style paths, with / or \, but when we reconstruct
	// the URI, we want to have / only.
	if runtime.GOOS == "windows" {
		// seems that sometimes we end up with
		// double-backslashes
		path = filepath.ToSlash(path)
	}

	return &uri{
		scheme:        "file",
		haveAuthority: true,
		authority:     "",
		path:          path,
		query:         "",
		fragment:      "",
	}
}
