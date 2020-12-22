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
		path = strings.ReplaceAll(path, "\\\\", "/")
		path = strings.ReplaceAll(path, "\\", "/")
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

// Parent gets the parent of a URI by splitting it along '/' separators and
// removing the last item.
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

// Child appends a new path element to a URI, separated by a '/' character.
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

// Exists will return true if the resource the URI refers to exists, and false
// otherwise. If an error occurs while checking, false is returned as the first
// return.
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
