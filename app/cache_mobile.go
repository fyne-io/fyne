//go:build android || ios || mobile

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func rootCacheDir(a fyne.App) string {
	return a.(*fyneApp).storageRoot()
}

func rootCacheURI(a fyne.App) fyne.URI {
	return storage.NewFileURI(rootCacheDir(a))
}
