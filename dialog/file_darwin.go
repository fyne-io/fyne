//go:build !ios && !android && !wasm && !js

package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	return storage.Child(homeURI, name)
}
