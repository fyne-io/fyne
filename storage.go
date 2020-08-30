package fyne

// Storage is used to manage file storage inside an application sandbox
type Storage interface {
	RootURI() URI
}
