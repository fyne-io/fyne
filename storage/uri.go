package storage

import (
	"bufio"
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne"
	"fyne.io/fyne/storage/repository"
)

// Declare conformance with fyne.URI interface.
var _ fyne.URI = &uri{}

// For backwards-compatibility with the now-deprecated ListableURI type, we
// also declare conformance with that.
var _ fyne.ListableURI = &uri{}

type uri struct {
	raw string
}

// NewFileURI creates a new URI from the given file path.
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
	return &uri{raw: "file://" + path}
}

func (u *uri) Extension() string {
	return filepath.Ext(u.raw)
}

func (u *uri) Name() string {
	return filepath.Base(u.raw)
}

func (u *uri) MimeType() string {
	mimeTypeFull := mime.TypeByExtension(u.Extension())
	if mimeTypeFull == "" {
		mimeTypeFull = "text/plain"
		readCloser, err := fyne.CurrentApp().Driver().FileReaderForURI(u)
		if err == nil {
			defer readCloser.Close()
			scanner := bufio.NewScanner(readCloser)
			if scanner.Scan() && !utf8.Valid(scanner.Bytes()) {
				mimeTypeFull = "application/octet-stream"
			}
		}
	}

	return strings.Split(mimeTypeFull, ";")[0]
}

func (u *uri) Scheme() string {
	pos := strings.Index(u.raw, ":")
	if pos == -1 {
		return ""
	}

	return strings.ToLower(u.raw[:pos])
}

func (u *uri) String() string {
	return u.raw
}

func (u *uri) Authority() string {
	// NOTE: we verified in ParseURI() that this would not error.
	r, _ := url.Parse(u.raw)

	a := ""
	if len(r.User.String()) > 0 {
		a = r.User.String() + "@"
	}
	a = a + r.Host

	return a
}

func (u *uri) Path() string {
	// NOTE: we verified in ParseURI() that this would not error.
	r, _ := url.Parse(u.raw)

	return r.Path
}

func (u *uri) Query() string {
	// NOTE: we verified in ParseURI() that this would not error.
	r, _ := url.Parse(u.raw)

	return r.RawQuery
}

func (u *uri) Fragment() string {
	// NOTE: we verified in ParseURI() that this would not error.
	r, _ := url.Parse(u.raw)

	return r.Fragment
}

func (u *uri) List() ([]fyne.URI, error) {
	return List(u)
}

// NewURI creates a new URI from the given string representation. This could be
// a URI from an external source or one saved from URI.String()
//
// Deprecated - use ParseURI instead
func NewURI(s string) fyne.URI {
	u, _ := ParseURI(s)
	return u
}

// ParseURI creates a new URI instance by parsing a URI string, which must
// conform to IETF RFC3986.
//
// Since 2.0.0
func ParseURI(s string) (fyne.URI, error) {

	if len(s) > 5 && s[:5] == "file:" {
		path := s[5:]
		if len(path) > 2 && path[:2] == "//" {
			path = path[2:]
		}

		// this looks weird, but it makes sure that we still pass
		// url.Parse()
		s = NewFileURI(path).String()
	}

	_, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return &uri{raw: s}, nil
}

// Parent returns a URI referencing the parent resource of the resource
// referenced by the URI. For example, the Parent() of 'file://foo/bar.baz' is
// 'file://foo'. The URI which is returned will be listable.
//
// NOTE: it is not a given that Parent() return a parent URI with the same
// Scheme(), though this will normally be the case.
//
// This can fail in several ways:
//
// * If the URI refers to a filesystem root, then the Parent() implementation
//   must return (nil, URIRootError).
//
// * If the URI refers to a resource which does not exist in a hierarchical
//   context (e.g. the URI references something which does not have a
//   semantically meaningful "parent"), the Parent() implementation may return
//   an error.
//
// * If determining the parent of the referenced resource requires
//   interfacing with some external system, failures may propagate
//   through the Parent() implementation. For example if determining
//   the parent of a file:// URI requires reading information from
//   the filesystem, it could fail with a permission error.
//
// * If the scheme of the given URI does not have a registered
//   HierarchicalRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// NOTE: since 2.0.0, Parent() is backed by the repository system - this
// function is a helper which calls into an appropriate repository instance for
// the scheme of the URI it is given.
//
// Since: 1.4
func Parent(u fyne.URI) (fyne.URI, error) {

	repo, err := repository.ForURI(u)
	if err != nil {
		return nil, err
	}

	hrepo, ok := repo.(repository.HierarchicalRepository)
	if !ok {
		return nil, repository.OperationNotSupportedError
	}

	return hrepo.Parent(u)
}

