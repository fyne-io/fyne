package repository

// Declare conformance with Error interface
var _ error = OperationNotSupportedError

type operationNotSupportedError string

// OperationNotSupportedError may be thrown by certain functions in the storage
// or repository packages which operate on URIs if an operation is attempted
// that is not supported for the scheme relevant to the URI, normally because
// the underlying repository has either not implemented the relevant function,
// or has explicitly returned this error.
const OperationNotSupportedError operationNotSupportedError = operationNotSupportedError("Operation not supported for this URI.")

func (e operationNotSupportedError) Error() string {
	return string(OperationNotSupportedError)
}

// Declare conformance with Error interface
var _ error = URIRootError

type uriRootError string

// URIRootError should be thrown by fyne.URI implementations when the caller
// attempts to take the parent of the root. This way, downstream code that
// wants to programmatically walk up a URIs parent's will know when to stop
// iterating.
const URIRootError uriRootError = uriRootError("Cannot take the parent of the root element in a URI")

func (e uriRootError) Error() string {
	return string(URIRootError)
}
