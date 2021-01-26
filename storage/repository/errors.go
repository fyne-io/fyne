package repository

import (
	"errors"
)

var (
	// ErrOperationNotSupported may be thrown by certain functions in the storage
	// or repository packages which operate on URIs if an operation is attempted
	// that is not supported for the scheme relevant to the URI, normally because
	// the underlying repository has either not implemented the relevant function,
	// or has explicitly returned this error.
	//
	// Since: 2.0
	ErrOperationNotSupported = errors.New("operation not supported for this URI")

	// ErrURIRoot should be thrown by fyne.URI implementations when the caller
	// attempts to take the parent of the root. This way, downstream code that
	// wants to programmatically walk up a URIs parent's will know when to stop
	// iterating.
	//
	// Since: 2.0
	ErrURIRoot = errors.New("cannot take the parent of the root element in a URI")
)
