package repository

import (
	"bufio"
	"mime"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
)

// Declare conformance with fyne.URI interface.
var _ fyne.URI = &uri{}

type uri struct {
	scheme    string
	authority string
	// haveAuthority lets us distinguish between a present-but-empty
	// authority, and having no authority. This is needed because net/url
	// incorrectly handles scheme:/absolute/path URIs.
	haveAuthority bool
	path          string
	query         string
	fragment      string
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

	return strings.Split(mimeTypeFull, ";")[0]
}

func (u *uri) Scheme() string {
	return u.scheme
}

func (u *uri) String() string {
	// NOTE: this string reconstruction is mandated by IETF RFC3986,
	// section 5.3, pp. 35.

	s := u.scheme + ":"
	if u.haveAuthority {
		s += "//" + u.authority
	}
	s += u.path
	if len(u.query) > 0 {
		s += "?" + u.query
	}
	if len(u.fragment) > 0 {
		s += "#" + u.fragment
	}
	return s
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
