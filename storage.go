package fyne

// Storage is used to manage file storage inside an application sandbox.
// The files managed by this interface are unique to the current application.
type Storage interface {
	RootURI() URI

	Create(name string) (URIWriteCloser, error)
	Open(name string) (URIReadCloser, error)
	Save(name string) (URIWriteCloser, error)
	Remove(name string) error

	List() []string
}
