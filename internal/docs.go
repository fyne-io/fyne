package internal

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

var errNoAppID = errors.New("storage API requires a unique ID, use app.NewWithID()")

// Docs is an internal implementation of the document features of the Storage interface.
// It is based on top of the current `file` repository and is rooted at RootDocURI.
type Docs struct {
	RootDocURI fyne.URI
}

// Create will create a new document ready for writing, you must write something and close the returned writer
// for the create process to complete.
// If the document for this app with that name already exists an error will be returned.
func (d *Docs) Create(name string) (fyne.URIWriteCloser, error) {
	if d.RootDocURI == nil {
		return nil, errNoAppID
	}

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

// List returns all documents that have been saved by the current application.
// Remember to use `app.NewWithID` so that your storage is unique.
func (d *Docs) List() []string {
	if d.RootDocURI == nil {
		return nil
	}

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

// Open will grant access to the contents of the named file. If an error occurs it is returned instead.
func (d *Docs) Open(name string) (fyne.URIReadCloser, error) {
	if d.RootDocURI == nil {
		return nil, errNoAppID
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return nil, err
	}

	return storage.Reader(u)
}

// Remove will delete the document with the specified name, if it exists
func (d *Docs) Remove(name string) error {
	if d.RootDocURI == nil {
		return errNoAppID
	}

	u, err := storage.Child(d.RootDocURI, name)
	if err != nil {
		return err
	}

	return storage.Delete(u)
}

// Save will open a document ready for writing, you close the returned writer for the save to complete.
// If the document for this app with that name does not exist an error will be returned.
func (d *Docs) Save(name string) (fyne.URIWriteCloser, error) {
	if d.RootDocURI == nil {
		return nil, errNoAppID
	}

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
