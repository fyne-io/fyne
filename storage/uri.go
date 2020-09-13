package storage

import (
	"bufio"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne"
)

type uri struct {
	raw string
}

// NewURI creates a new URI from the given string representation.
// This could be a URI from an external source or one saved from URI.String()
func NewURI(u string) fyne.URI {
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

// Parent gets the parent of a URI by splitting it along '/' separators and
// removing the last item.
func Parent(u fyne.URI) (fyne.URI, error) {
	s := u.String()

	// trim trailing slash
	if s[len(s)-1] == '/' {
		s = s[0 : len(s)-1]
	}

	// trim the scheme (and +1 for the :)
	s = s[len(u.Scheme())+1 : len(s)]

	// Completely empty URI with just a scheme
	if len(s) == 0 {
		return nil, URIRootError
	}

	// trim leading forward slashes
	trimmed := 0
	for s[0] == '/' {
		s = s[1:len(s)]
		trimmed++

		// if all we have left is an empty string, than this URI
		// pointed to a UNIX-style root
		if len(s) == 0 {
			return nil, URIRootError
		}
	}

	// handle Windows drive letters
	r := regexp.MustCompile("[A-Za-z][:]")
	components := strings.Split(s, "/")
	if len(components) == 1 && r.MatchString(components[0]) && trimmed <= 2 {
		// trimmed <= 2 makes sure we handle UNIX-style paths on
		// Windows correctly
		return nil, URIRootError
	}

	parent := u.Scheme() + "://"
	if trimmed > 2 && len(components) > 1 {
		// Because we trimmed all the leading '/' characters, for UNIX
		// style paths we want to insert one back in. Presumably we
		// trimmed two instances of / for the scheme.
		parent = parent + "/"
	}
	parent = parent + strings.Join(components[0:len(components)-1], "/") + "/"
	return NewURI(parent), nil
}

// Child appends a new path element to a URI, separated by a '/' character.
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
func Exists(u fyne.URI) (bool, error) {
	if u.Scheme() != "file" {
		return false, fmt.Errorf("Don't know how to check existence of %s scheme", u.Scheme())
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
