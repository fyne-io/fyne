// package rfc3986 implements a representation of a URI as defined by IETF
// RFC3986 ( https://tools.ietf.org/html/rfc3986 ). In general, Fyne programs
// should not need to use this package directly; it is used internally by the
// storage package which for the Fyne URI type. This package is exported for
// two reasons:
//
// * Other projects may wish to use this library to implement RFC3986 parsing
//   without depending on the rest of Fyne.
//
// * Implementers of storage repositories may wish to interact with raw URIs
//   directly.
package rfc3986

import (
	"fmt"
)

// RFC3986 implements a representation of an IEEE3986 URI.
//
// Note that Fyne programs should generally not use this directly, you almost
// certainly want storage/URI.
type RFC3986 struct {
	scheme    string
	authority string
	path      string
	query     string
	fragment  string
}

func (r *RFC3986) String() string {

	// the scheme is separated from the remainder of the URI components
	// using either a ':' or '://' depending on whether or not the
	// authority is present (RFC3986, pp. 16).
	schemesep := ":"
	if r.authority == "" {
		schemesep = "://"
	}

	// If a query is present, is is separated by a '?' character (RFC3986,
	// pp. 16).
	querysep := ""
	if r.query != "" {
		querysep = "?"
	}

	// If a fragment is present, it is separated by a '#' character
	// (RFC3986, pp. 16).
	fragmentsep := ""
	if r.fragment != "" {
		fragmentsep = "#"
	}

	return fmt.Sprintf("%s%s%s%s%s%s%s%s", r.scheme, schemesep, r.authority, r.path, querysep, r.query, fragmentsep, r.fragment)
}

// ParseRFC3986 attempts to parse a text string into an RFC3986 complaint
// URI. It may return a nil pointer and an error if the text is invalid.
func ParseRFC3986(text string) (*RFC3986, error) {

}
