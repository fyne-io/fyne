package storage

import (
	"bufio"
	"mime"
	"os"
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

func (u *uri) MimeType() string {
	ext := filepath.Ext(u.raw[1:])
	mimeTypeFull := mime.TypeByExtension(ext)
	if mimeTypeFull == "" {
		mimeTypeFull = "text/plain"
		file, err := os.Open(u.raw)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			if scanner.Scan() && !utf8.Valid(scanner.Bytes()) {
				mimeTypeFull = "application/octet-stream"
			}
		}
	}
	mimeSubTypeSplit := strings.Split(mimeTypeFull, ";")
	mimeTypeFull = mimeSubTypeSplit[0]

	return mimeTypeFull
}

func (u *uri) Scheme() string {
	pos := strings.Index(u.raw, ":")
	if pos == -1 {
		return ""
	}

	return u.raw[:pos]
}

func (u *uri) String() string {
	return u.raw
}
