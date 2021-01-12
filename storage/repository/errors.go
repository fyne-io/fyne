package repository

// Declare conformance with Error interface
var _ error = OperationNotSupportedError

type operationNotSupportedError string

// operationNotSupported may be thrown by certain functions in the storage or
// repository packages which operate on URIs if an operation is attempted that
// is not supported for the scheme relevant to the URI, normally because the
// underlying repository has either not implemented the relevant function, or
// has explicitly returned this error.
const OperationNotSupportedError operationNotSupportedError = operationNotSupportedError("Operation not supported for this URI.")

func (e operationNotSupportedError) Error() string {
	return string(OperationNotSupportedError)
}
