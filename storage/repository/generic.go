package repository

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
)

// GenericParent can be used as a common-case implementation of
// HierarchicalRepository.Parent(). It will create a parent URI based on
// IETF RFC3986.
//
// In short, the URI is separated into it's component parts, the path component
// is split along instances of '/', and the trailing element is removed. The
// result is concatenated and parsed as a new URI.
//
// If the URI path is empty or '/', then a duplicate of the URI is returned,
// along with URIRootError.
//
// The parseURI parameter should be a function that you wish GenericParent to
// use to generate any new URIs it creates. Usually, storage.ParseURI is
// suitable here.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since 2.0.0
func GenericParent(u fyne.URI, parseURI func(string) (fyne.URI, error)) (fyne.URI, error) {
	p := u.Path()

	if p == "" || p == "/" {
		parent, err := parseURI(u.String())
		if err != nil {
			return nil, err
		}
		return parent, URIRootError
	}

	components := strings.Split(u.Path(), "/")

	newURI := u.Scheme() + "://" + u.Authority()

	// there will be at least one component, since we know we don't have
	// '/' or ''.
	if len(components) == 1 {
		// the immediate parent is the root
		newURI += "/"
	} else {
		newURI += strings.Join(components[:len(components)-1], "/")
	}

	// stick the query and fragment back on the end
	q := u.Query()
	if len(q) > 0 {
		newURI += "?" + q
	}

	f := u.Fragment()
	if len(f) > 0 {
		newURI += "#" + f
	}

	return parseURI(newURI)
}

// GenericChild can be used as a common-case implementation of
// HierarchicalRepository.Child(). It will create a child URI by separating the
// URI into it's component parts as described in IETF RFC 3986, then appending
// "/" + component to the path, then concatenating the result and parsing it as
// a new URI.
//
// The parseURI parameter should be a function that you wish GenericParent to
// use to generate any new URIs it creates. Usually, storage.ParseURI is
// suitable here.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since 2.0.0
func GenericChild(u fyne.URI, component string, parseURI func(string) (fyne.URI, error)) (fyne.URI, error) {

	// split into components and add the new one
	components := strings.Split(u.Path(), "/")
	components = append(components, component)

	// generate the scheme, authority, and path
	newURI := u.Scheme() + "://" + u.Authority()
	newURI += "/" + strings.Join(components[:len(components)-1], "/")

	// stick the query and fragment back on the end
	if len(u.Query()) > 0 {
		newURI += "?" + u.Query()
	}
	if len(u.Fragment()) > 0 {
		newURI += "#" + u.Fragment()
	}

	return parseURI(newURI)
}

// GenericCopy can be used a common-case implementation of
// CopyableRepository.Copy(). It will perform the copy by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination.
//
// For obvious reasons, the destination URI must have a registered
// WriteableRepository.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since 2.0.0
func GenericCopy(source fyne.URI, destination fyne.URI) error {
	return fmt.Errorf("TODO")
}

// GenericMove can be used a common-case implementation of
// MoveableRepository.Move(). It will perform the move by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination. Following this, the source
// will be deleted using WriteableRepository.Delete.
//
// For obvious reasons, the source and destination URIs must both be writable.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since 2.0.0
func GenericMove(source fyne.URI, destination fyne.URI) error {
	return fmt.Errorf("TODO")
}
