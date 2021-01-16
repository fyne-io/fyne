package repository

import (
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne"
)

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
	scheme = strings.ToLower(scheme)

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

		// Windows files can break authority checks, so just return the parsed file URI
		return NewFileURI(path), nil
	}

	repo, err := ForURI(&uri{scheme: scheme})
	if err == nil {
		// If the repository registered for this scheme implements a parser
		if c, ok := repo.(CustomURIRepository); ok {
			return c.ParseURI(s)
		}
	}

	// There was no repository registered, or it did not provide a parser
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
