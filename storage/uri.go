package storage

import (
	"bufio"
	"mime"
	"path/filepath"
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
		s = s[0 : len(s)-2]
	}

	// We want to specifically make sure you don't take the parent of root.
	//
	// Note that we compare components[0] against "" since strings.Split()
	// will strip out the single / on root.
	//
	// Also note that components[0] will be something like 'file:'.
	components := strings.Split(s, "/")
	if (len(components) == 2 && components[1] == "") || (len(components) == 1) {
		return nil, URIRootError
	}

	parent := strings.Join(components[0:len(components)-1], "/") + "/"
	return NewURI(parent), nil
}