// Child returns a URI referencing a resource nested hierarchically below the
// given URI, identified by a string. For example, the child with the string
// component 'quux' of 'file://foo/bar' is 'file://foo/bar/quux'.
//
// This can fail in several ways:
//
// * If the URI refers to a resource which does not exist in a hierarchical
//   context (e.g. the URI references something which does not have a
//   semantically meaningful "child"), the Child() implementation may return an
//   error.
//
// * If generating a reference to a child of the referenced resource requires
//   interfacing with some external system, failures may propagate through the
//   Child() implementation. It is expected that this case would occur very
//   rarely if ever.
//
// * If the scheme of the given URI does not have a registered
//   HierarchicalRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// NOTE: since 2.0.0, Child() is backed by the repository system - this
// function is a helper which calls into an appropriate repository instance for
// the scheme of the URI it is given.
//
// Since: 1.4
func Child(u fyne.URI, component string) (fyne.URI, error) {
	// While as implemented this does not need to return an error, it is
	// reasonable to expect that future implementations of this, especially
	// once it gets moved into the URI interface will need to do so. This
	// also brings it in line with Parent().

	s := u.String()

	// guarantee that there will be a path separator
	if s[len(s)-1:] != "/" {
		s += "/"
	}

	return ParseURI(s + component)
}

// Exists determines if the resource referenced by the URI exists.
//
// This can fail in several ways:
//
// * If checking the existence of a resource requires interfacing with some
//   external system, then failures may propagate through Exists(). For
//   example, checking the existence of a resource requires reading a directory
//   may result in a permissions error.
//
// It is understood that a non-nil error value signals that the existence or
// non-existence of the resource cannot be determined and is undefined.
//
// NOTE: since 2.0.0, Exists is backed by the repository system - this function
// calls into a scheme-specific implementation from a registered repository.
//
// may call into either a generic implementation, or into a scheme-specific
// implementation depending on which storage repositories have been registered.
//
// Since: 1.4
func Exists(u fyne.URI) (bool, error) {
	repo, err := repository.ForURI(u)
	if err != nil {
		return false, err
	}

	return repo.Exists(u)

	// TODO: this needs to move to the file:// repository
	// if u.Scheme() != "file" {
	//         return false, fmt.Errorf("don't know how to check existence of %s scheme", u.Scheme())
	// }
	//
	// _, err := os.Stat(u.String()[len(u.Scheme())+3:])
	// if os.IsNotExist(err) {
	//         return false, nil
	// }
	//
	// if err != nil {
	//         return false, err
	// }
	//
	// return true, nil
}

// Delete destroys, deletes, or otherwise removes the resource referenced
// by the URI.
//
// This can fail in several ways:
//
// * If removing the resource requires interfacing with some external system,
//   failures may propagate through Destroy(). For example, deleting a file may
//   fail with a permissions error.
//
// * If the referenced resource does not exist, attempting to destroy it should
//   throw an error.
//
// * If the scheme of the given URI does not have a registered
//   WriteableRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// Delete is backed by the repository system - this function calls
// into a scheme-specific implementation from a registered repository.
//
// Since: 2.0.0
func Delete(u fyne.URI) error {
	repo, err := repository.ForURI(u)
	if err != nil {
		return err
	}

	wrepo, ok := repo.(repository.WriteableRepository)
	if !ok {
		return repository.OperationNotSupportedError
	}

	return wrepo.Delete(u)

}

