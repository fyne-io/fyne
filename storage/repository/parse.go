package repository

import (
	"errors"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/goos"
)

const domainLabelPattern = "[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?"

var rxHostName = regexp.MustCompile("^" + domainLabelPattern + `(?:\.` + domainLabelPattern + ")*$")

// NewFileURI implements the back-end logic to storage.NewFileURI, which you
// should use instead. This is only here because other functions in repository
// need to call it, and it prevents a circular import.
//
// Since: 2.0
func NewFileURI(path string) fyne.URI {
	// URIs are supposed to use forward slashes. On Windows, it
	// should be OK to use the platform native filepath with UNIX
	// or NT style paths, with / or \, but when we reconstruct
	// the URI, we want to have / only.
	if runtime.GOOS == goos.Windows {
		// seems that sometimes we end up with
		// double-backslashes
		path = filepath.ToSlash(path)
	}

	return &uri{url.URL{
		Scheme: fyne.URISchemeFile,
		Path:   path,
	}}
}

// ParseURI implements the back-end logic for storage.ParseURI, which you
// should use instead. This is only here because other functions in repository
// need to call it, and it prevents a circular import.
//
// Since: 2.0
func ParseURI(s string) (fyne.URI, error) {
	// Extract the scheme.
	scheme, path, ok := strings.Cut(s, fyne.URISchemeSeparator)
	if !ok {
		return nil, errors.New("invalid URI, scheme must be present")
	}

	if strings.EqualFold(scheme, "urn") {
		f := append(strings.SplitN(path, "#", 2), "")
		q := append(strings.SplitN(f[0], "?", 2), "")
		return &uri{url.URL{
			Scheme:   scheme,
			Opaque:   q[0],
			RawQuery: q[1],
			Fragment: f[1],
		}}, nil
	}

	if runtime.GOOS == goos.Windows && len(scheme) == 1 {
		path = scheme + ":" + filepath.ToSlash(path)
		scheme = fyne.URISchemeFile
	}

	if strings.EqualFold(scheme, fyne.URISchemeFile) {
		path = strings.TrimPrefix(path, fyne.URIAuthorityPrefix)
		if path == "" {
			return nil, errors.New("invalid file URI, path cannot be empty")
		}
		p, err := url.PathUnescape(path)
		if err != nil {
			return nil, err
		}
		return NewFileURI(p), nil
	}

	scheme = strings.ToLower(scheme)
	repo, err := ForScheme(scheme)
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

	if l.Host == "" {
		return &uri{*l}, nil
	}

	host := l.Hostname()
	if net.ParseIP(host) != nil {
		return &uri{*l}, nil
	}
	if !rxHostName.MatchString(host) {
		return nil, errors.New("failed to validate host")
	}
	return &uri{*l}, nil
}
