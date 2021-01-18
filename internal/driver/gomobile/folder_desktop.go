// +build !ios,!android

package gomobile

import (
	"fyne.io/fyne/v2"
)

func canListURI(fyne.URI) bool {
	// no-op as we use the internal FileRepository
	return false
}

func listURI(fyne.URI) ([]fyne.URI, error) {
	// no-op as we use the internal FileRepository
	return nil, nil
}
