package storage

// Declare conformance with Error interface
var _ error = URIOperationNotSupportedError

type uriOperationNotSupportedError string

// URIOperationNotSupported may be thrown by certain functions in the storage
// package which operate on URIs if an operation is attempted that is not
// supported for the scheme relevant to the URI, normally because the
// underlying repository has either not implemented the relevant function, or
// has explicitly returned this error.
const URIOperationNotSupportedError uriOperationNotSupportedError = uriOperationNotSupportedError("Operation not supported for this URI.")

func (e uriOperationNotSupportedError) Error() string {
	return string(URIOperationNotSupportedError)
}
