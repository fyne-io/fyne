package repository

import (
	"fmt"

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
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since 2.0.0
func GenericParent(u fyne.URI) (fyne.URI, error) {
	return nil, fmt.Errorf("TODO")
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
// Since 2.0.0
func GenericChild(u fyne.URI, component string) (fyne.URI, error) {
	return nil, fmt.Errorf("TODO")
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
