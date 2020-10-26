// +build ios

package gomobile

import "fyne.io/fyne"

func uriCanList(fyne.URI) bool {
	return false
}

func uriList(fyne.URI) ([]fyne.URI, error) {
	return nil, nil
}
