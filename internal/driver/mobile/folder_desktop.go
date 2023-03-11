//go:build !ios && !android
// +build !ios,!android

package mobile

import (
	"fyne.io/fyne/v2"
)

func canListURI(fyne.URI) bool {
	// no-op as we use the internal FileRepository
	return false
}

func createListableURI(fyne.URI) error {
	// no-op as we use the internal FileRepository
	return nil
}

func listURI(fyne.URI) ([]fyne.URI, error) {
	// no-op as we use the internal FileRepository
	return nil, nil
}
