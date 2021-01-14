package repository

import (
	"io"
	"strings"

	"fyne.io/fyne"
)

// splitNonEmpty works exactly like strings.Split(), but only returns non-empty
// components.
func splitNonEmpty(str, sep string) []string {
	components := []string{}
	for _, v := range strings.Split(str, sep) {
		if len(v) > 0 {
			components = append(components, v)
		}
	}
	return components
}

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
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0.0
func GenericParent(u fyne.URI) (fyne.URI, error) {
	p := u.Path()

	if p == "" || p == "/" {
		parent, err := ParseURI(u.String())
		if err != nil {
			return nil, err
		}
		return parent, URIRootError
	}

	components := splitNonEmpty(u.Path(), "/")

	newURI := u.Scheme() + "://" + u.Authority()

	// there will be at least one component, since we know we don't have
	// '/' or ''.
	newURI += "/"
	if len(components) > 1 {
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

	// NOTE: we specifically want to use ParseURI, rather than &uri{},
	// since the repository for the URI we just created might be a
	// CanonicalRepository that implements it's own ParseURI.
	return ParseURI(newURI)
}

// GenericChild can be used as a common-case implementation of
// HierarchicalRepository.Child(). It will create a child URI by separating the
// URI into it's component parts as described in IETF RFC 3986, then appending
// "/" + component to the path, then concatenating the result and parsing it as
// a new URI.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0.0
func GenericChild(u fyne.URI, component string) (fyne.URI, error) {

	// split into components and add the new one
	components := splitNonEmpty(u.Path(), "/")
	components = append(components, component)

	// generate the scheme, authority, and path
	newURI := u.Scheme() + "://" + u.Authority()
	newURI += "/" + strings.Join(components, "/")

	// stick the query and fragment back on the end
	if len(u.Query()) > 0 {
		newURI += "?" + u.Query()
	}
	if len(u.Fragment()) > 0 {
		newURI += "#" + u.Fragment()
	}

	// NOTE: we specifically want to use ParseURI, rather than &uri{},
	// since the repository for the URI we just created might be a
	// CanonicalRepository that implements it's own ParseURI.
	return ParseURI(newURI)
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
// Since: 2.0.0
func GenericCopy(source fyne.URI, destination fyne.URI) error {
	// Look up repositories for the source and destination.
	srcrepo, err := ForURI(source)
	if err != nil {
		return err
	}

	dstrepo, err := ForURI(destination)
	if err != nil {
		return err
	}

	// The destination must be writeable.
	destwrepo, ok := dstrepo.(WriteableRepository)
	if !ok {
		return OperationNotSupportedError
	}

	// Create a reader and a writer.
	srcReader, err := srcrepo.Reader(source)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	dstWriter, err := destwrepo.Writer(destination)
	if err != nil {
		return err
	}
	defer dstWriter.Close()

	// Perform the copy.
	_, err = io.Copy(dstWriter, srcReader)
	return err
}

// GenericMove can be used a common-case implementation of
// MovableRepository.Move(). It will perform the move by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination. Following this, the source
// will be deleted using WriteableRepository.Delete.
//
// For obvious reasons, the source and destination URIs must both be writable.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0.0
func GenericMove(source fyne.URI, destination fyne.URI) error {
	// This looks a lot like GenericCopy(), but I duplicated the code
	// to avoid having to look up the repositories more than once.

	// Look up repositories for the source and destination.
	srcrepo, err := ForURI(source)
	if err != nil {
		return err
	}

	dstrepo, err := ForURI(destination)
	if err != nil {
		return err
	}

	// The source and destination must both be writable, since the source
	// is being deleted, which requires WriteableRepository.
	destwrepo, ok := dstrepo.(WriteableRepository)
	if !ok {
		return OperationNotSupportedError
	}

	srcwrepo, ok := srcrepo.(WriteableRepository)
	if !ok {
		return OperationNotSupportedError
	}

	// Create the reader and writer to perform the copy operation.
	srcReader, err := srcrepo.Reader(source)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	dstWriter, err := destwrepo.Writer(destination)
	if err != nil {
		return err
	}
	defer dstWriter.Close()

	// Perform the copy.
	_, err = io.Copy(dstWriter, srcReader)
	if err != nil {
		return err
	}

	// Finally, delete the source only if the move finished without error.
	return srcwrepo.Delete(source)
}
