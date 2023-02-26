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

	ErrPathNotFound = errors.New("contents at path not found")

	ErrAndroidResourceNotFound = &ErrResourceNotFound{}
)

// ErrResourceNotFound is an error that is specific to Android and indicates that a resource was not found.
type ErrResourceNotFound struct{}

func (e *ErrResourceNotFound) Error() string {
	return "resource not found at URI"
}

func (e *ErrResourceNotFound) Is(target error) bool {
	return target == e || target == ErrPathNotFound || target.Error() == e.Error()
}
