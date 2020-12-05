// +build !ios,!android

package gomobile

import "fyne.io/fyne"

func canListURI(fyne.URI) bool {
	return false
}

func listURI(fyne.URI) ([]fyne.URI, error) {
	return nil, nil
}
