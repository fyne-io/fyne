package mobile

import (
	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/storage/repository"
)

// declare conformance with repository types
var (
	_ repository.Repository             = (*mobileFileRepo)(nil)
	_ repository.HierarchicalRepository = (*mobileFileRepo)(nil)
	_ repository.ListableRepository     = (*mobileFileRepo)(nil)
	_ repository.WritableRepository     = (*mobileFileRepo)(nil)
	_ repository.AppendableRepository   = (*mobileFileRepo)(nil)
)

type mobileFileRepo struct{}

func (*mobileFileRepo) CanList(u fyne.URI) (bool, error) {
	return canListURI(u), nil
}

func (*mobileFileRepo) CanRead(fyne.URI) (bool, error) {
	return true, nil // TODO check a file can be read
}

func (*mobileFileRepo) CanWrite(fyne.URI) (bool, error) {
	return true, nil // TODO check a file can be written
}

func (*mobileFileRepo) Child(u fyne.URI, name string) (fyne.URI, error) {
	if u == nil || u.Scheme() != fyne.URISchemeFile {
		return nil, repository.ErrOperationNotSupported
	}

	return repository.GenericChild(u, name)
}

func (*mobileFileRepo) CreateListable(u fyne.URI) error {
	return createListableURI(u)
}

func (*mobileFileRepo) Delete(u fyne.URI) error {
	return deleteURI(u)
}

func (*mobileFileRepo) Destroy(string) {
}

func (*mobileFileRepo) Exists(u fyne.URI) (bool, error) {
	return existsURI(u)
}

func (*mobileFileRepo) List(u fyne.URI) ([]fyne.URI, error) {
	return listURI(u)
}

func (*mobileFileRepo) Parent(u fyne.URI) (fyne.URI, error) {
	if u == nil || u.Scheme() != fyne.URISchemeFile {
		return nil, repository.ErrOperationNotSupported
	}

	return repository.GenericParent(u)
}

func (*mobileFileRepo) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	return fileReaderForURI(u)
}

func (*mobileFileRepo) Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	return fileWriterForURI(u, true)
}

func (*mobileFileRepo) Appender(u fyne.URI) (fyne.URIWriteCloser, error) {
	return fileWriterForURI(u, false)
}
