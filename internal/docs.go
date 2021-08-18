package internal

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

type Docs struct {
	RootDocURI fyne.URI
}

func (d *Docs) Create(name string) (fyne.URIWriteCloser, error) {
	err := d.ensureRootExists()
	if err != nil {
		return nil, err
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	exists, err := storage.Exists(u)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("document with name " + name + " already exists")
	}

	return storage.Writer(u)
}

func (d *Docs) List() []string {
	var ret []string
	uris, err := storage.List(d.RootDocURI)
	if err != nil {
		return ret
	}

	for _, u := range uris {
		ret = append(ret, u.Name())
	}

	return ret
}

func (d *Docs) Open(name string) (fyne.URIReadCloser, error) {
	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	return storage.Reader(u)
}

func (d *Docs) Remove(name string) error {
	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return err
	}

	return storage.Delete(u)
}

func (d *Docs) Save(name string) (fyne.URIWriteCloser, error) {
	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	exists, err := storage.Exists(u)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("document with name " + name + " does not exist")
	}

	return storage.Writer(u)
}

func (d *Docs) ensureRootExists() error {
	exists, err := storage.Exists(d.RootDocURI)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return storage.CreateListable(d.RootDocURI)
}