// Reader returns URIReadCloser set up to read from the resource that the
// URI references.
//
// This method can fail in several ways:
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
// Reader is backed by the repository system - this function calls
// into a scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	repo, err := repository.ForURI(u)
	if err != nil {
		return nil, err
	}

	return repo.Reader(u)
}

// CanRead determines if a given URI could be written to using the Reader()
// method. It is preferred to check if a URI is writable using this method
// before calling Reader(), because the underlying operations required to
// attempt to write and then report an error may be slower than the operations
// needed to test if a URI is writable.
//
// CanRead is backed by the repository system - this function calls into a
// scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func CanRead(u fyne.URI) (bool, error) {
	repo, err := repository.ForURI(u)
	if err != nil {
		return false, err
	}

	return repo.CanRead(u)
}

// Writer returns URIWriteCloser set up to write to the resource that the
// URI references.
//
// This method can fail in several ways:
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
// * If the scheme of the given URI does not have a registered
//   WriteableRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// Writer is backed by the repository system - this function calls into a
// scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	repo, err := repository.ForURI(u)
	if err != nil {
		return nil, err
	}

	wrepo, ok := repo.(repository.WriteableRepository)
	if !ok {
		return nil, repository.OperationNotSupportedError
	}

	return wrepo.Writer(u)
}

// CanWrite determines if a given URI could be written to using the Writer()
// method. It is preferred to check if a URI is writable using this method
// before calling Writer(), because the underlying operations required to
// attempt to write and then report an error may be slower than the operations
// needed to test if a URI is writable.
//
// CanWrite is backed by the repository system - this function calls into a
// scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func CanWrite(u fyne.URI) (bool, error) {
	repo, err := repository.ForURI(u)
	if err != nil {
		return false, err
	}

	wrepo, ok := repo.(repository.WriteableRepository)
	if !ok {
		return false, repository.OperationNotSupportedError
	}

	return wrepo.CanWrite(u)
}

// Copy given two URIs, 'src', and 'dest' both of the same scheme , will copy
// one to the other.
//
// This method may fail in several ways:
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
// * If the scheme of the given URI does not have a registered
//   CopyableRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// Copy is backed by the repository system - this function calls into a
// scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func Copy(source fyne.URI, destination fyne.URI) error {
	return fmt.Errorf("TODO: implement this function")
}

// Move returns a method that given two URIs, 'src' and 'dest' both of the same
// scheme this will move src to dest. This means the resource referenced by
// src will be copied into the resource referenced by dest, and the resource
// referenced by src will no longer exist after the operation is complete.
//
// This method may fail in several ways:
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
// * If the scheme of the given URI does not have a registered
//   MoveableRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// Move is backed by the repository system - this function calls into a
// scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func Move(source fyne.URI, destination fyne.URI) error {
	return fmt.Errorf("TODO: implement this function")
}

// CanList will determine if the URI is listable or not.
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
// * If the scheme of the given URI does not have a registered
//   ListableRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// CanList is backed by the repository system - this function calls into a
// scheme-specific implementation from a registered repository.
//
// Since 2.0.0
func CanList(u fyne.URI) (bool, error) {
	return false, fmt.Errorf("TODO: implement this function")
}

// List returns a list of URIs that reference resources which are nested below
// the resource referenced by the argument. For example, listing a directory on
// a filesystem should return a list of files and directories it contains.
//
// The returned method may fail in several ways:
//
// * Different permissions or credentials are required to obtain a
//   listing for the given URI.
//
// * This URI scheme could represent some resources that can be listed,
//   but this specific URI is not one of them (e.g. a file on a
//   filesystem, as opposed to a directory). This can be tested in advance
//   using the Listable() function.
//
// * Obtaining the listing depended on a lower level operation such as
//   network or filesystem access that has failed in some way.
//
// * If the scheme of the given URI does not have a registered
//   ListableRepository instance, then this method will fail with a
//   repository.OperationNotSupportedError.
//
// List is backed by the repository system - this function either calls into a
// scheme-specific implementation from a registered repository, or fails with a
// URIOperationNotSupported error.
//
// Since 2.0.0
func List(u fyne.URI) ([]fyne.URI, error) {
	return nil, fmt.Errorf("TODO")
}
