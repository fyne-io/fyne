package storage

import (
	"bufio"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne"
)

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

// NewURI creates a new URI from the given string representation.
// This could be a URI from an external source or one saved from URI.String()
func NewURI(u string) fyne.URI {
	if len(u) > 5 && u[:5] == "file:" {
		path := u[5:]
		if len(path) > 2 && path[:2] == "//" {
			path = path[2:]
		}
		return NewFileURI(path)
	}

	return &uri{raw: u}
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

// parentGeneric is a generic function that returns the last element of a
// path after splitting it on "/". It should be suitable for most URIs.
func parentGeneric(location string) (string, error) {

	// trim leading forward slashes
	trimmed := 0
	for location[0] == '/' {
		location = location[1:]
		trimmed++

		// if all we have left is an empty string, than this URI
		// pointed to a UNIX-style root
		if len(location) == 0 {
			return "", URIRootError
		}
	}

	components := strings.Split(location, "/")

	if len(components) == 1 {
		return "", URIRootError
	}

	parent := ""
	if trimmed > 2 && len(components) > 1 {
		// Because we trimmed all the leading '/' characters, for UNIX
		// style paths we want to insert one back in. Presumably we
		// trimmed two instances of / for the scheme.
		parent = parent + "/"
	}
	parent = parent + strings.Join(components[0:len(components)-1], "/") + "/"

	return parent, nil
}

// Parent returns a URI referencing the parent resource of the resource
// referenced by the URI. For example, the Parent() of 'file://foo/bar.baz' is
// 'file://foo'. The URI which is returned will be listable.
//
// NOTE: it is not required that the implementation return a parent URI with
// the same Scheme(), though this will normally be the case.
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
// NOTE: since 2.0.0, Parent() is backed by the repository system - this
// function may call into either a generic implementation, or into a
// scheme-specific implementation depending on which storage repositories have
// been registered.
//
// Since: 1.4
func Parent(u fyne.URI) (fyne.URI, error) {
	s := u.String()

	// trim trailing slash
	if s[len(s)-1] == '/' {
		s = s[0 : len(s)-1]
	}

	// trim the scheme
	s = s[len(u.Scheme())+3:]

	// Completely empty URI with just a scheme
	if len(s) == 0 {
		return nil, URIRootError
	}

	parent := ""
	if u.Scheme() == "file" {
		// use the system native path resolution
		parent = filepath.Dir(s)
		if parent[len(parent)-1] != filepath.Separator {
			parent += "/"
		}

		// only root is it's own parent
		if filepath.Clean(parent) == filepath.Clean(s) {
			return nil, URIRootError
		}

	} else {
		var err error
		parent, err = parentGeneric(s)
		if err != nil {
			return nil, err
		}
	}

	return NewURI(u.Scheme() + "://" + parent), nil
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
// NOTE: since 2.0.0, Child() is backed by the repository system - this
// function may call into either a generic implementation, or into a
// scheme-specific implementation depending on which storage repositories have
// been registered.
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

	return NewURI(s + component), nil
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
// Since: 1.4
func Exists(u fyne.URI) (bool, error) {
	if u.Scheme() != "file" {
		return false, fmt.Errorf("don't know how to check existence of %s scheme", u.Scheme())
	}

	_, err := os.Stat(u.String()[len(u.Scheme())+3:])
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Destroy destroys, deletes, or otherwise removes the resource referenced
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
// Since: 2.0.0
func Destroy(u fyne.URI) error {
	return fmt.Errorf("TODO: implement this function")
}
