package storage

import "errors"

var (
	// ErrNotExists may be thrown by docs. E.g., save an unknown document.
	//
	// Since: 2.3
	ErrNotExists = errors.New("document does not exist")
)
