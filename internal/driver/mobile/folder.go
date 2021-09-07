package mobile

import (
	"fmt"

	"fyne.io/fyne/v2"
)

type lister struct {
	fyne.URI
}

func (l *lister) List() ([]fyne.URI, error) {
	return listURI(l)
}

func listerForURI(uri fyne.URI) (fyne.ListableURI, error) {
	if !canListURI(uri) {
		return nil, fmt.Errorf("specified URI is not listable")
	}

	return &lister{uri}, nil
}
