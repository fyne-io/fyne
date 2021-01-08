package storage

// Declare conformance with Error interface
var _ error = URIOperationImpossible

type uriOperationImpossible string

// URIOperationImpossible occurs when an operation is attempted on a URI which
// is supported by the underlying implementation, but which violates some
// internal constraints of that particular repository implementation. For
// example, creating an un-representable state might cause this error to occur.
const URIOperationImpossible uriOperationImpossible = uriOperationImpossible("The requested URI operation is not possible.")

func (e uriOperationImpossible) Error() string {
	return string(URIOperationImpossible)
}
