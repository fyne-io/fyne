// +build !ios,!android

package gomobile

import (
	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func canListURI(u fyne.URI) bool {
	listable, _ := storage.CanList(u)
	return listable
}

func listURI(u fyne.URI) ([]fyne.URI, error) {
	return storage.List(u)
}
