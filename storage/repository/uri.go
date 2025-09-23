package repository

import (
	"bufio"
	"mime"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
)

// EqualURI returns true if the two URIs are equal.
//
// Since: 2.6
func EqualURI(t1, t2 fyne.URI) bool {
	if t1 == nil || t2 == nil {
		return t1 == t2
	}

	u1, ok1 := t1.(*uri)
	u2, ok2 := t2.(*uri)
	if ok1 && ok2 {
		// Knowing the type, pointers are either the same or fields are the same.
		return u1 == u2 || *u1 == *u2
	}

	return t1 == t2 || t1.String() == t2.String()
}

// Declare conformance with fyne.URI interface.
var _ fyne.URI = &uri{}

type uri struct {
	scheme    string
	authority string
	path      string
	query     string
	fragment  string
}

func (u *uri) Extension() string {
	return filepath.Ext(u.path)
}

func (u *uri) Name() string {
	return filepath.Base(u.path)
}

func (u *uri) MimeType() string {
	mimeTypeFull := mime.TypeByExtension(u.Extension())
	if mimeTypeFull == "" {
		mimeTypeFull = "text/plain"

		repo, err := ForURI(u)
		if err != nil {
			return "application/octet-stream"
		}

		readCloser, err := repo.Reader(u)
		if err == nil {
			defer readCloser.Close()
			scanner := bufio.NewScanner(readCloser)
			if scanner.Scan() && !utf8.Valid(scanner.Bytes()) {
				mimeTypeFull = "application/octet-stream"
			}
		}
	}

	mimeType, _, _ := strings.Cut(mimeTypeFull, ";")
	return mimeType
}

func (u *uri) Scheme() string {
	return u.scheme
}

func (u *uri) String() string {
	// NOTE: this string reconstruction is mandated by IETF RFC3986,
	// section 5.3, pp. 35.
	s := strings.Builder{}
	s.Grow(len(u.scheme) + len(u.authority) + len(u.path) + len(u.query) + len(u.fragment) + len("://?#"))

	s.WriteString(u.scheme)
	s.WriteString("://")
	s.WriteString(u.authority)
	s.WriteString(u.path)

	if len(u.query) > 0 {
		s.WriteByte('?')
		s.WriteString(u.query)
	}
	if len(u.fragment) > 0 {
		s.WriteByte('#')
		s.WriteString(u.fragment)
	}
	return s.String()
}

func (u *uri) Authority() string {
	return u.authority
}

func (u *uri) Path() string {
	return u.path
}

func (u *uri) Query() string {
	return u.query
}

func (u *uri) Fragment() string {
	return u.fragment
}
