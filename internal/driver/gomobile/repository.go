// +build ios android

package gomobile

import (
	"fyne.io/fyne"
)

type mobileFileRepo struct {
	driver *mobileDriver
}

func (m *mobileFileRepo) Exists(u fyne.URI) (bool, error) {
	return true, nil // TODO check a file exists
}

func (m *mobileFileRepo) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	return m.driver.FileReaderForURI(u)
}

func (m *mobileFileRepo) CanRead(u fyne.URI) (bool, error) {
	return true, nil // TODO check a file can be read
}

func (m *mobileFileRepo) Destroy(string) {
}

func (m *mobileFileRepo) CanList(u fyne.URI) (bool, error) {
	return canListURI(u), nil
}

func (m *mobileFileRepo) List(u fyne.URI) ([]fyne.URI, error) {
	return listURI(u)
}

// TODO add write support (not yet supported on mobile)
