//go:build !ios

package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

const folderVideos = "Movies"

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	return storage.Child(homeURI, name)
}
