package repository

import (
	"bufio"
	"mime"
	"net/url"
	"path"
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
	url.URL
}

func (u *uri) Extension() string {
	return path.Ext(u.URL.Path)
}

func (u *uri) Name() string {
	return path.Base(u.URL.Path)
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
	return u.URL.Scheme
}

func (u *uri) Authority() string {
	if u.User != nil {
		s := u.User.String()
		if s != "" {
			return s + "@" + u.Host
		}
	}
	return u.Host
}

func (u *uri) Path() string {
	return u.URL.Path
}

func (u *uri) Query() string {
	return u.RawQuery
}

func (u *uri) Fragment() string {
	return u.URL.Fragment
}
