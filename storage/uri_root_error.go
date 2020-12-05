package storage

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
