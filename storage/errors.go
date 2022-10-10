package storage

import "errors"

var (
	// ErrAlreadyExists may be thrown by docs. E.g., save a document twice.
	//
	// Since: 2.3
	ErrAlreadyExists = errors.New("document already exists")

	// ErrNotExists may be thrown by docs. E.g., save an unknown document.
	//
	// Since: 2.3
	ErrNotExists = errors.New("document does not exist")
)
